package transport

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/merenguessss/dracarys/codec"
	"github.com/merenguessss/dracarys/pool/conn_pool"
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
var DefaultClientTransport = NewClientDefault()

// NewDefault 默认ClientTransport的创建函数.
var NewClientDefault = func() ClientTransport {
	return &defaultClientTransport{
		clientOptions: &ClientOptions{
			pool: conn_pool.DefaultConnPool,
		},
	}
}

// defaultClientTransport 一个默认的ClientTransport.
type defaultClientTransport struct {
	clientOptions *ClientOptions
}

// Send 默认客户端发送消息到服务端.
func (ct *defaultClientTransport) Send(ctx context.Context, req []byte, options ...ClientOption) ([]byte, error) {
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
	return nil, errors.New("network not support")
}

// multiplexed 客户端多路复用连接实现.
func (ct *defaultClientTransport) multiplexed(ctx context.Context, req []byte) ([]byte, error) {
	return nil, nil
}

// sendTCP 发送TCP消息.
func (ct *defaultClientTransport) sendTCP(ctx context.Context, req []byte) ([]byte, error) {
	var conn net.Conn
	var err error
	address := ct.clientOptions.Addr
	network := string(ct.clientOptions.Network)
	var timeout time.Duration
	var rep []byte

	t, ok := ctx.Deadline()
	if ok {
		timeout = time.Until(t)
	}
	if ct.clientOptions.DisableConnPool {
		conn, err = net.DialTimeout(network, address, timeout)
		if err != nil {
			return nil, errors.New("direct connect error :" + err.Error())
		}
	} else {
		conn, err = ct.clientOptions.pool.Get(ctx, network, address)
		if err != nil {
			return nil, errors.New("connection pool error :" + err.Error())
		}
	}
	defer conn.Close()

	// 发送消息
	err = sendTCPMsg(ctx, conn, req)
	if err != nil {
		return nil, err
	}

	framer := codec.DefaultFramerBuilder.New(conn)
	rep, err = framer.ReadFrame()
	if err != nil {
		return nil, err
	}
	// TODO 接收消息 Framer
	return rep, nil
}

// sendUDP 发送UDP消息.
func (ct *defaultClientTransport) sendUDP(ctx context.Context, req []byte) ([]byte, error) {
	return nil, nil
}
