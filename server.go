package dracarys

import (
	"github.com/merenguessss/dracarys-go/config"
	"github.com/merenguessss/dracarys-go/server"
)

// NewServer 创建server端,用于注册service以及监听.
func NewServer(opts ...server.Option) *server.Server {
	o, err := config.GetServer()
	if err != nil {
		// todo log or panic
		o = &server.Options{}
	}

	for _, op := range opts {
		op(o)
	}

	if err = o.PluginFactory.Setup(o.PluginFactoryOptions...); err != nil {
		// todo log or panic
	}

	srv := &server.Server{
		ServiceMap: make(map[string]server.Service),
		Options:    o,
	}
	return srv
}
