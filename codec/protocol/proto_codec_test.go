package protocol

import (
	"fmt"
	"testing"

	"github.com/merenguessss/Dracarys-go/codec"
)

func TestClientEncode(t *testing.T) {
	pbCCodec := GetClientCodec(codec.Proto)
	pbSCodec := GetServerCodec(codec.Proto)
	msg := codec.MsgBuilder.WithServerServiceName("server service").
		WithRequestID(uint8(123)).WithRPCMethodName("method").Build()
	b := []byte("ronaldo")
	req, err := pbCCodec.Encode(msg, b)
	if err != nil {
		t.Log("client encode err")
	}
	t.Log(msg)
	rep, err := pbSCodec.Decode(msg, req)
	if err != nil {
		t.Log("server decode error")
	}
	fmt.Println(string(rep) == "ronaldo")
	t.Log(string(rep))
	t.Log(msg)
}
