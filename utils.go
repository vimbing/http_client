package http

import (
	"bytes"
	"encoding/json"
	"math/rand"
)

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func marshalAndEncodeBody(body RequestJsonBody) (*bytes.Buffer, error) {
	marshaled, err := json.Marshal(body)

	if err != nil {
		return &bytes.Buffer{}, err
	}

	return bytes.NewBuffer(marshaled), nil
}
