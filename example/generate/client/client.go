package main

import (
	"context"

	"github.com/merenguessss/dracarys/client"
	pb "github.com/merenguessss/dracarys/example/generate/helloworld"
	"github.com/merenguessss/dracarys/log"
)

func main() {
	opts := []client.Option{
		client.WithAddr("127.0.0.1:8000"),
		client.WithSerializerType("proto"),
	}
	c := pb.NewGreeterClient(opts...)
	in := &pb.HelloRequest{
		Msg: "hello",
	}
	rep, err := c.SayHello(context.Background(), in)
	if err != nil {
		log.Fatal(rep)
	}
	log.Info(rep)
}
