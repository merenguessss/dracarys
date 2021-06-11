package dracarys

import (
	"context"
	"errors"
	"strings"

	"github.com/merenguessss/dracarys-go/client"
	"github.com/merenguessss/dracarys-go/config"
)

type Client struct {
	c    client.Client
	opts []client.Option
}

type Method func(...interface{}) (interface{}, error)

func NewClient(opts ...client.Option) *Client {
	o, err := config.GetClient()
	if err != nil {
		// todo log err
		o = &client.Options{}
	}
	return &Client{
		c:    client.New(o),
		opts: opts,
	}
}

func (c *Client) Service(name string) {
	c.opts = append(c.opts, client.WithService(name))
}

// CallStruct 请求直接返回结构体的方法,传入结构体.
func (c *Client) CallStruct(methodName string, rep interface{}, req ...interface{}) error {
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
