package dracarys

import (
	"context"
	"errors"
	"strings"

	"github.com/merenguessss/dracarys/client"
	"github.com/merenguessss/dracarys/config"
	"github.com/merenguessss/dracarys/log"
)

type Client struct {
	c    client.Client
	opts []client.Option
}

type Method func(...interface{}) (interface{}, error)

// NewClientOpts 创建client的配置.
func NewClientOpts(opts ...client.Option) *client.Options {
	o, err := config.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	for _, op := range opts {
		op(o)
	}

	if err = o.PluginFactory.Setup(o.PluginFactoryOptions...); err != nil {
		log.Fatal(err)
	}
	return o
}

// NewClientWithOpts 通过client的配置生成Client进行调用.
func NewClientWithOpts(opts *client.Options) *Client {
	return &Client{
		c: client.New(opts),
	}
}

func NewClient(opts ...client.Option) *Client {
	o, err := config.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	for _, op := range opts {
		op(o)
	}

	if err = o.PluginFactory.Setup(o.PluginFactoryOptions...); err != nil {
		log.Fatal(err)
	}
	return &Client{
		c: client.New(o),
	}
}

func (c *Client) Service(name string) {
	c.opts = append(c.opts, client.WithService(c.wrapServiceName(name)))
}

// CallWithReturnValue 将返回值传入进行调用,可以直接解析出返回值,而不需要再进行类型断言.
// rep 可以为任意形式参数,并不局限于结构体.
func (c *Client) CallWithReturnValue(methodName string, rep interface{}, req ...interface{}) error {
	c.opts = append(c.opts, client.WithMethod(methodName))
	return c.call(context.Background(), c.opts, rep, req...)
}

func (c *Client) Invoke(ctx context.Context, req, rep interface{}, option ...client.Option) error {
	c.opts = append(c.opts, option...)
	return c.c.Invoke(ctx, req, rep, c.opts...)
}

// Call 通过MethodName定位请求到具体的req.
func (c *Client) Call(methodName string, req ...interface{}) (interface{}, error) {
	c.opts = append(c.opts, client.WithMethod(methodName))
	rep := new(interface{})
	if err := c.call(context.Background(), c.opts, rep, req...); err != nil {
		return nil, err
	}
	return *rep, nil
}

// Method 获取具体Method.
func (c *Client) Method(name string) Method {
	c.opts = append(c.opts, client.WithMethod(name))
	return func(req ...interface{}) (interface{}, error) {
		rep := new(interface{})
		if err := c.call(context.Background(), c.opts, rep, req...); err != nil {
			return nil, err
		}
		return *rep, nil
	}
}

// ServiceAndMethod 通过service和method查询函数.
func (c *Client) ServiceAndMethod(name string) (Method, error) {
	serviceName, methodName, err := c.parseServicePath(name)
	if err != nil {
		return nil, err
	}
	c.opts = append(c.opts, client.WithService(serviceName), client.WithMethod(methodName))
	return func(req ...interface{}) (interface{}, error) {
		rep := new(interface{})
		if err := c.call(context.Background(), c.opts, rep, req...); err != nil {
			return nil, err
		}
		return *rep, nil
	}, nil
}

// call 调用具体invoke方法前先进行参数判断,用于区分单参数和多参数.
func (c *Client) call(ctx context.Context, opts []client.Option, rep interface{}, req ...interface{}) error {
	var r interface{}
	r = req
	if len(req) == 1 {
		r = req[0]
	}
	return c.c.Invoke(ctx, r, rep, opts...)
}

var serviceNamePrefix = "dracarys.service."

// 解析serviceName和MethodName.
func (c *Client) parseServicePath(path string) (string, string, error) {
	index := strings.LastIndex(path, "/")
	if index == 0 || index == -1 {
		return "", "", errors.New("invalid path")
	}

	serviceName := path[0:index]
	if !strings.HasPrefix(path, serviceNamePrefix) {
		serviceName = serviceNamePrefix + serviceName
	}
	return serviceName, path[index+1:], nil
}

// wrapServiceName 包装serviceName,使其规范化.
func (c *Client) wrapServiceName(name string) string {
	if !strings.HasPrefix(name, serviceNamePrefix) {
		return serviceNamePrefix + name
	}
	return name
}
