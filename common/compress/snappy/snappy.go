package snappy

import (
	"github.com/Kingson4Wu/fast_proxy/common/compress/interfaces"
	"github.com/golang/snappy"
)

type Snappy struct {
}

var _ interfaces.Compress = &Snappy{}

func (s *Snappy) Encode(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

func (s *Snappy) Decode(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
