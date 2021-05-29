package main

import (
	"fmt"
	"github.com/merenguessss/dracarys-go"
	"github.com/merenguessss/dracarys-go/server"
	"time"
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

func (h *Hello) World(s string) (string, error) {
	return "hello world " + s, nil
}
