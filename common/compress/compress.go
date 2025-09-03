package compress

import (
	"github.com/Kingson4Wu/fast_proxy/common/compress/gzip"
	"github.com/Kingson4Wu/fast_proxy/common/compress/interfaces"
	"github.com/Kingson4Wu/fast_proxy/common/compress/snappy"
)

var algorithms map[Algorithm]interfaces.Compress

func init() {

	algorithms = make(map[Algorithm]interfaces.Compress)
	algorithms[Snappy] = new(snappy.Snappy)
	algorithms[Gzip] = new(gzip.Gzip)

}

type Algorithm int32

const (
	Snappy Algorithm = iota
	Gzip
)

func use(algorithm Algorithm) interfaces.Compress {
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
