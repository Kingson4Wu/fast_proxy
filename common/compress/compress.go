package compress

import (
	"github.com/Kingson4Wu/fast_proxy/common/compress/gzip"
	"github.com/Kingson4Wu/fast_proxy/common/compress/snappy"
)

var algorithms map[Algorithm]Compress

func init() {

	algorithms = make(map[Algorithm]Compress)
	algorithms[Snappy] = new(snappy.Snappy)
	algorithms[Gzip] = new(gzip.Gzip)

}

type Algorithm int32

const (
	Snappy Algorithm = iota
	Gzip
)

func use(algorithm Algorithm) Compress {
	if v, ok := algorithms[algorithm]; ok {
		return v
	}
	return nil
}

func Encode(data []byte, algorithm int32) ([]byte, error) {

	return use(Algorithm(algorithm)).Encode(data)
}

func Decode(data []byte, algorithm int32) ([]byte, error) {

	return use(Algorithm(algorithm)).Decode(data)
}

type Compress interface {
	Encode(data []byte) ([]byte, error)
	Decode(data []byte) ([]byte, error)
}
