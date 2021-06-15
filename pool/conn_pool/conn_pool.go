package conn_pool

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrorOutMaxConn = errors.New("out of maximum connection error ")
	ErrorConnClose  = errors.New("connection is close .... ")
)

var DefaultConnPool = newDefaultPool()

var newDefaultPool = func(options ...Option) Pool {
	opts := &Options{
		MaxIdle:     10000,
		MaxActive:   1000,
		CoreIdle:    1000,
		Wait:        true,
		IdleTimeout: 2 * time.Minute,
		DialTimeout: 800 * time.Millisecond,
	}
	for _, o := range options {
		o(opts)
	}

	return &pool{
		opts:    opts,
		connMap: &sync.Map{},
	}
}

var poolMap = make(map[string]Pool)

func init() {
	RegisterPool("default", DefaultConnPool)
}

func RegisterPool(name string, pool Pool) {
	if poolMap == nil {
		poolMap = make(map[string]Pool)
	}
	poolMap[name] = pool
}

type Pool interface {
	Get(ctx context.Context, network string, address string) (net.Conn, error)
}

type pool struct {
	opts    *Options
	connMap *sync.Map
}

// Get 连接池对外接口.
func (p *pool) Get(ctx context.Context, network string, address string) (net.Conn, error) {
	if v, ok := p.connMap.Load(address); ok {
		if conn, ok := v.(*channelPool); ok {
			return conn.get(ctx)
		}
	}

	dial := func(ctx context.Context) (net.Conn, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		timeout := p.opts.DialTimeout
		v, ok := ctx.Deadline()
		if ok {
			timeout = time.Until(v)
		}
		return net.DialTimeout(network, address, timeout)
	}

	newPool := &channelPool{
		maxIdle:     p.opts.MaxIdle,
		coreIdle:    p.opts.CoreIdle,
		maxActive:   p.opts.MaxActive,
		wait:        p.opts.Wait,
		idleTimeout: p.opts.IdleTimeout,
		dialTimeout: p.opts.DialTimeout,
		dial:        dial,
		ch:          make(chan *poolConn, p.opts.MaxActive),
	}

	v, ok := p.connMap.LoadOrStore(address, newPool)
	if !ok {
		newPool.checkFreeConn(2*time.Second, newPool.checkConn)
	}
	return v.(*channelPool).get(ctx)
}

// channelPool 连接池
type channelPool struct {
	net.Conn
	maxIdle     int
	coreIdle    int
	maxActive   int
	wait        bool
	idleTimeout time.Duration
	dialTimeout time.Duration
	dial        func(ctx context.Context) (net.Conn, error)
	ch          chan *poolConn
	mu          sync.RWMutex
	activeNum   int32
}

