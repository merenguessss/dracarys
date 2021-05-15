package codec

import "sync"

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
	return defaultCodec
}

var defaultCodec = NewDefaultCodec()

var NewDefaultCodec = func() Codec {
	return &codec{}
}

type codec struct{}

func (cc *codec) Encode(msg Msg, bytes []byte) ([]byte, error) {
	return nil, nil
}

func (cc *codec) Decode(msg Msg, bytes []byte) ([]byte, error) {
	return nil, nil
}
