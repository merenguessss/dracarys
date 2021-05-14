package codec

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
