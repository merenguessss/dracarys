package transport

import (
	"context"
	"time"
)

type ServerOptions struct {
	network         Network
	address         string
	keepAlivePeriod time.Duration
	handler         Handler
}

type ServerOption func(*ServerOptions)

type Handler interface {
	handle(context.Context, []byte) ([]byte, error)
}

func WithNetwork(network string) ServerOption {
	return func(so *ServerOptions) {
		so.network = Network(network)
	}
}

func WithAddress(addr string) ServerOption {
	return func(so *ServerOptions) {
		so.address = addr
	}
}

func WithKeepAlivePeriod(t time.Duration) ServerOption {
	return func(so *ServerOptions) {
		so.keepAlivePeriod = t
	}
}

func WithHandler(handler Handler) ServerOption {
	return func(so *ServerOptions) {
		so.handler = handler
	}
}
