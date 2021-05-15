package client

import (
	"context"
	"github.com/merenguessss/Dracarys-go/codec"
	"github.com/merenguessss/Dracarys-go/interceptor"
	"github.com/merenguessss/Dracarys-go/serialization"
)

type Client interface {
	Invoke(ctx context.Context, req, rep interface{}, path string, option ...Option) error
}

var DefaultClient = New()

func New() *defaultClient {
	return &defaultClient{}
}

type defaultClient struct {
	option *Options
}

func (c *defaultClient) Invoke(ctx context.Context, req, rep interface{}, path string,
	option ...Option) error {
	for _, op := range option {
		op(c.option)
	}
	return interceptor.Invoke(ctx, req, rep, c.invoke, c.option.beforeHandle)
}

func (c *defaultClient) invoke(ctx context.Context, req, rep interface{}) error {
	serializer := serialization.Get(c.option.serializerType)
	reqBuf, err := serializer.Marshal(req)
	if err != nil {
		return err
	}

	msg := codec.Msg{}
	coder := codec.Get(c.option.codecType)
	//reqBody,err := coder.Decode(msg,reqBuf)
	_, err = coder.Decode(msg, reqBuf)
	if err != nil {
		return err
	}

	//addr := c.findAddress()

	return nil
}

func (c *defaultClient) findAddress() string {
	return ""
}
