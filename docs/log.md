## 日志系统

`Dracarys`提供自带的日志系统

### 使用

日志支持以下接口中支持的方法：
```go
type Log interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}
```

代码中直接调用即可。
```go
log.Info("info")
log.Error("error")
```

### 配置

日志系统只支持全局配置，见[全局配置文档](config.md)。

配置参数如下：

**path**

日志系统输出的位置，默认为空，为空即输出到命令行。

**level**

日志级别，日志只会打印等级高于`level`的日志。

（日志等级按接口方法，从上而下。）

### 自定义日志系统

自定义一个日志系统很简单，只需两个步骤：

1. 继承日志系统接口
2. 指定自定义日志系统

示例:
```go
// 假定myLogger结构体继承了日志系统的接口
myLogger := NewLogger()
// 直接指定默认日志系统即可
log.DefaultLogger = myLogger
```