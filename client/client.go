// client client端,主要包含client端调用的核心内容.
package client

import (
	"context"
	"errors"

	"github.com/merenguessss/dracarys-go/codec"
	"github.com/merenguessss/dracarys-go/codec/protocol"
	"github.com/merenguessss/dracarys-go/interceptor"
	"github.com/merenguessss/dracarys-go/serialization"
	"github.com/merenguessss/dracarys-go/transport"
)

// ErrorServiceNotExist 服务不存在错误.
var ErrorServiceNotExist = errors.New("service not exist")

// Client 调用接口.
type Client interface {
	Invoke(ctx context.Context, req, rep interface{}, option ...Option) error
}

func New(o *Options) *defaultClient {
	return &defaultClient{
		option: o,
	}
}

type defaultClient struct {
	option *Options
}

// Invoke 动态代理函数,先加载client配置,再动态代理执行client的invoke函数.
func (c *defaultClient) Invoke(ctx context.Context, req, rep interface{},
	option ...Option) error {
	for _, op := range option {
		op(c.option)
	}
	return interceptor.ClientInvoke(ctx, req, rep, c.invoke, c.option.beforeHandle)
}

// invoke client动态代理函数,client的核心函数.
// 包括client所做的所有操作.
func (c *defaultClient) invoke(ctx context.Context, req, rep interface{}) error {
	msg := c.getMsg()

	serializer := serialization.Get(msg.SerializerType())
	reqBuf, err := serializer.Marshal(req)
	if err != nil {
		return err
	}

	protocolCoder := protocol.GetClientCodec(msg.PackageType())
	reqBuf, err = protocolCoder.Encode(msg, reqBuf)
	if err != nil {
		return err
	}

	coder := codec.Get(c.option.CodecType)
	reqBody, err := coder.Encode(msg, reqBuf)
	if err != nil {
		return err
	}

	addr, err := c.findAddress()
	if err != nil {
		return err
	}

	transportOption := []transport.ClientOption{
		transport.WithAddr(addr),
		transport.WithDisableConnPool(c.option.DisableConnPool),
		transport.WithEnableMultiplexed(c.option.EnableMultiplexed),
		transport.WithNetWork(transport.Network(c.option.NetWork)),
	}
	clientTransport := c.NewClientTransport()
	repBody, err := clientTransport.Send(ctx, reqBody, transportOption...)
	if err != nil {
		return err
	}

	repBuf, err := coder.Decode(msg, repBody)
	if err != nil {
		return err
	}

	protocolCoder = protocol.GetClientCodec(msg.PackageType())
	repBuf, err = protocolCoder.Decode(msg, repBuf)
	if err != nil {
		return err
	}

	serializer = serialization.Get(msg.SerializerType())
	err = serializer.Unmarshal(repBuf, rep)
	if err != nil {
		return err
	}
	return nil
}

func (c *defaultClient) NewClientTransport() transport.ClientTransport {
	return transport.GetClientTransport("default")
}

// findAddress client端通过serviceName查询地址.
// 其中应该包含服务发现->路由策略->负载均衡, 最终得到具体结点地址.
func (c *defaultClient) findAddress() (string, error) {
	if c.option.Addr != "" {
		return c.option.Addr, nil
	}
	slt := c.option.PluginFactory.GetSelector()

	if err := slt.RegisterClient(c.option.ClientName, c.option.Addr); err != nil {
		return "", err
	}

	nodes, err := slt.Select(c.option.serviceName)
	if err != nil {
		return "", err
	}

	if nodes.Length <= 0 {
		return "", ErrorServiceNotExist
	}

	balancer := c.option.PluginFactory.GetBalancer()
	node := balancer.Get(nodes)
	return node.Value, nil
}

// 通过client中的配置生成Msg.
func (c *defaultClient) getMsg() codec.Msg {
	mb := codec.NewMsgBuilder()
	return mb.WithCompressType(c.option.CompressType).
		WithSerializerType(c.option.SerializerType).
		WithPackageType(c.option.CompressType).
		WithServerServiceName(c.option.serviceName).
		WithRPCMethodName(c.option.methodName).Build()
}
