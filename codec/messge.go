package codec

type MessageBuilder interface {
	Default() Msg
	WithSerializerType(string) Msg
	WithCompressType(uint8) Msg
	WithPackageType(uint8) Msg
	WithMsgType(uint8) Msg
	WithReqType(uint8) Msg
	WithClientServiceName(string) Msg
	WithServerServiceName(string) Msg
	WithRPCMethodName(string) Msg
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

	WithSerializerType(string)
	WithCompressType(uint8)
	WithPackageType(uint8)
	WithMsgType(uint8)
	WithReqType(uint8)
	WithClientServiceName(string)
	WithServerServiceName(string)
	WithRPCMethodName(string)
}

type msg struct {
	serializerType    string
	compressType      uint8
	packageType       uint8
	msgType           uint8
	reqType           uint8
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
