package serialization

import "encoding/json"

func init() {
	Register(Json, &jsonSerializer{})
}

type jsonSerializer struct{}

func (js *jsonSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (js *jsonSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
