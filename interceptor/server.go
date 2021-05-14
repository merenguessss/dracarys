package interceptor

import "context"

type Handler func(ctx context.Context, req interface{}) (rep interface{}, err error)

type ServerHandler func(ctx context.Context, req interface{}, handler Handler) (rep interface{}, err error)

func ServerHandle(ctx context.Context, serverHandles []ServerHandler, handler Handler, req interface{}) (rep interface{}, err error) {
	if len(serverHandles) == 0 {
		return handler(ctx, req)
	}
	return serverHandles[0](ctx, req, getHandler(0, serverHandles, handler))
}

func getHandler(cur int, serverHandles []ServerHandler, handler Handler) Handler {
	if cur == len(serverHandles)-1 {
		return handler
	}
	return getHandler(cur+1, serverHandles, handler)
}
