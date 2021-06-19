package server

import (
	"time"

	"github.com/merenguessss/dracarys-go/interceptor"
	"github.com/merenguessss/dracarys-go/plugin"
)

type Options struct {
	ServerName           string        `yaml:"name"`
	Port                 string        `yaml:"port"`
	Network              string        `yaml:"network"`
	KeepAlivePeriod      time.Duration `yaml:"keep_alive_period"`
	SerializerType       string        `yaml:"serializer"`
	CodecType            string        `yaml:"codec"`
	CompressType         string        `yaml:"compress"`
	Address              string
	PluginFactoryOptions []plugin.Option
	PluginFactory        *plugin.Factory
	beforeHandle         []interceptor.ServerHandler
	afterHandle          []interceptor.ServerHandler
}

type Option func(*Options)

func WithPort(port string) Option {
	return func(o *Options) {
		o.Address = "127.0.0.1:" + port
		o.Port = port
	}
}

func WithAddress(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

func WithNetWork(network string) Option {
	return func(o *Options) {
		o.Network = network
	}
}

func WithCodecType(codecType string) Option {
	return func(o *Options) {
		o.CodecType = codecType
	}
}

func WithBeforeHandle(iter ...interceptor.ServerHandler) Option {
	return func(o *Options) {
		o.beforeHandle = append(o.beforeHandle, iter...)
	}
}

func WithAfterHandle(iter ...interceptor.ServerHandler) Option {
	return func(o *Options) {
		o.afterHandle = append(o.afterHandle, iter...)
	}
}

func WithSerializerType(serializerType string) Option {
	return func(o *Options) {
		o.SerializerType = serializerType
	}
}

func WithKeepAlivePeriod(t time.Duration) Option {
	return func(o *Options) {
		o.KeepAlivePeriod = t
	}
}

func WithPluginFactory(p *plugin.Factory) Option {
	return func(o *Options) {
		o.PluginFactory = p
	}
}

func WithPluginFactoryOptions(o []plugin.Option) Option {
	return func(options *Options) {
		options.PluginFactoryOptions = o
	}
}
