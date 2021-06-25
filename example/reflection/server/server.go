package main

import (
	"context"
	"time"

	"github.com/merenguessss/dracarys"
	"github.com/merenguessss/dracarys/log"
	"github.com/merenguessss/dracarys/server"
)

func main() {
	opts := []server.Option{
		server.WithAddress("localhost:8000"),
		server.WithNetWork("tcp"),
		server.WithKeepAlivePeriod(time.Second * 200),
		server.WithSerializerType("json"),
	}
	srv := dracarys.NewServer(opts...)
	err := srv.RegisterService("Hello", &Hello{})
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Serve(); err != nil {
		log.Error(err)
	}
}

type Hello struct {
}

func (h *Hello) World(ctx context.Context, s string) (string, error) {
	return "hello world " + s, nil
}
