package transport

import (
	"github.com/merenguessss/Dracarys-go/pool/conn_pool"
)

type ClientOptions struct {
	Addr              string
	Network           Network
	pool              conn_pool.ConnPool
	EnableMultiplexed bool
	DisableConnPool   bool
}

type ClientOption func(*ClientOptions)

func WithAddr(addr string) ClientOption {
	return func(options *ClientOptions) {
		options.Addr = addr
	}
}

func WithNetWork(network Network) ClientOption {
	return func(options *ClientOptions) {
		options.Network = network
	}
}

func WithConnPool(connPool conn_pool.ConnPool) ClientOption {
	return func(options *ClientOptions) {
		options.pool = connPool
	}
}
