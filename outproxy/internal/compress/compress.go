package compress

import (
	"github.com/Kingson4Wu/fast_proxy/common/snappy"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/logger"
)

var SnappyCompress *snappy.SnappyCompress

func init() {
	SnappyCompress = &snappy.SnappyCompress{
		Log: logger.GetLogger(),
	}
}
