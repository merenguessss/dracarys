package client

import (
	"github.com/merenguessss/dracarys-go/interceptor"
	"github.com/merenguessss/dracarys-go/plugin"
)

type Options struct {
	ClientName        string `yaml:"name"`
	Addr              string `yaml:"address"`
	SerializerType    string `yaml:"serializer"`
	CodecType         string `yaml:"codec"`
	EnableMultiplexed bool   `yaml:"enable_multiplexed"`
	DisableConnPool   bool   `yaml:"disable_connection_pool"`
	NetWork           string `yaml:"network"`
	CompressType      string `yaml:"compress"`
	serviceName       string
	methodName        string
	PluginFactory     plugin.Factory
	beforeHandle      []interceptor.ClientInvoker
	afterHandle       []interceptor.ClientInvoker
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
		options.serviceName = serviceName
	}
}

func WithMethod(methodName string) Option {
	return func(options *Options) {
		options.methodName = methodName
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
		options.CodecType = codecType
	}
}

func WithSerializerType(serializerType string) Option {
	return func(options *Options) {
		options.SerializerType = serializerType
	}
}

func WithBeforeHandle(h []interceptor.ClientInvoker) Option {
	return func(options *Options) {
		options.beforeHandle = h
	}
}

func WithAfterHandle(h []interceptor.ClientInvoker) Option {
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
