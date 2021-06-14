package plugin

import "github.com/merenguessss/dracarys-go/plugin/selector"

type Options struct {
	Selector *selector.Options `yaml:"selector"`
}

type Option func(*Options)

func WithSelectorName(name string) Option {
	return func(o *Options) {
		o.Selector.SelectorName = name
	}
}

func WithSelectorAddress(address string) Option {
	return func(o *Options) {
		o.Selector.Address = address
	}
}

func WithSelectorScheme(scheme string) Option {
	return func(o *Options) {
		o.Selector.Scheme = scheme
	}
}

func WithSelectorEnableHeartbeat(b bool) Option {
	return func(o *Options) {
		o.Selector.EnableHeartbeat = b
	}
}

func WithHeartbeatPort(port string) Option {
	return func(o *Options) {
		o.Selector.HeartbeatOptions.Port = port
	}
}

func WithHeartbeatTimeout(timeout string) Option {
	return func(o *Options) {
		o.Selector.HeartbeatOptions.Timeout = timeout
	}
}

func WithHeartbeatInterval(interval string) Option {
	return func(o *Options) {
		o.Selector.HeartbeatOptions.Interval = interval
	}
}

func WithSelectorDeregisterCriticalServiceAfter(t string) Option {
	return func(o *Options) {
		o.Selector.HeartbeatOptions.DeregisterCriticalServiceAfter = t
	}
}
