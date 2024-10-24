package http

import (
	"encoding/json"
	"io"
)

func Decode[T any](input io.Reader) (*T, error) {
	v := new(T)
	err := json.NewDecoder(input).Decode(v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func Encode[T any](w io.Writer, v T) {
	_ = json.NewEncoder(w).Encode(v)
}
