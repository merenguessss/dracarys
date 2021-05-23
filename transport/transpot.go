package transport

import (
	"context"
	"net"
)

type ServerTransport interface {
	ListenAndServe(context.Context, ...ServerOption) error
}

type ClientTransport interface {
	Send(context.Context, []byte, ...ClientOption) ([]byte, error)
}

type Network string

const (
	TCP Network = "tcp"
	UDP Network = "udp"
)

// sendTCPMsg 发送TCP信息包的函数，通过循环写入连接中.
func sendTCPMsg(ctx context.Context, conn net.Conn, b []byte) (err error) {
	sendNum := 0
	addNum := 0
	for sendNum < len(b) {
		addNum, err = conn.Write(b[sendNum:])
		if err != nil {
			return err
		}
		sendNum += addNum

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}
