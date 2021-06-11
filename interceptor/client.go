package interceptor

import "context"

type Invoker func(ctx context.Context, req, rep interface{}) error

type ClientInvoker func(ctx context.Context, req, rep interface{}, invoker Invoker) error

func ClientInvoke(ctx context.Context, req, rep interface{}, invoke Invoker, clientInvoker []ClientInvoker) error {
	if len(clientInvoker) == 0 {
		return invoke(ctx, req, rep)
	}
	return clientInvoker[0](ctx, req, rep, getInvoke(0, clientInvoker, invoke))
}

func getInvoke(cur int, clientInvoker []ClientInvoker, invoke Invoker) Invoker {
	if cur == len(clientInvoker)-1 {
		return invoke
	}
	return getInvoke(cur, clientInvoker, invoke)
}
