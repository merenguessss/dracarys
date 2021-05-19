package codec

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

func init() {
	DefaultFramerBuilder = &defaultFramerBuilder{}
}

var DefaultFramerBuilder FramerBuilder

type FramerBuilder interface {
	New(conn net.Conn) Framer
}

type Framer interface {
	ReadFrame() ([]byte, error)
}

type defaultFramerBuilder struct{}

func (fb *defaultFramerBuilder) New(conn net.Conn) Framer {
	return &framer{
		conn:       conn,
		readBuffer: make([]byte, DefaultBufferLength),
	}
}

type framer struct {
	conn       net.Conn
	readBuffer []byte
}

func (f *framer) ReadFrame() ([]byte, error) {
	frameHeader := make([]byte, FrameHeaderLen)
	if n, err := io.ReadFull(f.conn, frameHeader); n != FrameHeaderLen || err == io.EOF {
		return nil, errors.New("frame header error")
	}
	if magic := frameHeader[0]; magic != Magic {
		return nil, errors.New("magic number error ")
	}
	if version := frameHeader[1]; Version != version {
		return nil, errors.New("version error")
	}

	length := binary.BigEndian.Uint32(frameHeader)
	if length > MaxPayloadLength {
		return nil, errors.New("payload beyond max length")
	}
	if length > DefaultBufferLength {
		f.readBuffer = make([]byte, length)
	}
	if n, err := io.ReadFull(f.conn, f.readBuffer); n != int(length) || err == io.EOF {
		return nil, errors.New("read payload error")
	}
	return f.readBuffer[:length], nil
}

const DefaultBufferLength = 1024
const MaxPayloadLength = 4 * 1024 * 1024

type FrameType uint8
type MsgType uint8
type ReqType uint8
type CompressType uint8

// FrameHeader 帧头.
type FrameHeader struct {
	// Magic 魔数.
	Magic uint8
	// Version 版本号.
	Version uint8
	// MsgType 消息类型.
	MsgType MsgType
	// ReqType 请求类型.
	ReqType ReqType
	// CompressType 压缩类型.
	CompressType CompressType
	// StreamID 暂时没用，后续扩展.
	StreamID uint8
	// FrameType 用于帧解析出Header.
	FrameType FrameType
	// Length 帧长度.
	Length uint32
	// Reserved 保留位.
	Reserved uint32
}

const Magic = 0x12
const Version = 0
const FrameHeaderLen = 15

const (
	Proto FrameType = iota
	Thrift
	arvo
)

const (
	NoneCompress CompressType = iota
)

const (
	GeneralMsg MsgType = 0x0
	HeartMsg   MsgType = 0x1
)

const (
	SendAndRecv ReqType = iota
	SendOnly
	LongConn
	StreamTrans
)