// get 从连接池中获取一个conn连接.
func (cp *channelPool) get(ctx context.Context) (*poolConn, error) {
	if cp.ch == nil {
		return nil, ErrorConnClose
	}

	// 无限循环保证一定会取得一个连接，无法取得连接的情况只会是拒绝策略.
	for {
		select {
		case pc := <-cp.ch:
			if pc != nil && !pc.unusable {
				return pc, nil
			}
		default:
			cp.mu.RLock()
			if cp.ch == nil {
				cp.ch = make(chan *poolConn, cp.maxActive)
			}

			if int(cp.activeNum)+len(cp.ch) >= cp.maxIdle {
				// TODO 拒绝策略
				cp.mu.RUnlock()
				return nil, ErrorOutMaxConn
			}
			pc, err := cp.getPoolConn(ctx)
			if err != nil {
				cp.mu.RUnlock()
				return nil, err
			}
			atomic.AddInt32(&cp.activeNum, 1)
			return pc, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// getPoolConn 获取PoolConn连接.
func (cp *channelPool) getPoolConn(ctx context.Context) (*poolConn, error) {
	conn, err := cp.dial(ctx)
	if err != nil {
		return nil, err
	}
	pc := cp.wrapConn(conn)
	return pc, nil
}

// checkFreeConn 检查连接池中的空闲连接.
func (cp *channelPool) checkFreeConn(interval time.Duration, checker func(conn *poolConn) bool) {
	if interval > 0 && checker != nil {
		go func() {
			for {
				time.Sleep(interval)
				length := len(cp.ch)

				for i := 0; i < length; i++ {
					select {
					case conn := <-cp.ch:
						if checker(conn) {
							conn.recycling()
							_ = conn.Close()
						} else {
							_ = cp.Put(conn)
						}
					default:
					}
				}
			}
		}()
	}
}

// checkConn 检查连接是否有效,无效即断开连接.
func (cp *channelPool) checkConn(conn *poolConn) bool {
	// 判断连接超时时间是否存在并且 连接超时时间+连接创建时间<当前时间? 成立 即返回true.
	if cp.idleTimeout > 0 && conn.t.Add(cp.idleTimeout).Before(time.Now()) {
		return true
	}

	if !conn.remoteSend() {
		return true
	}
	return false
}

// Close 关闭线程池.
func (cp *channelPool) Close() {
	cp.mu.Lock()
	ch := cp.ch
	cp.ch = nil
	cp.dial = nil
	cp.mu.Unlock()

	if ch == nil {
		return
	}
	close(ch)
	for conn := range ch {
		conn.recycling()
		conn.Close()
	}
}

// Put 像线程池中加入一个连接.
func (cp *channelPool) Put(conn *poolConn) error {
	if conn == nil {
		return ErrorConnClose
	}
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	if cp.ch == nil || len(cp.ch) >= cp.coreIdle {
		conn.recycling()
		return conn.Close()
	}

	select {
	case cp.ch <- conn:
		return nil
	default:
		conn.recycling()
		return conn.Close()
	}
}

type poolConn struct {
	net.Conn
	c        *channelPool
	unusable bool
	// 连接时间
	t           time.Time
	dialTimeout time.Duration
	mu          sync.RWMutex
}

var oneByte = make([]byte, 1)

// remoteSend 尝试发送一个字节.
func (conn *poolConn) remoteSend() bool {
	_ = conn.Conn.SetDeadline(time.Now().Add(time.Millisecond))
	defer func() {
		_ = conn.Conn.SetDeadline(time.Time{})
	}()

	if n, err := conn.Read(oneByte); err == io.EOF || n == 0 {
		return false
	}
	return true
}

// Close 关闭一个连接，如果标识符unusable为true,关闭连接。否则重新尝试放入连接池.
func (conn *poolConn) Close() error {
	conn.mu.RLock()
	defer conn.mu.RUnlock()
	if conn.unusable {
		if conn.Conn != nil {
			return conn.Conn.Close()
		}
	}
	conn.Conn.SetDeadline(time.Time{})
	return conn.c.Put(conn)
}

// Recycling 将连接标记为正在回收,将其中的unusable字段标记为true,准备回收掉.
func (conn *poolConn) recycling() {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	conn.unusable = true
	atomic.AddInt32(&conn.c.activeNum, -1)
}

// wrapConn 将net.Conn类型封装为poolConn类型.
func (cp *channelPool) wrapConn(conn net.Conn) *poolConn {
	pc := &poolConn{
		c:           cp,
		dialTimeout: cp.dialTimeout,
		Conn:        conn,
		t:           time.Now(),
	}
	return pc
}

// Write 为poolConn重写Conn接口的Write方法.
func (conn *poolConn) Write(b []byte) (int, error) {
	if conn.unusable {
		return 0, ErrorConnClose
	}
	n, err := conn.Conn.Write(b)
	if err != nil {
		conn.recycling()
		conn.Conn.Close()
	}
	return n, err
}

// Read 为poolConn重写Conn接口的Read方法.
func (conn *poolConn) Read(b []byte) (int, error) {
	if conn.unusable {
		return 0, ErrorConnClose
	}
	n, err := conn.Conn.Read(b)
	if err != nil {
		conn.recycling()
		conn.Conn.Close()
	}
	return n, err
}
