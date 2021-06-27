## 负载均衡

`Dracarys`中提供了三种负载均衡的方式(后续还会继续更新🏆)

- 随机
- 轮询
- 加权轮询

### 配置

负载均衡器和服务中心一样，属于`Dracarys`的插件。所以其中的配置方式和服务中心一样，所以这里就不给出示例了。

负载均衡的配置参数如下：

**name**

负载均衡器的名称，用于指定负载均衡器，默认为`weight_poll`。

### 自定义负载均衡器

如需使用自定义负载均衡器，只需三个步骤：
1. 继承负载均衡接口，实现方法。
2. 将负载均衡注册到`Dracarys`中。
3. 插件工厂中放入自定义负载均衡器。

**1.继承接口**

```go
type Balancer interface {
	Get(*selector.ServiceNodes) *selector.Node
}
```

**2.注册负载均衡**

```go
myBalancer := NewBalancer()
balance.Register("my balancer",myBalancer)
```

**3.载入负载均衡**

现在只需将插件配置到`Dracarys`即可。

配置方式见上文（[配置](#配置)）


### 服务节点

服务节点具体信息参考[服务中心文档](selector.md)