package im

import "encoding/json"

type Encoder interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

type JsonEncoder struct{}

func NewJsonEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

func (e *JsonEncoder) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (e *JsonEncoder) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
