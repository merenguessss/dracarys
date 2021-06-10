package dracarys

import (
	"context"
	"errors"
	"strings"

	"github.com/merenguessss/dracarys-go/client"
)

type Client struct {
	c    client.Client
	opts []client.Option
}

type Method func(...interface{}) (interface{}, error)

func NewClient(opts ...client.Option) *Client {
	return &Client{
		c:    client.New(),
		opts: opts,
	}
}

func (c *Client) Service(name string) {
	c.opts = append(c.opts, client.WithService(name))
}

// Call 通过MethodName定位请求到具体的req.
func (c *Client) Call(methodName string, req ...interface{}) (interface{}, error) {
	c.opts = append(c.opts, client.WithMethod(methodName))
	return c.c.Invoke(context.Background(), req, c.opts...)
}

// Method 获取具体Method.
func (c *Client) Method(name string) Method {
	c.opts = append(c.opts, client.WithMethod(name))
	return func(req ...interface{}) (interface{}, error) {
		return c.c.Invoke(context.Background(), req, c.opts...)
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
		return c.c.Invoke(context.Background(), req, c.opts...)
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
