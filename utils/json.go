package utils

import (
	"io"
	"encoding/json"
)

func EncodeJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func DecodeJSON(r io.Reader, v interface{}) error {
	data, err := io.ReadAll(r)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}