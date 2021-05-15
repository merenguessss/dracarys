package serialization

type Serialization interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

const (
	Proto   = "proto"
	Msgpack = "msgpack"
	Json    = "json"
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

func (s *pbSerialization) Marshal(interface{}) ([]byte, error) {
	return nil, nil
}

func (s *pbSerialization) Unmarshal([]byte, interface{}) error {
	return nil
}
