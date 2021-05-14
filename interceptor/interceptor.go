package interceptor

import "context"

type Invoker func(ctx context.Context, req, rep interface{}) error

type Interceptor func(ctx context.Context, req, rep interface{}, invoker Invoker) error

func Invoke(ctx context.Context, req, rep interface{}, invoke Invoker, interceptors []Interceptor) error {
	if len(interceptors) == 0 {
		return invoke(ctx, req, rep)
	}
	return interceptors[0](ctx, req, rep, getInvoke(0, interceptors, invoke))
}

func getInvoke(cur int, interceptors []Interceptor, invoke Invoker) Invoker {
	if cur == len(interceptors)-1 {
		return invoke
	}
	return getInvoke(cur, interceptors, invoke)
}
