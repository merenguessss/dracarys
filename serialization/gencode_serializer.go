package serialization

func init() {
	Register(Gencode, &gencode_serializer{})
}

type gencode interface {
	Size() (s uint64)
	Marshal(buf []byte) ([]byte, error)
	Unmarshal(buf []byte) (uint64, error)
}

type gencode_serializer struct {
}

// Marshal gencode序列化方法.(注: 传入v必须为point类型).
func (gs *gencode_serializer) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.(gencode); ok {
		return data.Marshal(nil)
	}
	return nil, ErrorMissingGencodeMethod
}

func (gs *gencode_serializer) Unmarshal(b []byte, v interface{}) error {
	if data, ok := v.(gencode); ok {
		if _, err := data.Unmarshal(b); err != nil {
			return err
		}
		return nil
	}
	return ErrorMissingGencodeMethod
}
