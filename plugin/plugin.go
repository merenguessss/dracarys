package plugin

import (
	"github.com/merenguessss/dracarys/plugin/balance"
	"github.com/merenguessss/dracarys/plugin/selector"
)

type Factory struct {
	s       selector.Selector
	b       balance.Balancer
	options *Options
}

var New = func(opts *Options) *Factory {
	return &Factory{
		options: opts,
	}
}

var DefaultConfig = func() *Options {
	return &Options{
		Selector: &selector.Options{
			SelectorName:    "consul",
			Address:         "127.0.0.1:8500",
			EnableHeartbeat: true,
			Scheme:          "http",
			HeartbeatOptions: &selector.HeartbeatOptions{
				Port:                           "8001",
				Timeout:                        "5s",
				Interval:                       "5s",
				DeregisterCriticalServiceAfter: "20s",
			},
		},
		Balancer: &balance.Options{
			BalancerName: balance.WeightPoll,
		},
	}
}

// Setup 加载插件工厂中的插件.
func (f *Factory) Setup(opts ...Option) error {
	for _, o := range opts {
		o(f.options)
	}

	f.s = selector.Get(f.options.Selector.SelectorName)
	err := f.s.LoadConfig(f.options.Selector)
	if err != nil {
		return err
	}

	f.b = balance.Get(f.options.Balancer.BalancerName)
	return nil
}

// GetSelector 获取服务发现插件.
func (f *Factory) GetSelector() selector.Selector {
	return f.s
}

// GetBalancer 获取负载均衡插件.
func (f *Factory) GetBalancer() balance.Balancer {
	return f.b
}
