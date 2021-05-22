package protocol

import (
	"testing"

	"github.com/merenguessss/Dracarys-go/codec"
)

func TestClientEncode(t *testing.T) {
	pbCCodec := GetClientCodec(codec.Proto)
	pbSCodec := GetServerCodec(codec.Proto)
	msg := codec.MsgBuilder.WithServerServiceName("server service").
		WithRequestID(uint8(123)).WithRPCMethodName("method").Build()

	testStr := []string{
		"ronaldo",
		"",
		"1231asd5",
	}
	for _, v := range testStr {
		b := []byte(v)
		req, err := pbCCodec.Encode(msg, b)
		if err != nil {
			t.Log("client encode err")
		}
		rep, err := pbSCodec.Decode(msg, req)
		if err != nil {
			t.Log("server decode error")
		}
		if string(rep) != v || msg.RequestID() != uint8(123) {
			t.Error("codec error")
		}
	}
}

func TestServerEncode(t *testing.T) {
	pbCCodec := GetClientCodec(codec.Proto)
	pbSCodec := GetServerCodec(codec.Proto)
	msg := codec.MsgBuilder.WithRetCode(uint32(0)).WithRetMsg("success").
		WithRequestID(uint8(123)).Build()

	testStr := []string{
		"ronaldo",
		"",
		"dasilkdjmsa,123",
	}
	for _, v := range testStr {
		b := []byte(v)
		req, err := pbSCodec.Encode(msg, b)
		if err != nil {
			t.Log("client encode err")
		}
		rep, err := pbCCodec.Decode(msg, req)
		if err != nil {
			t.Log("server decode error")
		}
		if string(rep) != v || msg.RequestID() != uint8(123) {
			t.Error("codec error")
		}
	}
}
