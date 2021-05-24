package server

import (
	"time"

	"github.com/merenguessss/dracarys-go/interceptor"
	"github.com/merenguessss/dracarys-go/plugin"
)

type Options struct {
	address         string
	network         string
	keepAlivePeriod time.Duration
	PluginFactory   plugin.Factory
	serializerType  string
	codecType       string
	beforeHandle    []interceptor.Interceptor
	afterHandle     []interceptor.Interceptor
}

type Option func(*Options)

func WithAddress(addr string) Option {
	return func(o *Options) {
		o.address = addr
	}
}

func WithNetWork(network string) Option {
	return func(o *Options) {
		o.network = network
	}
}

func WithKeepAlivePeriod(t time.Duration) Option {
	return func(o *Options) {
		o.keepAlivePeriod = t
	}
}
