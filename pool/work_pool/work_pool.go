package work_pool

import (
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/merenguessss/dracarys/log"
	"github.com/merenguessss/dracarys/pool"
)

type Pool interface {
	Start()
	Stop()
	clean(*[]*workChan)
	Serve(net.Conn)
	getWorkChan() *workChan
	workerFunc(*workChan)
	release(*workChan) bool
}

type workChan struct {
	LastUseTime time.Time
	ch          chan net.Conn
}

type workpool struct {
	initWorkersCount int
	maxWorkerTime    time.Duration
	mu               sync.Mutex
	rejectPolicy     pool.RejectPolicy
	ready            []*workChan
	maxWorkersCount  int
	worksCount       int
	stopCh           chan struct{}
	handle           func(net.Conn) error
	workChanPool     sync.Pool
	stop             bool
}

var workerChanCap = func() int {
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}
	return 1
}()

func (wp *workpool) Start() {
	if wp.stopCh != nil {
		panic("stop channel not nil")
	}

	wp.stopCh = make(chan struct{})
	stopCh := wp.stopCh
	wp.workChanPool.New = func() interface{} {
		return &workChan{
			ch: make(chan net.Conn, workerChanCap),
		}
	}

	go func() {
		cleanCh := make([]*workChan, 0)
		for {
			wp.clean(&cleanCh)
			select {
			case <-stopCh:
				break
			default:
				time.Sleep(wp.getMaxWorkerTime())
			}
		}
	}()
}

func (wp *workpool) Stop() {
	if wp.stopCh == nil {
		panic("stop channnel is nil")
	}
	close(wp.stopCh)
	wp.stopCh = nil

	wp.mu.Lock()
	ready := wp.ready
	for i := range wp.ready {
		ready[i].ch <- nil
		ready[i] = nil
	}
	ready = ready[:0]
	wp.stop = true
	wp.mu.Unlock()
}

func (wp *workpool) clean(wc *[]*workChan) {
	maxWorkerTime := wp.getMaxWorkerTime()
	criticalTime := time.Now().Add(-maxWorkerTime)

	wp.mu.Lock()
	ready := wp.ready
	n := len(ready)
	l, r, mid := 0, n-1, 0
	for l <= r {
		mid = l + (r-l)/2
		if criticalTime.After(ready[mid].LastUseTime) {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	i := r + 1
	if i == 0 {
		wp.mu.Unlock()
		return
	}

	*wc = append((*wc)[:0], ready[:i]...)
	m := copy(ready, ready[i:])
	for i = m; i < n; i++ {
		ready[i] = nil
	}
	wp.ready = ready[:m]
	wp.mu.Unlock()

	tmp := *wc
	for t := range tmp {
		tmp[t].ch = nil
		tmp[t] = nil
	}
}

func (wp *workpool) getMaxWorkerTime() time.Duration {
	if wp.maxWorkerTime <= 0 {
		return 10 * time.Second
	}
	return wp.maxWorkerTime
}

func (wp *workpool) Serve(conn net.Conn) bool {
	ch := wp.getWorkChan()
	if ch == nil {
		return false
	}
	ch.ch <- conn
	return true
}

func (wp *workpool) getWorkChan() *workChan {
	var ch *workChan
	wp.mu.Lock()
	ready := wp.ready
	n := len(ready) - 1
	if n < 0 {
		if wp.worksCount < wp.maxWorkersCount {
			wp.worksCount++
			wp.mu.Unlock()
			return wp.createWorkChan()
		}
		wp.mu.Unlock()
		return ch
	}

	ch = ready[n]
	ready[n] = nil
	ready = ready[:0]
	wp.mu.Unlock()
	return ch
}

func (wp *workpool) createWorkChan() *workChan {
	ch := wp.workChanPool.Get().(*workChan)
	go func() {
		wp.workerFunc(ch)
		wp.workChanPool.Put(ch)
	}()
	return ch
}

func (wp *workpool) release(ch *workChan) bool {
	ch.LastUseTime = time.Now()
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.stop {
		return false
	}
	wp.ready = append(wp.ready, ch)
	return true
}

func (wp *workpool) workerFunc(ch *workChan) {
	var conn net.Conn
	for conn = range ch.ch {
		if conn == nil {
			break
		}

		if err := wp.handle(conn); err != nil {
			log.Error(err)
		}
		_ = conn.Close()
		conn = nil
		if !wp.release(ch) {
			break
		}
	}

	wp.mu.Lock()
	wp.worksCount--
	wp.mu.Unlock()
}
