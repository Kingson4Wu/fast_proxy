package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
)

func LoadYamlConfig(configBytes []byte) Config {

	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewBuffer(configBytes)); err != nil {
		//fmt.Errorf("wrap error: %w", err)
		panic(err.Error())
	}

	return &yamlConfig{}
}

type yamlConfig struct {
}

func (c *yamlConfig) ForwardAddress() string {
	return viper.GetString("proxy.forwardAddress")
}

func (c *yamlConfig) ServerPort() int {
	return viper.GetInt("application.port")
}

func (c *yamlConfig) ServiceRpcHeaderName() string {
	return viper.GetString("rpc.serviceHeaderName")
}

func (c *yamlConfig) GetServiceConfig(serviceName string) *ServiceConfig {

	return &ServiceConfig{
		EncryptKeyName: viper.GetString(fmt.Sprintf("serviceConfig.%s.encryptKeyName", serviceName)),
		SignKeyName:    viper.GetString(fmt.Sprintf("serviceConfig.%s.signKeyName", serviceName)),
		EncryptEnable:  viper.GetBool(fmt.Sprintf("serviceConfig.%s.encryptEnable", serviceName)),
		SignEnable:     viper.GetBool(fmt.Sprintf("serviceConfig.%s.signEnable", serviceName)),
		CompressEnable: viper.GetBool(fmt.Sprintf("serviceConfig.%s.compressEnable", serviceName)),
	}
}

func (c *yamlConfig) GetSignKey(serviceConfig *ServiceConfig) string {
	return viper.GetString(fmt.Sprintf("signKeyConfig.%s", serviceConfig.SignKeyName))
}

func (c *yamlConfig) GetEncryptKey(serviceConfig *ServiceConfig) string {
	return viper.GetString(fmt.Sprintf("encryptKeyConfig.%s", serviceConfig.EncryptKeyName))
}

func (c *yamlConfig) GetEncryptKeyByName(name string) string {
	return viper.GetString(fmt.Sprintf("encryptKeyConfig.%s", name))
}

func (c *yamlConfig) GetTimeoutConfigByName(name string, uri string) int {
	return viper.GetInt(fmt.Sprintf("serviceTimeoutConfig.%s.%s", name, uri))
}
