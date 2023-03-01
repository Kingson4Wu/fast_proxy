package snappy

import (
	"github.com/golang/snappy"
)

type Snappy struct {
}

func (s *Snappy) Encode(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

func (s *Snappy) Decode(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
