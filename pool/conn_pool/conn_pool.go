package conn_pool

import (
	"context"
	"io"
	"net"
	"sync"
	"time"
)

type Pool interface {
	Get(ctx context.Context, network string, address string) (net.Conn, error)
}

type pool struct {
	opts    *Options
	connMap sync.Map
}

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
			timeout = v.Sub(time.Now())
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
	}

	v, ok := p.connMap.LoadOrStore(address, newPool)
	if !ok {
		newPool.checkFreeConn(5*time.Second, newPool.checkConn)
	}
	return v.(*channelPool).get(ctx)
}

type channelPool struct {
	net.Conn
	maxIdle     int
	coreIdle    int
	maxActive   int
	wait        bool
	idleTimeout time.Duration
	dialTimeout time.Duration
	dial        func(ctx context.Context) (net.Conn, error)
	ch          chan struct{}
	mu          sync.RWMutex
	activeNum   int
	connList    chan *poolConn
}

// get 从连接池中获取一个conn连接.
func (cp *channelPool) get(ctx context.Context) (*poolConn, error) {

	return nil, nil
}

// checkFreeConn 检查连接池中的空闲连接.
func (cp *channelPool) checkFreeConn(interval time.Duration, checker func(conn *poolConn) bool) {
	if interval > 0 && checker != nil {
		go func() {
			for {
				time.Sleep(interval)
				// TODO 检查chan中的连接是否空闲
				cp.mu.Unlock()
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

var oneByte = make([]byte, 1)

// remoteSend 尝试发送一个字节.
func (conn *poolConn) remoteSend() bool {
	_ = conn.Conn.SetDeadline(time.Now().Add(time.Millisecond))
	defer func() {
		_ = conn.Conn.SetDeadline(time.Time{})
	}()

	if n, err := conn.Conn.Read(oneByte); err == io.EOF || n == 0 {
		return false
	}
	return true
}

type poolConn struct {
	net.Conn
	c        *channelPool
	unusable bool
	// 连接时间
	t          time.Time
	prev, next *poolConn
}
