package config

import (
	"fmt"
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

func (c *yamlConfig) ServerContextPath() string {
	return viper.GetString("application.contextPath")
}

func (c *yamlConfig) GetCallTypeConfigByName(name string, uri string) int {
	return viper.GetInt(fmt.Sprintf("serviceCallTypeConfig.%s.%s.callType", name, uri))
}

func (c *yamlConfig) ContainsCallPrivilege(name string, uri string) bool {
	return c.GetCallTypeConfigByName(name, uri) > 0
}
