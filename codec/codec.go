package codec

import "sync"

// Codec
type Codec interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
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

func GetCodeC(name string) Codec {
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

func (cc *codec) Encode([]byte) ([]byte, error) {
	return nil, nil
}

func (cc *codec) Decode([]byte) ([]byte, error) {
	return nil, nil
}
