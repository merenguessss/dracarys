package main

import (
	"context"
	"github.com/merenguessss/dracarys/log"

	"github.com/merenguessss/dracarys"
	pb "github.com/merenguessss/dracarys/example/generate/helloworld"
	"github.com/merenguessss/dracarys/server"
)

func main() {
	opts := []server.Option{
		server.WithAddress("127.0.0.1:8000"),
		server.WithSerializerType("proto"),
	}
	s := dracarys.NewServer(opts...)
	pb.RegisterGreeterServer(s, &service{})
	if err := s.Serve(); err != nil {
		log.Error(err)
	}
}

type service struct {
}

func (s *service) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Msg: r.Msg + "world",
	}, nil
}
