package protocol

import (
	"errors"

	"github.com/merenguessss/dracarys-go/codec"
	"github.com/merenguessss/dracarys-go/protocol"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterClientCodec(codec.Proto, &pbClientCodec{})
	RegisterServerCodec(codec.Proto, &pbServerCodec{})
}

var clientMap = make(map[uint8]codec.Codec)
var serverMap = make(map[uint8]codec.Codec)

func RegisterClientCodec(t uint8, c codec.Codec) {
	if clientMap == nil {
		clientMap = make(map[uint8]codec.Codec)
	}
	clientMap[t] = c
}

func GetClientCodec(t uint8) codec.Codec {
	if v, ok := clientMap[t]; ok {
		return v
	}
	return &pbClientCodec{}
}

func RegisterServerCodec(t uint8, c codec.Codec) {
	if serverMap == nil {
		serverMap = make(map[uint8]codec.Codec)
	}
	serverMap[t] = c
}

func GetServerCodec(t uint8) codec.Codec {
	if v, ok := serverMap[t]; ok {
		return v
	}
	return &pbServerCodec{}
}

type pbClientCodec struct {
}

var serializerType string

func (p *pbClientCodec) Encode(msg codec.Msg, bytes []byte) ([]byte, error) {
	metadata := map[string][]byte{
		serializerType: []byte(msg.SerializerType()),
	}
	req := &protocol.Request{
		RequestId:   uint32(msg.RequestID()),
		ServiceName: msg.ServerServiceName(),
		MethodName:  msg.RPCMethodName(),
		Metadata:    metadata,
		Payload:     bytes,
	}
	return proto.Marshal(req)
}

func (p *pbClientCodec) Decode(msg codec.Msg, bytes []byte) ([]byte, error) {
	rep := &protocol.Response{}
	err := proto.Unmarshal(bytes, rep)
	if err != nil {
		return nil, err
	}

	if rep.RetCode != 0 {
		return nil, errors.New(rep.RetMsg)
	}
	msg.WithRequestID(uint8(rep.RequestId))
	msg.WithSerializerType(string(rep.Metadata[serializerType]))
	return rep.Payload, nil
}

type pbServerCodec struct {
}

func (p *pbServerCodec) Encode(msg codec.Msg, b []byte) ([]byte, error) {
	rep := &protocol.Response{
		RetCode:   msg.RetCode(),
		RetMsg:    msg.RetMsg(),
		RequestId: uint32(msg.RequestID()),
		Metadata: map[string][]byte{
			serializerType: []byte(msg.SerializerType()),
		},
		Payload: b,
	}
	return proto.Marshal(rep)
}

func (p *pbServerCodec) Decode(msg codec.Msg, b []byte) ([]byte, error) {
	req := &protocol.Request{}
	err := proto.Unmarshal(b, req)
	if err != nil {
		return nil, err
	}
	msg.WithSerializerType(string(req.Metadata[serializerType]))
	msg.WithServerServiceName(req.ServiceName)
	msg.WithRPCMethodName(req.MethodName)
	msg.WithRequestID(uint8(req.RequestId))
	return req.Payload, nil
}
