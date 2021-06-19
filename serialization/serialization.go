package serialization

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

type Serialization interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

const (
	Proto   = "proto"
	Msgpack = "msgpack"
	Json    = "json"
	Gencode = "gencode"
)

var serializationMap = make(map[string]Serialization)

func init() {
	serializationMap[Proto] = defaultSerialization
}

func Register(name string, serialization Serialization) {
	if serializationMap == nil {
		serializationMap = make(map[string]Serialization)
	}
	serializationMap[name] = serialization
}

func Get(name string) Serialization {
	if v, ok := serializationMap[name]; ok {
		return v
	}
	return defaultSerialization
}

var defaultSerialization = newSerialization()

var newSerialization = func() Serialization {
	return &pbSerialization{}
}

type pbSerialization struct {
}

var ErrorMissingProtoMessage = errors.New("missing proto message")
var ErrorMissingGencodeMethod = errors.New("missing gencode generate code")

func (s *pbSerialization) Marshal(data interface{}) ([]byte, error) {
	v, ok := data.(proto.Message)
	if !ok {
		return nil, ErrorMissingProtoMessage
	}
	return proto.Marshal(v)
}

func (s *pbSerialization) Unmarshal(b []byte, data interface{}) error {
	v, ok := data.(proto.Message)
	if !ok {
		return ErrorMissingProtoMessage
	}
	return proto.Unmarshal(b, v)
}
