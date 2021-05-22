package transport

import (
	"context"
)

type ServerTransport interface {
	ListenAndServe(context.Context, ...ServerOption) error
}

type ClientTransport interface {
	Send(context.Context, []byte, ...ClientOption) ([]byte, error)
}

type Network string

const (
	TCP Network = "tcp"
	UDP Network = "udp"
)
