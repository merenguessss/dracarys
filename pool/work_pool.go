package pool

import (
	"net"
	"time"
)

type WorkPool interface {
	Start()
	clean(*[]*WorkChan)
	Serve(net.Conn)
	getWorkChan() *WorkChan
	workerFunc(*WorkChan)
	release(*WorkChan) bool
}

type WorkChan struct {
	LastUseTime time.Time
	ch          chan net.Conn
}

type RejectPolicy interface {
}
