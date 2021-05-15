package transport

import (
	"context"
	"errors"
)

func init() {
	ClientRegister("default", DefaultClientTransport)
}

// ClientTransport 存放的map,可以通过name取出具体的ClientTransport.
var clientTransportMap = make(map[string]ClientTransport)

// ClientRegister ClientTransport的注册函数.
func ClientRegister(name string, clientTransport ClientTransport) {
	if clientTransportMap == nil {
		clientTransportMap = make(map[string]ClientTransport)
	}
	clientTransportMap[name] = clientTransport
}

// GetClientTransport 通过name获取具体的ClientTransport.
func GetClientTransport(name string) ClientTransport {
	if v, ok := clientTransportMap[name]; ok {
		return v
	}
	return DefaultClientTransport
}

// 默认 ClientTransport
var DefaultClientTransport = NewDefault()

// NewDefault 默认ClientTransport的创建函数.
var NewDefault = func() ClientTransport {
	return &defaultClientTransport{}
}

// defaultClientTransport 一个默认的ClientTransport.
type defaultClientTransport struct {
	clientOptions *ClientOptions
}

// Send 默认客户端发送消息到服务端.
func (ct *defaultClientTransport) Send(ctx context.Context, req []byte, options ...ClientOption) error {
	for _, v := range options {
		v(ct.clientOptions)
	}

	if ct.clientOptions.EnableMultiplexed {
		return ct.multiplexed(ctx, req)
	}

	if ct.clientOptions.Network == TCP {
		return ct.sendTCP(ctx, req)
	}

	if ct.clientOptions.Network == UDP {
		return ct.sendUDP(ctx, req)
	}
	return errors.New("network not support")
}

// multiplexed 客户端多路复用连接实现.
func (ct *defaultClientTransport) multiplexed(ctx context.Context, req []byte) error {
	return nil
}

// sendTCP 发送TCP消息.
func (ct *defaultClientTransport) sendTCP(ctx context.Context, req []byte) error {

	return nil
}

// sendUDP 发送UDP消息.
func (ct *defaultClientTransport) sendUDP(ctx context.Context, req []byte) error {
	return nil
}
