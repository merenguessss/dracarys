package transport

import "context"

type ClientTransport interface {
	Send(context.Context, []byte, ...ClientOptions) error
}

type ClientOptions struct {
	Addr string
}
