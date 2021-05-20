package codec

import (
	"bytes"
	"encoding/binary"
	"sync"
)

// Codec
type Codec interface {
	Encode(Msg, []byte) ([]byte, error)
	Decode(Msg, []byte) ([]byte, error)
}

var (
	codecMap = make(map[string]Codec)
	lock     sync.RWMutex
)

func RegisterCodeC(name string, codec Codec) {
	lock.Lock()
	defer lock.Unlock()
	if codecMap == nil {
		codecMap = make(map[string]Codec)
	}
	codecMap[name] = codec
}

func Get(name string) Codec {
	if codec, ok := codecMap[name]; ok {
		return codec
	}
	return DefaultCodec
}

var DefaultCodec = NewDefaultCodec()

var NewDefaultCodec = func() Codec {
	return &codec{}
}

type codec struct{}

func (cc *codec) Encode(msg Msg, b []byte) ([]byte, error) {
	if msg.CompressType() != NoneCompress {
		// TODO compress
	}

	length := FrameHeaderLen + len(b)
	buffer := bytes.NewBuffer(make([]byte, length))

	frameHeader := &FrameHeader{
		Magic:        Magic,
		Version:      Version,
		MsgType:      msg.MsgType(),
		ReqType:      msg.ReqType(),
		CompressType: msg.CompressType(),
		PackageType:  msg.PackageType(),
		Length:       uint32(length),
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.Magic); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.MsgType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.ReqType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.CompressType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.StreamID); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.PackageType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frameHeader.Reserved); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, b); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (cc *codec) Decode(msg Msg, b []byte) ([]byte, error) {
	msg.WithMsgType(b[2])
	msg.WithReqType(b[3])
	msg.WithCompressType(b[4])
	msg.WithPackageType(b[6])
	if msg.CompressType() != NoneCompress {
		// TODO compress
	}
	return b[FrameHeaderLen:], nil
}
