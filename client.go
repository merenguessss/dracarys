package dracarys

import (
	"context"
	"time"

	"github.com/merenguessss/dracarys-go/client"
)

type Client struct {
	c    client.Client
	opts []client.Option
}

func NewClient(opts ...client.Option) *Client {
	return &Client{
		c:    client.New(),
		opts: opts,
	}
}

func (c *Client) Service(name string) {
	c.opts = append(c.opts, client.WithService(name))
}

func (c *Client) Call(methodName string, req interface{}) (interface{}, error) {
	c.opts = append(c.opts, client.WithMethod(methodName))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return c.c.Invoke(ctx, req, "", c.opts...)
}
