package transport

import (
	"github.com/merenguessss/Dracarys-go/pool/conn_pool"
)

type ClientOptions struct {
	Addr              string
	Network           Network
	pool              conn_pool.Pool
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

func WithConnPool(connPool conn_pool.Pool) ClientOption {
	return func(options *ClientOptions) {
		options.pool = connPool
	}
}

func WithEnableMultiplexed(enableMultiplexed bool) ClientOption {
	return func(options *ClientOptions) {
		options.EnableMultiplexed = enableMultiplexed
	}
}

func WithDisableConnPool(disableConnPool bool) ClientOption {
	return func(options *ClientOptions) {
		options.DisableConnPool = disableConnPool
	}
}
