package outconfig

import (
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/spf13/viper"
)

func LoadYamlConfig(configBytes []byte) Config {

	return &yamlConfig{
		config.LoadYamlConfig(configBytes),
	}
}

type yamlConfig struct {
	config.Config
}

func (c *yamlConfig) ForwardAddress() string {
	return viper.GetString("proxy.forwardAddress")
}
