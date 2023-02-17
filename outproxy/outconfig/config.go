package outconfig

import (
	"github.com/Kingson4Wu/fast_proxy/common/config"
)

type Config interface {
	ForwardAddress() string
	config.Config
}

var Configuration Config

func Get() Config {
	return Configuration
}
