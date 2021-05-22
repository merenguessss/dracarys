package transport

import (
	"context"
	"time"
)

type ServerOptions struct {
	network         Network
	address         string
	keepAlivePeriod time.Duration
	handle          Handler
}

type ServerOption func(*ServerOptions)

type Handler interface {
	handle(context.Context, []byte) ([]byte, error)
}
