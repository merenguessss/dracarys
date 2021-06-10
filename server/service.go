package server

import (
	"context"
	"errors"

	"github.com/merenguessss/dracarys-go/codec"
	"github.com/merenguessss/dracarys-go/interceptor"
	"github.com/merenguessss/dracarys-go/serialization"
	"github.com/merenguessss/dracarys-go/transport"
)

type Service interface {
	Register(string, FilterFunc)
	Serve(o *Options) error
	Close()
	handle(codec.Msg, []byte) ([]byte, error)
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

type FilterFunc func(ctx context.Context, parse func(interface{}) error,
	handlers []interceptor.Interceptor) (rep interface{}, err error)

type service struct {
	ctx         context.Context
	cancel      context.CancelFunc
	serviceName string
	handles     map[string]FilterFunc
	opt         *Options
}

func (s *service) Register(methodName string, method FilterFunc) {
	if s.handles == nil {
		s.handles = make(map[string]FilterFunc)
	}
	s.handles[methodName] = method
}

func (s *service) Serve(o *Options) error {
	s.opt = o

	tsOpt := []transport.ServerOption{
		transport.WithAddress(s.opt.address),
		transport.WithNetwork(s.opt.network),
		transport.WithKeepAlivePeriod(s.opt.keepAlivePeriod),
		transport.WithHandler(DefaultDispatcher),
	}
	st := transport.DefaultServerTransport

	s.ctx, s.cancel = context.WithCancel(context.Background())
	if err := st.ListenAndServe(s.ctx, tsOpt...); err != nil {
		return errors.New(s.serviceName + " service transport error " + err.Error())
	}

	<-s.ctx.Done()
	return nil
}

func (s *service) Close() {
	if s.cancel != nil {
		s.cancel()
	}
}

// handle 调用具体RPC函数.
func (s *service) handle(msg codec.Msg, reqBuf []byte) ([]byte, error) {
	serializer := serialization.Get(msg.SerializerType())

	parser := func(req interface{}) error {
		if err := serializer.Unmarshal(reqBuf, req); err != nil {
			return err
		}
		return nil
	}

	handle := s.handles[msg.RPCMethodName()]
	rep, err := handle(s.ctx, parser, s.opt.beforeHandle)
	if err != nil {
		return nil, err
	}

	// 刷新msg的内容
	s.updateMsg(msg)
	serializer = serialization.Get(msg.SerializerType())
	repBuf, err := serializer.Marshal(rep)
	if err != nil {
		return nil, err
	}
	return repBuf, nil
}

// updateMsg 刷新msg中的内容.
func (s *service) updateMsg(msg codec.Msg) {
	// todo msg.WithCompressType()
	msg.WithSerializerType(s.opt.serializerType)
	msg.WithPackageType(codec.StrToPackageType(s.opt.codecType))
	msg.WithReqType(codec.SendOnly)
}
