package server

import (
	"context"

	"github.com/merenguessss/dracarys-go/codec"
	"github.com/merenguessss/dracarys-go/codec/protocol"
)

type dispatcher struct {
	serviceMap map[string]Service
}

var DefaultDispatcher = NewDispatcher()

var NewDispatcher = func() *dispatcher {
	return &dispatcher{
		serviceMap: make(map[string]Service),
	}
}

// RegisterService 将service注册到dispatcher中.
func (d *dispatcher) RegisterService(name string, s Service) {
	if d.serviceMap == nil {
		d.serviceMap = make(map[string]Service)
	}
	d.serviceMap[name] = s
}

// Handle 收到请求之后定位到具体的service进行调用.
func (d *dispatcher) Handle(ctx context.Context, b []byte) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	msg := codec.NewMsgBuilder().Build()
	var rep []byte

	coder := codec.DefaultCodec
	reqBody, err := coder.Decode(msg, b)
	if err != nil {
		return nil, err
	}

	protocolCoder := protocol.GetServerCodec(msg.PackageType())
	reqBuf, err := protocolCoder.Decode(msg, reqBody)
	if err != nil {
		return nil, err
	}

	if s, ok := d.serviceMap[msg.ServerServiceName()]; ok {
		rep, err = s.handle(msg, reqBuf)
		if err != nil {
			return nil, err
		}
	}

	// 重新选择包名解析器，因为服务端的解析器可能不同.
	protocolCoder = protocol.GetServerCodec(msg.PackageType())
	repBuf, err := protocolCoder.Encode(msg, rep)
	if err != nil {
		return nil, err
	}

	repBody, err := coder.Encode(msg, repBuf)
	if err != nil {
		return nil, err
	}

	return repBody, nil
}
