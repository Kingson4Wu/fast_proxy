package outconfig

import (
	"github.com/Kingson4Wu/fast_proxy/common/config"
)

type Config interface {
	ForwardAddress() string
	config.Config
}

var configuration Config

func Get() Config {
	return configuration
}

func Read(c Config) {
	configuration = c
}
