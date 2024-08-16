package snappy

import (
	"github.com/Kingson4Wu/fast_proxy/common/compress"
	"github.com/golang/snappy"
)

type Snappy struct {
}

var _ compress.Compress = &Snappy{}

func (s *Snappy) Encode(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

func (s *Snappy) Decode(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
