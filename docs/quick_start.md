## 快速开始

下面介绍`Dracarys`的两种调用方式：

- [反射调用](#反射调用)
- [代码生成](#代码生成)

### 反射调用

`Dracarys`允许使用反射方式进行调用。用户可直接定义一个`service`对象，并直接为其编写可导出方法，之后直接在`server`中注册即可。

#### 示例
发起一个服务只需要三个步骤：
1. 定义一个服务结构体。
2. 通过一个server发布服务。
3. client端调用服务。

**1. 定义一个服务结构体**
```go
package hello

import "context"

type Hello struct {
}

func (h *Hello) World(ctx context.Context, s string) (string, error) {
	return "hello world " + s, nil
}
```
**注：注册服务方法的第一个入参必须为```context.Context```**

**2. 通过一个server发布服务**
```go
func main() {
	opts := []server.Option{
		server.WithAddress("127.0.0.1:8000"),
		server.WithNetWork("tcp"),
		server.WithKeepAlivePeriod(time.Second * 200),
		server.WithSerializerType("json"),
	}
	srv := dracarys.NewServer(opts...)
	err := srv.RegisterService("dracarys.service.Hello", &Hello{})
	if err != nil {
		fmt.Println(err)
	}

	if err := srv.Serve(); err != nil {
		fmt.Println(err)
	}
}
```
**3. 通过client调用服务**
```go
func main() {
	opts := []client.Option{
		client.WithAddr("localhost:8000"),
		client.WithNetWork("tcp"),
		client.WithSerializerType("json"),
		client.WithCodecType("proto"),
	}
	c := dracarys.NewClient(opts...)
	c.Service("Hello")
	res, err := c.Call("World", "1111")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
```
反射调用，目前支持三种序列化方式：

1. `json`
2. `msgpack`
3. `gencode`

如果想要使用某种序列化方式只需在配置中加入：
```go
SerializerType := "you want"
client.WithSerializerType(SerializerType)
server.WithSerializerType(SerializerType)
```
如需了解更多或者使用自定义序列化器请参考[`Dracarys`序列化](serializer.md)

### 代码生成

`Dracarys`允许使用`protobuf`的代码生成器生成本地IDL文件进行调用，**此时只允许使用`proto`进行序列化**。

#### 示例

通过代码生成一个服务只需要六个步骤。

1. 编写proto文件。
2. 安装[`protoc-gen-go-dracarys`代码生成工具](https://github.com/merenguessss/protoc-gen-go-dracarys)（需要先安装好protoc工具）。
3. 生成代码。
4. 定义service结构体继承rpc方法。
5. 通过server发布一个服务。
6. 通过client进行调用。

**1. 编写proto文件**
```protobuf
syntax = "proto3";

package helloworld;
option go_package="github.com/merenguessss/dracarys/example/generate/helloworld";

service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
    string msg = 1;
}

message HelloReply {
    string msg = 1;
}
```

**2. 安装[`protoc-gen-go-dracarys`代码生成工具](https://github.com/merenguessss/protoc-gen-go-dracarys)**

使用以下命令进行安装：
```
go get github.com/merenguessss/protoc-gen-go-dracarys
go install github.com/merenguessss/protoc-gen-go-dracarys
```

执行以上命令后，会在`$GOPATH/bin`目录下生成可执行程序`protoc-gen-go-dracarys`

**3. 生成代码**

执行命令：
```
protoc --go-dracarys_out=. helloworld.proto
```
即可生成文件helloworld.pb.proto

**4.定义service结构体继承rpc方法**

在生成`helloworld.pb.proto`文件后，会生成一个接口对应`proto`文件中的`service`。需要定义结构体继承该接口才能使用。

```go
type service struct {
}

func (s *service) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Msg: r.Msg + "world",
	}, nil
}
```

**5. 通过server发布一个服务**

```go
func main() {
	opts := []server.Option{
		server.WithAddress("127.0.0.1:8000"),
		server.WithSerializerType("proto"),
	}
	s := dracarys.NewServer(opts...)
	pb.RegisterGreeterServer(s, &service{})
	if err := s.Serve(); err != nil {
		log.Error(err)
	}
}
```

**6. 通过client进行调用**

```go
func main() {
	opts := []client.Option{
		client.WithAddr("127.0.0.1:8000"),
		client.WithSerializerType("proto"),
	}
	c := pb.NewGreeterClient(opts...)
	in := &pb.HelloRequest{
		Msg: "hello",
	}
	rep, err := c.SayHello(context.Background(), in)
	if err != nil {
		log.Fatal(rep)
	}
	log.Info(rep)
}
```

## 名字规范

`Dracarys`对一个服务的命名规范为：

 dracarys.service.myservice

 例如定义一个服务helloworld，其服务名为dracarys.service.helloworld。

## 调用方式

**以下均为反射调用时可自行调用的函数。**

`Dracarys`支持多种客户端调用方式。

### 调用函数

1. 可以直接使用与`grpc`相似的`Invoke`函数进行调用:

```go
func Invoke(ctx context.Context, methodName string, req, rep interface{}, option ...client.Option) error
```

**如果使用该函数，必须在调用之前使用`Service`确定`rpc`函数。**

例：
```go
c := dracarys.NewClient()
c.Service("dracarys.service.helloworld")
if err := c.Invoke(ctx, "Method", req, rep, opts...);err != nil{
    return err
}
```

2. `CallWithReturnValue`函数

```go
func CallWithReturnValue(methodName string, rep interface{}, req ...interface{}) error
```

**如果使用该函数进行调用，需要先使用Service函数确定服务。**

例：
```go
c := dracarys.NewClient()
c.Service("dracarys.service.helloworld")
if err := c.CallWithReturnValue("Hello", rep, req...);err!=nil{
    return err
}
```

值得一提的是，方法支持**多参数**，可以使用连续的多个参数

3. `Call`函数

```go
func Call(methodName string, req ...interface{}) (interface{}, error)
```

**如果使用该函数进行调用，需要先使用Service函数确定服务。**

例：
```go
c := dracarys.NewClient()
c.Service("dracarys.service.helloworld")
rep, err := c.Call("Hello", req...)
if err != nil{
    return err
}
```

4. `Method`调用

`Dracarys`支持生成自定义方法：
```go
type Method func(...interface{}) (interface{}, error)
```

该方法支持多参数,单返回值。

4.1. `Method`函数

```go
func Method(name string) Method
```

**如果使用该函数进行调用，需要先使用Service函数确定服务。**

例：
```go
c := dracarys.NewClient()
c.Service("dracarys.service.helloworld")
Hello := c.Method("Hello")
rep, err := Hello(req)
if err != nil{
	return err
}
```

4.2. `ServiceAndMethod`函数

```go
func ServiceAndMethod(name string) (Method, error)
```

**此函数的参数`name`为服务名和函数名的组合**

写法格式为  **服务名/函数名**

例：
```go
c := dracarys.NewClient()
Hello, err := c.ServiceAndMethod("dracarys.service.helloworld/Hello")
if err != nil{
	return err
}
rep, err := Hello(req)
if err != nil{
	return err
}
```

### 注意问题

`CallWithReturnValue`和`Invoke`函数将`rpc`远程方法的返回值作为入参接收结果。而其余方法都是直接返回interface{}类型。此时会出现**类型断言**。

如果`rpc`函数的返回值是结构体，返回的interface{}将会被解析为`map`类型，将**不能进行类型断言为所需结构体。如果返回值为结构体，建议使用`CallWithReturnValue`和`Invoke`函数。**