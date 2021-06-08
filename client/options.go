package client

import (
	"github.com/merenguessss/dracarys-go/interceptor"
	"github.com/merenguessss/dracarys-go/plugin"
)

type Options struct {
	ClientName        string
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
	CompressType      string
}

type Option func(*Options)

func WithAddr(addr string) Option {
	return func(options *Options) {
		options.Addr = addr
	}
}

func WithClientName(clientName string) Option {
	return func(options *Options) {
		options.ClientName = clientName
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

func WithDisableConnPool(b bool) Option {
	return func(options *Options) {
		options.DisableConnPool = b
	}
}

func WithEnableMultiplexed(b bool) Option {
	return func(options *Options) {
		options.EnableMultiplexed = b
	}
}

func WithCodecType(codecType string) Option {
	return func(options *Options) {
		options.codecType = codecType
	}
}

func WithSerializerType(serializerType string) Option {
	return func(options *Options) {
		options.serializerType = serializerType
	}
}

func WithBeforeHandle(h []interceptor.Interceptor) Option {
	return func(options *Options) {
		options.beforeHandle = h
	}
}

func WithAfterHandle(h []interceptor.Interceptor) Option {
	return func(options *Options) {
		options.afterHandle = h
	}
}

func WithPluginFactory(p plugin.Factory) Option {
	return func(options *Options) {
		options.PluginFactory = p
	}
}

func WithCompressType(compressType string) Option {
	return func(options *Options) {
		options.CompressType = compressType
	}
}
