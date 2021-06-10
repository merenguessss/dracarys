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

func WithCodecType(codecType string) Option {
	return func(o *Options) {
		o.codecType = codecType
	}
}

func WithBeforeHandle(iter ...interceptor.Interceptor) Option {
	return func(o *Options) {
		o.beforeHandle = append(o.beforeHandle, iter...)
	}
}

func WithAfterHandle(iter ...interceptor.Interceptor) Option {
	return func(o *Options) {
		o.afterHandle = append(o.afterHandle, iter...)
	}
}

func WithSerializerType(serializerType string) Option {
	return func(o *Options) {
		o.serializerType = serializerType
	}
}

func WithKeepAlivePeriod(t time.Duration) Option {
	return func(o *Options) {
		o.keepAlivePeriod = t
	}
}
