package server

import (
	"context"
	"reflect"

	"github.com/merenguessss/dracarys-go/interceptor"
)

type Server struct {
	Options    *Options
	ServiceMap map[string]Service
}

type emptyHandlerType interface{}

func (s *Server) RegisterService(serviceName string, service interface{}, options ...Option) error {
	sd := &ServiceDesc{
		Svr:         service,
		HandlerType: (*emptyHandlerType)(nil),
		ServiceName: serviceName,
	}

	methods, err := s.getServiceMethod(service)
	if err != nil {
		return err
	}

	sd.Methods = methods
	s.Register(sd, service, options...)
	return nil
}

func (s *Server) getServiceMethod(service interface{}) ([]*Method, error) {
	srvType := reflect.TypeOf(service)
	n := srvType.NumMethod()
	methods := make([]*Method, n)

	for i := 0; i < n; i++ {
		method := srvType.Method(i)

		methodFilter := func(ctx context.Context, parse func(interface{}) error,
			beforeHandle []interceptor.ServerHandler) (interface{}, error) {
			in := make([]interface{}, 0)
			var params []reflect.Value

			if err := parse(&in); err != nil {
				return nil, err
			}

			params = append(params, reflect.ValueOf(service))
			for _, v := range in {
				params = append(params, reflect.ValueOf(v))
			}

			handler := func(ctx context.Context, reqBody interface{}) (interface{}, error) {
				value := method.Func.Call(params)
				// todo 多返回值
				return value[0].Interface(), nil
			}

			return interceptor.ServerHandle(ctx, beforeHandle, handler, in)
		}

		methods[i] = &Method{
			Name: method.Name,
			Func: methodFilter,
		}
	}
	return methods, nil
}

func (s *Server) Register(srvDesc *ServiceDesc, srv interface{}, opts ...Option) {
	if srvDesc == nil || srv == nil {
		return
	}
	handlerType := reflect.TypeOf(srvDesc.HandlerType).Elem()
	srvType := reflect.TypeOf(srv)
	if !srvType.Implements(handlerType) {
		// log
		return
	}

	ser := &service{
		serviceName: srvDesc.ServiceName,
		handles:     make(map[string]FilterFunc),
	}
	for _, v := range srvDesc.Methods {
		ser.handles[v.Name] = v.Func
	}
	for _, o := range opts {
		o(ser.opt)
	}

	if _, ok := s.ServiceMap[ser.serviceName]; ok {
		// log has same service
	}

	s.ServiceMap[ser.serviceName] = ser
	DefaultDispatcher.serviceMap[ser.serviceName] = ser
}

func (s *Server) Serve() error {
	for _, v := range s.ServiceMap {
		if err := v.Serve(s.Options); err != nil {
			return err
		}
	}
	return nil
}
