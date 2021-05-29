package dracarys

import "github.com/merenguessss/dracarys-go/server"

func NewServer(opts ...server.Option) *server.Server {
	srv := &server.Server{
		ServiceMap: make(map[string]server.Service),
		Options:    &server.Options{},
	}

	for _, o := range opts {
		o(srv.Options)
	}
	return srv
}
