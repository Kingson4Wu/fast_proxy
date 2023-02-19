package compress

import (
	"github.com/Kingson4Wu/fast_proxy/common/compress/gzip"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
)

var c Compress

func init() {
	/*c = &snappy.Snappy{
		Log: logger.GetLogger(),
	}*/

	c = &gzip.Gzip{
		Log: logger.GetLogger(),
	}
}

func Encode(data []byte) ([]byte, error) {
	return c.Encode(data)
}

func Decode(data []byte) ([]byte, error) {
	return c.Decode(data)
}

type Compress interface {
	Encode(data []byte) ([]byte, error)
	Decode(data []byte) ([]byte, error)
}
