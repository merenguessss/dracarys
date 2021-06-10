package client

import (
	"context"

	"github.com/merenguessss/dracarys-go/codec"
	"github.com/merenguessss/dracarys-go/codec/protocol"
	"github.com/merenguessss/dracarys-go/interceptor"
	"github.com/merenguessss/dracarys-go/serialization"
	"github.com/merenguessss/dracarys-go/transport"
)

type Client interface {
	Invoke(ctx context.Context, req interface{}, path string, option ...Option) (interface{}, error)
}

func New() *defaultClient {
	return &defaultClient{
		option: &Options{},
	}
}

type defaultClient struct {
	option *Options
}

func (c *defaultClient) Invoke(ctx context.Context, req interface{}, path string,
	option ...Option) (interface{}, error) {
	for _, op := range option {
		op(c.option)
	}
	r := []interface{}{req}
	return interceptor.Invoke(ctx, r, c.invoke, c.option.beforeHandle)
}

func (c *defaultClient) invoke(ctx context.Context, req interface{}) (interface{}, error) {
	msg := codec.MsgBuilder.Default()
	serializer := serialization.Get(msg.SerializerType())
	reqBuf, err := serializer.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg.WithServerServiceName(c.option.ServiceName)
	msg.WithRPCMethodName(c.option.MethodName)
	protocolCoder := protocol.GetClientCodec(msg.PackageType())
	reqBuf, err = protocolCoder.Encode(msg, reqBuf)
	if err != nil {
		return nil, err
	}

	coder := codec.Get(c.option.codecType)
	reqBody, err := coder.Encode(msg, reqBuf)
	if err != nil {
		return nil, err
	}

	addr := c.findAddress()

	transportOption := []transport.ClientOption{
		transport.WithAddr(addr),
		transport.WithDisableConnPool(c.option.DisableConnPool),
		transport.WithEnableMultiplexed(c.option.EnableMultiplexed),
		transport.WithNetWork(transport.Network(c.option.NetWork)),
	}
	clientTransport := c.NewClientTransport()
	repBody, err := clientTransport.Send(ctx, reqBody, transportOption...)
	if err != nil {
		return nil, err
	}

	repBuf, err := coder.Decode(msg, repBody)
	if err != nil {
		return nil, err
	}

	protocolCoder = protocol.GetClientCodec(msg.PackageType())
	repBuf, err = protocolCoder.Decode(msg, repBuf)
	if err != nil {
		return nil, err
	}

	var rep interface{}
	serializer = serialization.Get(msg.SerializerType())
	err = serializer.Unmarshal(repBuf, &rep)
	if err != nil {
		return nil, err
	}
	return rep, nil
}

func (c *defaultClient) NewClientTransport() transport.ClientTransport {
	return transport.GetClientTransport("default")
}

func (c *defaultClient) findAddress() string {
	return c.option.Addr
}
