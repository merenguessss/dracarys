package pool

import (
	"context"
	"net"
)

type ConnPool interface {
	Get(ctx context.Context, network string, address string) (net.Conn, error)
}
