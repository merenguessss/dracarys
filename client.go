package dracarys

import (
	"context"
	"errors"
	"strings"

	"github.com/merenguessss/dracarys-go/client"
	"github.com/merenguessss/dracarys-go/config"
	"github.com/merenguessss/dracarys-go/log"
)

type Client struct {
	c    client.Client
	opts []client.Option
}

type Method func(...interface{}) (interface{}, error)

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
	c.opts = append(c.opts, client.WithService(name))
}

// CallWithReturnValue 将返回值传入进行调用,可以直接解析出返回值,而不需要再进行类型断言.
// rep 可以为任意形式参数,并不局限于结构体.
func (c *Client) CallWithReturnValue(methodName string, rep interface{}, req ...interface{}) error {
	c.opts = append(c.opts, client.WithMethod(methodName))
	return c.c.Invoke(context.Background(), req, rep, c.opts...)
}

// Call 通过MethodName定位请求到具体的req.
func (c *Client) Call(methodName string, req ...interface{}) (interface{}, error) {
	c.opts = append(c.opts, client.WithMethod(methodName))
	var rep interface{}
	if err := c.c.Invoke(context.Background(), req, rep, c.opts...); err != nil {
		return nil, err
	}
	return rep, nil
}

// Method 获取具体Method.
func (c *Client) Method(name string) Method {
	c.opts = append(c.opts, client.WithMethod(name))
	return func(req ...interface{}) (interface{}, error) {
		var rep interface{}
		if err := c.c.Invoke(context.Background(), req, rep, c.opts...); err != nil {
			return nil, err
		}
		return rep, nil
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
		var rep interface{}
		if err := c.c.Invoke(context.Background(), req, rep, c.opts...); err != nil {
			return nil, err
		}
		return rep, nil
	}, nil
}

// 解析serviceName和MethodName.
func (c *Client) parseServicePath(path string) (string, string, error) {
	index := strings.LastIndex(path, "/")
	if index == 0 || index == -1 || !strings.HasPrefix(path, "dracarys.service.") {
		return "", "", errors.New("invalid path")
	}
	return path[0:index], path[index+1:], nil
}
