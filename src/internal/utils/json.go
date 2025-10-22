package utils

import (
	"bytes"
	"encoding/json"
)

func JsonStrictUnmarshal[T any](data []byte) (T, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	var t T
	err := decoder.Decode(&t)
	return t, err
}
