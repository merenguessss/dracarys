package serialization

import (
	"bytes"

	"github.com/vmihailenco/msgpack"
)

func init() {
	Register(Msgpack, &msgpackSerializer{})
}

type msgpackSerializer struct {
}

func (ms *msgpackSerializer) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	coder := msgpack.NewEncoder(&buf)
	err := coder.Encode(v)
	return buf.Bytes(), err
}

func (ms *msgpackSerializer) Unmarshal(b []byte, v interface{}) error {
	buf := bytes.NewBuffer(b)
	coder := msgpack.NewDecoder(buf)
	return coder.Decode(v)
}
