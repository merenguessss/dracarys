package transport

import (
	"context"
)

type ClientTransport interface {
	Send(context.Context, []byte, ...ClientOption) error
	multiplexed(context.Context, []byte) error
	sendTCP(context.Context, []byte) error
	sendUDP(context.Context, []byte) error
}

type Network string

const (
	TCP Network = "tcp"
	UDP Network = "udp"
)
