package main

import (
	"fmt"
	"time"

	"github.com/merenguessss/dracarys-go"
	"github.com/merenguessss/dracarys-go/server"
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
		fmt.Println(err)
	}

	if err := srv.Serve(); err != nil {
		fmt.Println(err)
	}
}

type Hello struct {
}

func (h *Hello) World(s, t string) (Res, error) {
	return Res{
		"hello world " + s + t,
	}, nil
}

type Res struct {
	S string
}
