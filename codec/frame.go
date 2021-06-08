package codec

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

var (
	ErrorFrameHeaderRead  = errors.New("frame header read error")
	ErrorMagicNumber      = errors.New("magic number error ")
	ErrorRPCVersion       = errors.New("version error")
	ErrorPayloadOutLength = errors.New("payload beyond max length")
	ErrorPayloadRead      = errors.New("read payload error")
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
	counter    int
	conn       net.Conn
	readBuffer []byte
}

// ReadFrame 读取帧信息，先判断帧头信息，再对魔数、版本号进行判断，最后读取完整payload.
func (f *framer) ReadFrame() ([]byte, error) {
	var n int
	var err error
	frameHeader := make([]byte, FrameHeaderLen)
	n, err = io.ReadFull(f.conn, frameHeader)
	if n != FrameHeaderLen || err != nil {
		if n == 0 {
			return nil, io.EOF
		}
		return nil, ErrorFrameHeaderRead
	}

	if magic := frameHeader[0]; magic != Magic {
		return nil, ErrorMagicNumber
	}
	if version := frameHeader[1]; Version != version {
		return nil, ErrorRPCVersion
	}

	length := binary.BigEndian.Uint32(frameHeader[7:11]) - FrameHeaderLen
	if length > MaxPayloadLength {
		return nil, ErrorPayloadOutLength
	}
	if length > uint32(len(f.readBuffer)) && f.counter < 12 {
		f.readBuffer = make([]byte, len(f.readBuffer)*2)
		f.counter++
	}

	if n, err = io.ReadFull(f.conn, f.readBuffer[:length]); uint32(n) != length || err == io.EOF {
		return nil, ErrorPayloadRead
	}
	return append(frameHeader, f.readBuffer[:length]...), nil
}

const DefaultBufferLength = 1024
const MaxPayloadLength = 4 * 1024 * 1024

// FrameHeader 帧头.
type FrameHeader struct {
	// Magic 魔数.
	Magic uint8
	// Version 版本号.
	Version uint8
	// MsgType 消息类型.
	MsgType uint8
	// ReqType 请求类型.
	ReqType uint8
	// CompressType 压缩类型.
	CompressType uint8
	// StreamID 暂时没用，后续扩展.
	StreamID uint8
	// PackageType 用于包头压缩类型
	PackageType uint8
	// Length 帧长度.
	Length uint32
	// Reserved 保留位.
	Reserved uint32
}

const Magic = 0x12
const Version = 0
const FrameHeaderLen = 15

// 协议protocol打包类型.
const (
	Proto = iota
	Thrift
	Arvo
)

// 压缩类型.
const (
	NoneCompress = iota
)

// 消息类型.
const (
	GeneralMsg = 0x0
	HeartMsg   = 0x1
)

// 请求类型.
const (
	SendAndRecv = iota
	SendOnly
	LongConn
	StreamTrans
)

// PackageTypeToString 解析uint8类型的packageType帧头
func PackageTypeToString(t uint8) string {
	switch t {
	case 0:
		return "proto"
	case 1:
		return "thrift"
	default:
		return "arvo"
	}
}

// StrToPackageType string类型转到uint8放到帧头压缩.
func StrToPackageType(t string) uint8 {
	switch t {
	case "proto", "Proto":
		return 0
	case "thrift", "Thrift":
		return 1
	case "Arvo", "arvo":
		return 2
	default:
		return 3
	}
}
