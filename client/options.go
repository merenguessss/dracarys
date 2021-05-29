package client

import (
	"github.com/merenguessss/dracarys-go/interceptor"
	"github.com/merenguessss/dracarys-go/plugin"
)

type Options struct {
	ServiceName       string
	MethodName        string
	Addr              string
	PluginFactory     plugin.Factory
	beforeHandle      []interceptor.Interceptor
	afterHandle       []interceptor.Interceptor
	serializerType    string
	codecType         string
	EnableMultiplexed bool
	DisableConnPool   bool
	NetWork           string
}

type Option func(*Options)

func WithAddr(addr string) Option {
	return func(options *Options) {
		options.Addr = addr
	}
}

func WithService(serviceName string) Option {
	return func(options *Options) {
		options.ServiceName = serviceName
	}
}

func WithMethod(methodName string) Option {
	return func(options *Options) {
		options.MethodName = methodName
	}
}

func WithNetWork(network string) Option {
	return func(options *Options) {
		options.NetWork = network
	}
}
