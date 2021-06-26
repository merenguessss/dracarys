## 服务中心

`Dracarys`默认使用`consul`作为服务中心。如果不进行配置，将默认使用`consul`。

### 配置

服务中心被`Dracarys`用作插件使用，所以服务中心的配置在插件配置中。

服务中心的配置如下：

**name**

服务中心名称，用于使用特定服务中心，默认为`consul`。

**address**

服务中心地址，服务中心所在的`ip`地址，默认为`127.0.0.1:8500`。

**enable-heartbeat**

心跳检测开关配置，默认为`true`，开启心跳检测。

**scheme**

协议名，服务中心中使用的协议，默认为`http`。

**heartbeat**

心跳配置：
    
- `host`: 本机`host`,用于服务中心向本机发送心跳检测，默认为`server`监听的`ip`。
- `port`：本机心跳程序使用`port`，默认为`8001`。
- `timeout`：心跳超时时间，默认为`5s`。
- `interval`：心跳检测间隔时间，默认为`5s`。
- `deregister-critical-service-after`：超时注销服务时间，默认为`20s`。

#### 配置方法

`Dracarys`有两种方式进行配置，代码配置和[全局配置](config.md)

**代码配置**

示例：
```go
// 先设置插件工厂配置
pluginOpts := []plugin.Option{
    plugin.WithSelectorName("your selector"),
    plugin.WithSelectorAddress("127.0.0.1:8888"),
}
// 将插件工厂配置放入客户端配置中
clientOpts := []client.Option{
    client.WithPluginFactoryOptions(pluginOpts),
}
client := dracarys.NewClient(clientOpts...)
```

**全局配置**

全局配置需要使用`yml`进行配置，详细参考[全局配置](config.md)。

### 自定义服务中心

`Dracarys`的所有插件都是可插拔的，服务中心也不例外。

如需使用自定义插件工厂，只需三个步骤。
1. 继承服务中心接口，实现方法。
2. 将服务中心注册到`Dracarys`中。
3. 插件工厂中放入自定义服务中心。

**1.继承接口**

```go
type Selector interface {
    // 注册客户端
	RegisterClient(name, address string) error
    // 注册服务
	RegisterService(name, address string) error
    // 注册心跳
	RegisterHeartbeat()
    // 通过name获取服务结点
	Select(string) (*ServiceNodes, error)
    // 加载配置
	LoadConfig(*Options) error
}
```

自定义服务中心需继承以上接口才能使用。

**2.注册插件**

示例：
```go
mySelector := NewSelector()
selector.Register("selectorName",mySelector)
```

即可将服务中心注册到`Dracarys`中。

**3.加入插件**

现在只需将插件配置到`Dracarys`即可。

配置方式见上文（[配置](#配置)）

**服务中心名称改为当前自定义服务中心即可。**