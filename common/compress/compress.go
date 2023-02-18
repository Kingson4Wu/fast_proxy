package compress

import (
	"github.com/Kingson4Wu/fast_proxy/common/compress/snappy"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
)

// TODO 支持其他压缩方式， 抽象接口
var c *snappy.Compress

func init() {
	c = &snappy.Compress{
		Log: logger.GetLogger(),
	}
}

func Encode(data []byte) (result []byte, ok bool) {
	return c.Encode(data)
}

func Decode(data []byte) (result []byte, erro error) {
	return c.Decode(data)
}
