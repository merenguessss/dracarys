## 序列化

`Dracarys`支持多样化的序列化方式。

目前提供了四种序列化方式：

- `json`
- `proto`
- `gencode`
- `msgpack`

从性能上来看`gencode`>`proto`>`json`>`msgpack`

```
goos: windows
goarch: amd64
pkg: github.com/merenguessss/dracarys/serialization
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
BenchmarkJsonSerializer-8       1000000000               0.08141 ns/op         0 B/op          0 allocs/op
BenchmarkMsgPackSerializer-8    1000000000               0.1211 ns/op          0 B/op          0 allocs/op
BenchmarkGencodeSerializer-8    1000000000               0.01262 ns/op         0 B/op          0 allocs/op
BenchmarkProtoSerializer-8      1000000000               0.04083 ns/op         0 B/op          0 allocs/op
```

其中`gencode`>`proto`>`json`>`msgpack`

- `gencode`、`json`、`msgpack`用于反射调用。

- `proto`用于代码生成使用。

如果使用`gencode`也必须使用到代码生成，详细使用方式请参考[`gencode`文档](https://github.com/andyleap/gencode)

**`Dracarys`推荐用户使用`gencode`进行调用。**

`Dracarys`默认使用`json`进行序列化。

### 配置

`Dracarys`支持两种序列化配置方式，代码配置以及全局配置。

**代码配置**

示例：
```go
opts := []client.Option{
    client.WithSerializerType("gencode")
}
client := dracarys.NewClient(opts...)
```

**全局配置**

全局配置需要使用`yml`进行配置，详细参考[全局配置](config.md)。

### 自定义序列化

`Dracarys`的序列化同样支持插拔，定义一个自定义序列化只需三个步骤：
1. 继承序列化接口。
2. 将序列化器注册到`Dracarys`。
3. 修改配置使用自定义序列化。

**1. 继承序列化接口**

```go
type Serialization interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}
```
**2. 注册序列化**

```go
mySerializer := NewSerializer()
serialization.Register("mySerializer",mySerializer)
```

**3. 载入序列化器**

使用配置的方式，指定使用的序列化器名称为自定义序列化器名称，[配置方式](#配置)如上。