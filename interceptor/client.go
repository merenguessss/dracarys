package interceptor

//type Invoker func(ctx context.Context,req,rep interface{}) error
//
//type ClientInvokes func(ctx context.Context,req,rep interface{},invoker Invoker) error
//
//func ClientInvoke(ctx context.Context,req,rep interface{},invoke Invoker,clientInvokes []ClientInvokes) error{
//	if len(clientInvokes) == 0{
//		return invoke(ctx,req,rep)
//	}
//	return clientInvokes[0](ctx,req,rep,getInvoke(0,clientInvokes,invoke))
//}
//
//func getInvoke(cur int,clientInvokes []ClientInvokes,invoke Invoker)Invoker{
//	if cur == len(clientInvokes)-1{
//		return invoke
//	}
//	return getInvoke(cur,clientInvokes,invoke)
//}
