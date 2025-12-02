package im

import (
	"reflect"
	"testing"
)

type Payload struct {
	Key1 string   `json:"key1"`
	Key2 int      `json:"key2"`
	Key3 []string `json:"key3"`
}

func TestEncoder(t *testing.T) {
	e := NewJsonEncoder()
	original := Payload{"value1", 42, []string{"a", "b", "c"}}

	encoded, err := e.Encode(original)
	if err != nil {
		t.Fatalf("Encoding failed: %v", err)
	}

	var decoded Payload
	err = e.Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Fatalf("Decoded struct mismatch. Got: %+v, Want: %+v", decoded, original)
	}
	t.Logf("Original: %+v", original)
	t.Logf("Encoded: %s", string(encoded))
	t.Logf("Decoded: %+v", decoded)
}
