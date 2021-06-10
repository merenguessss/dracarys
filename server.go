package dracarys

import (
	"github.com/merenguessss/dracarys-go/config"
	"github.com/merenguessss/dracarys-go/server"
)

func NewServer(opts ...server.Option) *server.Server {
	o, err := config.GetServer()
	if err != nil {
		o = &server.Options{}
	}
	srv := &server.Server{
		ServiceMap: make(map[string]server.Service),
		Options:    o,
	}

	for _, o := range opts {
		o(srv.Options)
	}
	return srv
}
