package codec

type MessageBuilder interface {
	Default() Msg
	WithSerializerType(string) MessageBuilder
	WithCompressType(uint8) MessageBuilder
	WithPackageType(uint8) MessageBuilder
	WithMsgType(uint8) MessageBuilder
	WithReqType(uint8) MessageBuilder
	WithClientServiceName(string) MessageBuilder
	WithServerServiceName(string) MessageBuilder
	WithRPCMethodName(string) MessageBuilder
	WithRetCode(uint32) MessageBuilder
	WithRetMsg(string) MessageBuilder
	Build() Msg
}

var MsgBuilder = &messageBuilder{}

type messageBuilder struct {
	msg Msg
}

func (mb *messageBuilder) Default() Msg {
	mb.msg = &msg{
		serializerType: "json",
		compressType:   NoneCompress,
		packageType:    Proto,
		msgType:        GeneralMsg,
		reqType:        SendAndRecv,
	}
	return mb.msg
}

func (mb *messageBuilder) WithSerializerType(st string) *messageBuilder {
	mb.msg.WithSerializerType(st)
	return mb
}
func (mb *messageBuilder) WithCompressType(ct uint8) *messageBuilder {
	mb.msg.WithCompressType(ct)
	return mb
}
func (mb *messageBuilder) WithPackageType(ft uint8) *messageBuilder {
	mb.msg.WithPackageType(ft)
	return mb
}
func (mb *messageBuilder) WithMsgType(mt uint8) *messageBuilder {
	mb.msg.WithMsgType(mt)
	return mb
}
func (mb *messageBuilder) WithReqType(rt uint8) *messageBuilder {
	mb.msg.WithReqType(rt)
	return mb
}
func (mb *messageBuilder) WithClientServiceName(csn string) *messageBuilder {
	mb.msg.WithClientServiceName(csn)
	return mb
}
func (mb *messageBuilder) WithServerServiceName(ssn string) *messageBuilder {
	mb.msg.WithServerServiceName(ssn)
	return mb
}
func (mb *messageBuilder) WithRPCMethodName(mn string) *messageBuilder {
	mb.msg.WithRPCMethodName(mn)
	return mb
}

func (mb *messageBuilder) WithRequestID(id uint8) *messageBuilder {
	mb.msg.WithRequestID(id)
	return mb
}

func (mb *messageBuilder) WithRetCode(code uint32) *messageBuilder {
	mb.msg.WithRetCode(code)
	return mb
}

func (mb *messageBuilder) WithRetMsg(msg string) *messageBuilder {
	mb.msg.WithRetMsg(msg)
	return mb
}

// TODO msg
type Msg interface {
	SerializerType() string
	CompressType() uint8
	PackageType() uint8
	MsgType() uint8
	ReqType() uint8
	ClientServiceName() string
	ServerServiceName() string
	RPCMethodName() string
	RequestID() uint8
	RetCode() uint32
	RetMsg() string

	WithSerializerType(string)
	WithCompressType(uint8)
	WithPackageType(uint8)
	WithMsgType(uint8)
	WithReqType(uint8)
	WithClientServiceName(string)
	WithServerServiceName(string)
	WithRPCMethodName(string)
	WithRequestID(uint8)
	WithRetCode(uint32)
	WithRetMsg(string)
}

type msg struct {
	serializerType    string
	compressType      uint8
	packageType       uint8
	msgType           uint8
	reqType           uint8
	requestID         uint8
	retCode           uint32
	retMsg            string
	clientServiceName string
	serverServiceName string
	rpcMethodName     string
}

func (m *msg) SerializerType() string {
	return m.serializerType
}

func (m *msg) CompressType() uint8 {
	return m.compressType
}

func (m *msg) PackageType() uint8 {
	return m.packageType
}

func (m *msg) MsgType() uint8 {
	return m.msgType
}

func (m *msg) ReqType() uint8 {
	return m.reqType
}

func (m *msg) RequestID() uint8 {
	return m.requestID
}

func (m *msg) RetCode() uint32 {
	return m.retCode
}

func (m *msg) RetMsg() string {
	return m.retMsg
}

func (m *msg) ClientServiceName() string {
	return m.clientServiceName
}

func (m *msg) ServerServiceName() string {
	return m.serverServiceName
}

func (m *msg) RPCMethodName() string {
	return m.rpcMethodName
}

func (m *msg) WithSerializerType(s string) {
	m.serializerType = s
}

func (m *msg) WithCompressType(u uint8) {
	m.compressType = u
}

func (m *msg) WithPackageType(u uint8) {
	m.packageType = u
}

func (m *msg) WithMsgType(u uint8) {
	m.msgType = u
}

func (m *msg) WithReqType(u uint8) {
	m.reqType = u
}

func (m *msg) WithClientServiceName(s string) {
	m.clientServiceName = s
}

func (m *msg) WithServerServiceName(s string) {
	m.serverServiceName = s
}

func (m *msg) WithRPCMethodName(s string) {
	m.rpcMethodName = s
}

func (m *msg) WithRequestID(u uint8) {
	m.requestID = u
}

func (m *msg) WithRetCode(code uint32) {
	m.retCode = code
}

func (m *msg) WithRetMsg(rm string) {
	m.retMsg = rm
}
