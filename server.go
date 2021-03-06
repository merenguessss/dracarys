package dracarys

import (
	"github.com/merenguessss/dracarys/config"
	"github.com/merenguessss/dracarys/log"
	"github.com/merenguessss/dracarys/server"
)

// NewServer 创建server端,用于注册service以及监听.
func NewServer(opts ...server.Option) *server.Server {
	o, err := config.GetServer()
	if err != nil {
		log.Fatal(err)
	}

	for _, op := range opts {
		op(o)
	}

	if err = o.PluginFactory.Setup(o.PluginFactoryOptions...); err != nil {
		log.Fatal(err)
	}

	srv := &server.Server{
		ServiceMap: make(map[string]server.Service),
		Options:    o,
	}
	srv.PrintLogo()
	return srv
}
