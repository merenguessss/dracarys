package transport

import (
	"context"
)

type ClientTransport interface {
	Send(context.Context, []byte, ...ClientOption) ([]byte, error)
	multiplexed(context.Context, []byte) ([]byte, error)
	sendTCP(context.Context, []byte) ([]byte, error)
	sendUDP(context.Context, []byte) ([]byte, error)
}

type Network string

const (
	TCP Network = "tcp"
	UDP Network = "udp"
)
