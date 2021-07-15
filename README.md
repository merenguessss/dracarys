# Dracarys
**Dracarys**是一个轻量级💨、高性能🚀、开源🌟、可插拔🐲的RPC框架。

## 证书
`Dracarys`的源码允许用户在遵循[MIT 开源证书](LICENSE) 规则的前提下使用。
## 安装
在安装`Dracarys`之前，需要先部署Go环境(Go 1.16+)。

在Go module支持下，你只需要简单执行以下命令进行导入：
```
import "github.com/merenguessss/dracarys"
```
在你的代码中执行执行 go [build | run | test] 时会自动下载依赖项。

否则，请执行以下命令进行安装：
```
go get -u github.com/merenguessss/dracarys
```
## 快速开始
下面展示一个最简单的例子：
```
git clone https://github.com/merenguessss/dracarys.git
#server
cd example/reflection/server
go run main.go
#client
cd example/reflection/client
go run main.go
```
具体调用文档请参考[**`Dracarys`使用手册**](docs/quick_start.md)

## 功能特性
- [x] [服务发现](#服务发现)
- [x] [负载均衡](#负载均衡)
- [x] [日志系统](#日志系统)
- [x] [序列化](#序列化)
- [x] [全局配置](#全局配置)
- [x] 连接池
- [x] 可插拔式(所有功能都是可插拔的)
- [ ] 限流
- [ ] 熔断
- [ ] 链路追踪
- [ ] 后续功能持续开发中...

### 服务发现
`Dracarys`目前默认使用`consul`作为服务发现插件。

调用`consul`的api实现服务发现、服务注册、心跳检测等。

如需自定义服务中心，请参考[**`Dracarys`服务中心手册**](docs/selector.md)

（`Dracarys`自有的服务中心正在开发中😉）
### 负载均衡
`Dracarys`目前默认使用加权轮询的负载均衡方式。

支持的负载均衡器有：加权轮询、轮询、随机、一致性hash(在写了😬)

详细使用及自定义负载均衡器请参考[**`Dracarys`负载均衡手册**](docs/balancer.md)

可以通过client端配置插件配置进行使用：
```go
pluginOpts := []plugin.Option{
	WithBalancerName(balancerName),
}
opts := []client.Option{
	WithPluginFactoryOptions(pluginOpts),
}
client := dracarys.NewClient(opts...)
```
也可以通过[全局配置](#全局配置)进行配置。

如需自定义负载均衡器

### 日志系统
`Dracarys`支持多样化日志系统。

目前默认使用自带的日志系统。

如需自定义日志系统请参考[**`Dracarys`日志系统手册**](docs/log.md)

### 序列化
`Dracarys`支持多样化的序列化方式。

目前默认使用`json`进行代码反射调用序列化。如果使用代码生成，将会使用`protocol`序列化。

除此之外还支持[`gencode`](https://github.com/andyleap/gencode)（序列化速度更快）、`msgpack`作为反射序列化。

如需自定义序列化方式，请参考[**`Dracarys`序列化手册**](docs/serializer.md)

### 全局配置
`Dracarys`支持通过yml进行全局配置。

详细配置内容请参考[**`Dracarys`全局配置手册**](docs/config.md)

## 协议

协议格式如图：

![protocol](docs/pic/protocol.png)

### Header
`header`中包括以下几个部分：

![header](docs/pic/Header.png)

**Magic**

`Magic`为魔数，用于判断是否为`Dracarys`协议。魔数为常量0x12。

**Version**

`Version`为版本号，表示当前服务号。用于在某个服务更新后依旧需要提供旧版的服务所用。

**MsgType**

`MsgType`为消息类型，分为普通消息和心跳消息。

**ReqType**

`ReqType`为请求类型，分为发送并接收、只发送不接收、客户端流请求、服务端流请求、双向流请求。

**CompressType**

`CompressType`为压缩类型，用于接收方识别该消息的压缩类型(待扩展)。

**StreamID**

`StreamID`在使用流请求时使用（待扩展）。

**PackageType**

`PackageType`为包头类型，即编解码协议中Message内容时使用。

**Length**

`Length`为协议总长度。

**Reserved**

`Reserved`为保留位，待后续使用。

### Message

Message为协议报文头部，主要分为请求报文和响应报文。

**请求报文头**
```
message Request{
  uint32            request_id    // 一次请求的唯一ID
  string            service_name  // 服务名
  string            method_name   // 方法名
  map<string,bytes> metadata      // 传递的数据
  bytes             payload       // body
}
```
**响应报文头**
```
message Response{
  uint32            ret_code    // 错误码   成功请求为 0
  uint32            request_id  // 一次请求的唯一ID
  string            ret_msg     // 错误消息
  map<string,bytes> metadata    // 传递的数据
  bytes             payload     // body
}
```
### Body
序列化后的请求或响应内容。

## 性能

这里取`grpc`和`Dracarys`进行性能对比。

对两个框架进行三次测试，进行100w次请求，对三次测试结果取最大值。

`Dracarys`:
```
[DRACARYS]2021/06/27 12:10:51 proc.go:225: [INFO] total req  : 1000000
[DRACARYS]2021/06/27 12:10:51 proc.go:225: [INFO] success num: 1000000
[DRACARYS]2021/06/27 12:10:51 proc.go:225: [INFO] fail  num  : 0
[DRACARYS]2021/06/27 12:10:51 proc.go:225: [INFO] total time : 25733
[DRACARYS]2021/06/27 12:10:51 proc.go:225: [INFO] tps        : 38860
```

在同样的环境下，对`grpc`进行测试

`grpc`:
```
[DRACARYS]2021/06/27 12:17:04 proc.go:225: [INFO] total req  : 1000000
[DRACARYS]2021/06/27 12:17:04 proc.go:225: [INFO] success num: 1000000
[DRACARYS]2021/06/27 12:17:04 proc.go:225: [INFO] fail  num  : 0
[DRACARYS]2021/06/27 12:17:04 proc.go:225: [INFO] total time : 57381
[DRACARYS]2021/06/27 12:17:04 proc.go:225: [INFO] tps        : 17427
```