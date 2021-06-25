package server

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/merenguessss/dracarys/interceptor"
	"github.com/merenguessss/dracarys/log"
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
		ServiceName: s.wrapServiceName(serviceName),
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

		methodFilter := func(ctx context.Context, srv interface{}, parse func(interface{}) error,
			beforeHandle []interceptor.ServerHandler) (interface{}, error) {
			in := make([]interface{}, 0)
			var params []reflect.Value
			numIn := method.Type.NumIn()
			isPtr := make([]bool, numIn)

			for j := 2; j < numIn; j++ {
				reqType := method.Type.In(j)
				isPtr[j] = reqType.Kind() == reflect.Ptr
				if isPtr[j] {
					reqType = reqType.Elem()
				}
				in = append(in, reflect.New(reqType).Interface())
			}

			inData := make([]interface{}, 0)
			onlyParam := len(in) == 1
			if !onlyParam {
				if err := parse(&in); err != nil {
					return nil, err
				}
				inData = append(inData, in...)
			} else {
				parseData := in[0]
				if err := parse(parseData); err != nil {
					return nil, err
				}
				inData = append(inData, parseData)
			}

			params = append(params, reflect.ValueOf(service), reflect.ValueOf(ctx))
			for j, v := range inData {
				value := reflect.ValueOf(v)
				if !isPtr[j+2] {
					// 如果不是指针类型需要先进行解引用.
					value = value.Elem()
				}
				params = append(params, value)
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
		log.Error("handler type error")
		return
	}

	ser := &service{
		srv:         srv,
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
		log.Error("has same name service")
	}

	s.ServiceMap[ser.serviceName] = ser
	DefaultDispatcher.serviceMap[ser.serviceName] = ser
}

func (s *Server) Serve() error {
	var wg sync.WaitGroup
	wg.Add(len(s.ServiceMap))
	for _, v := range s.ServiceMap {
		go func(src Service) {
			if err := src.Serve(s.Options); err != nil {
				log.Error(err)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	return nil
}

// wrapServiceName 包装service name, 形成规范.
func (s *Server) wrapServiceName(name string) string {
	prefix := "dracarys.service."
	res := name
	if strings.HasPrefix(name, prefix) {
		res = name[len(prefix):]
	}

	if s.Options.ServerName != "" {
		res = s.Options.ServerName + "." + res
	}
	return prefix + res
}

func (s *Server) PrintLogo() {
	fmt.Println(
		"________________________________________________  _________\n" +
			"___  __ \\__  __ \\__    |_  ____/__    |__  __ \\ \\/ /_  ___/\n" +
			"__  / / /_  /_/ /_  /| |  /    __  /| |_  /_/ /_  /_____ \\ \n" +
			"_  /_/ /_  _, _/_  ___ / /___  _  ___ |  _, _/_  / ____/ / \n" +
			"/_____/ /_/ |_| /_/  |_\\____/  /_/  |_/_/ |_| /_/  /____/")
}
