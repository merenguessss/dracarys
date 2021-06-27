## 全局配置

`Dracarys`不仅可以支持代码配置，也可以通过`yml`的方式进行全局配置。(代码配置可以覆盖`yml`的全局配置)

### 配置方式

全局配置默认读取`classpath`目录下的`dracarys.yml`文件。

如需修改配置文件地址可以使用以下方式：

**`SetPath`函数**

在使用`Dracarys`之前调用`SetPath`函数更改配置文件位置。

注：

`SetPath`函数需在`Dracarys`加载配置之前执行才会有效。

`client`端加载配置时机：
```go
client := dracarys.NewClient()
```
`server`端加载配置时机：
```go
server := dracarys.NewServer()
```

**设置环境变量**

用户可自行在设置环境变量`DRACARYS_CONFIG`进行配置。

### 配置参数

完整示例如下：

```yml
client:
    // 客户端名称
    name: client
    // 客户端直接请求地址
    address: 127.0.0.1:8080
    // 序列化类型
    serializer: json
    // 协议编码类型
    codec: proto
    // 是否开启多路复用
    enable_multiplexed: false
    // 是否禁止线程池
    disable_connection_pool: false
    // 网络传输协议
    network: tcp
    // 压缩类型
    compress: none

server:
    // 服务名
    name: server
    // 监听端口号
    port: 8080
    // 网络传输协议
    network: tcp
    // tcp长连接的心跳间隔
    keep_alive_period: 200
    // 序列化类型
    serializer: json
    // 协议编码类型
    codec: proto
    // 压缩类型
    compress: none

plugin:
    selector:
        // 服务中心名
        name: consul
        // 服务中心地址
        address: 127.0.0.1:8500
        // 传输协议
        scheme: http
        // 是否开启心跳
        enable-heartbeat: true
        // 心跳配置
        heartbeat: 
            // 检测ip
            host: 127.0.0.1
            // 检测端口
            port: 8080
            // 超时时间
            timeout: 5s 
            // 心跳信息间隔
            interval: 5s
            // 服务失效时间
            deregister-critical-service-after: 20s
    balancer:
        // 负载均衡名
        name: weight_poll
log:
    // 日志存储地址
    path: classpath:/dracarys.log
    // 日志等级
    level: 2
```