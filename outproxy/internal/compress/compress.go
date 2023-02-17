package compress

import (
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/Kingson4Wu/fast_proxy/common/snappy"
)

var SnappyCompress *snappy.SnappyCompress

func init() {
	SnappyCompress = &snappy.SnappyCompress{
		Log: logger.GetLogger(),
	}
}
