package main

import (
	"fmt"

	"github.com/merenguessss/dracarys"
	"github.com/merenguessss/dracarys/client"
)

func main() {
	opts := []client.Option{
		client.WithAddr("localhost:8000"),
		client.WithNetWork("tcp"),
		client.WithSerializerType("json"),
		client.WithCodecType("proto"),
	}
	c := dracarys.NewClient(opts...)
	c.Service("Hello")
	res, err := c.Call("World", "1111")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
