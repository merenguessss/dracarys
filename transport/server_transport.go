package transport

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/merenguessss/Dracarys-go/codec"
)

func init() {
	RegisterServerTransport("default", DefaultServerTransport)
}

var serverTransportMap = make(map[string]ServerTransport)

var DefaultServerTransport = NewServerDefault()

var NewServerDefault = func() ServerTransport {
	return &defaultServerTransport{
		serverOptions: &ServerOptions{},
	}
}

func GetServerTransport(name string) ServerTransport {
	if v, ok := serverTransportMap[name]; ok {
		return v
	}
	return DefaultServerTransport
}

func RegisterServerTransport(name string, st ServerTransport) {
	if serverTransportMap == nil {
		serverTransportMap = make(map[string]ServerTransport)
	}
	serverTransportMap[name] = st
}

type defaultServerTransport struct {
	serverOptions *ServerOptions
}

func (st *defaultServerTransport) ListenAndServe(ctx context.Context, so ...ServerOption) error {
	for _, o := range so {
		o(st.serverOptions)
	}

	switch st.serverOptions.network {
	case TCP:
		return st.serveTCP(ctx)
	case UDP:
		return st.serveUDP(ctx)
	default:
		return errors.New("network not support")
	}
}

func (st *defaultServerTransport) serveTCP(ctx context.Context) error {
	network := string(st.serverOptions.network)
	addr := st.serverOptions.address

	lis, err := net.Listen(network, addr)
	if err != nil {
		return nil
	}

	go func() {
		if err := st.serve(ctx, lis); err != nil {
			log.Println("serve error")
		}
	}()
	return nil
}

func (st *defaultServerTransport) serve(ctx context.Context, lis net.Listener) error {
	var temporary time.Duration
	tcpLis, ok := lis.(*net.TCPListener)
	if !ok {
		return errors.New("network not support")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := tcpLis.AcceptTCP()
		if err != nil {
			netError, ok := err.(net.Error)
			if !ok {
				return err
			}

			if netError.Temporary() {
				if temporary == 0 {
					temporary = 1 * time.Millisecond
				} else {
					temporary *= 2
				}
				if max := 1 * time.Second; temporary > max {
					temporary = max
				}
				time.Sleep(temporary)
				continue
			}
			return err
		}

		if err := conn.SetKeepAlive(true); err != nil {
			return err
		}

		if err := conn.SetKeepAlivePeriod(st.serverOptions.keepAlivePeriod); err != nil {
			return err
		}

		go func() {
			if err := st.handleConn(ctx, conn); err != nil {
				log.Print(err)
			}
		}()
	}
}

func (st *defaultServerTransport) handleConn(ctx context.Context, conn net.Conn) error {

	defer conn.Close()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		framer := codec.DefaultFramerBuilder.New(conn)
		req, err := framer.ReadFrame()
		if err != nil {
			return errors.New("read frame error  " + err.Error())
		}

	}
}

func (st *defaultServerTransport) serveUDP(ctx context.Context) error {
	return nil
}
