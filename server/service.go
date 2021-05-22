package server

import (
	"context"

	"github.com/merenguessss/Dracarys-go/interceptor"
)

type Service interface {
	Register(methodName, method FilterFunc)
	Serve() error
	Close() error
}

type ServiceDesc struct {
	Svr         interface{}
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
	handles     map[string]FilterFunc
}

func (s *service) Register(methodName string, method FilterFunc) {
	if s.handles == nil {
		s.handles = make(map[string]FilterFunc)
	}
	s.handles[methodName] = method
}

func (s *service) Serve() error {

	return nil
}

func (s *service) Close() error {
	return nil
}
