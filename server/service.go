package server

import (
	"context"

	"github.com/merenguessss/Dracarys-go/interceptor"
)

type Service interface {
	Register(serviceDesc interface{}, serviceImpl interface{}) error
	Serve() error
	Close() error
}

type ServiceDesc struct {
	ServiceName string
	HandlerType interface{}
	Methods     []*Method
}

type Method struct {
	Name string
	Func FilterFunc
}

type FilterFunc func(svr interface{}, ctx context.Context, parse func(interface{}) error, handlers []interceptor.ServerHandler) (rep interface{}, err error)

type service struct {
	ctx         context.Context
	serviceName string
	handle      map[string]FilterFunc
}

func (s *service) Register(serviceDesc interface{}, serviceImpl interface{}) error {
	return nil
}

func (s *service) Serve() error {
	return nil
}

func (s *service) Close() error {
	return nil
}
