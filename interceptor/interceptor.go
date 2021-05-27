package interceptor

import "context"

type Invoker func(ctx context.Context, req interface{}) (rep interface{}, err error)

type Interceptor func(ctx context.Context, req interface{}, invoker Invoker) (rep interface{}, err error)

func Invoke(ctx context.Context, req interface{}, invoke Invoker, interceptors []Interceptor) (interface{}, error) {
	if len(interceptors) == 0 {
		return invoke(ctx, req)
	}
	return interceptors[0](ctx, req, getInvoke(0, interceptors, invoke))
}

func getInvoke(cur int, interceptors []Interceptor, invoke Invoker) Invoker {
	if cur == len(interceptors)-1 {
		return invoke
	}
	return getInvoke(cur, interceptors, invoke)
}
