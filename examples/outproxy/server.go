package main

import (
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/outproxy"
)

func main() {

	c := &Config{}
	outproxy.NewServer(c)
}

type Config struct {
}

func (c Config) ForwardAddress() string {
	return "http://127.0.0.1:8384/inProxy"
}

func (c *Config) ServerPort() int {
	return 8084
}

func (c *Config) ServerName() string {
	return "curry"
}

func (c *Config) ServiceRpcHeaderName() string {
	return "C_service_name"
}

func (c *Config) GetServiceConfig(serviceName string) *config.ServiceConfig {

	return &config.ServiceConfig{
		EncryptKeyName: "encrypt.key.room.v2",
		SignKeyName:    "sign.key.room.v1",
		EncryptEnable:  true,
		SignEnable:     true,
		CompressEnable: true,
	}
}

func (c *Config) GetSignKey(serviceConfig *config.ServiceConfig) string {
	return "abcd"
}

func (c *Config) GetSignKeyByName(name string) string {
	return "abcd"
}

func (c *Config) GetEncryptKey(serviceConfig *config.ServiceConfig) string {
	return "ABCDABCDABCDABCD"
}

func (c *Config) GetEncryptKeyByName(name string) string {
	return "abcd"
}

func (c *Config) GetTimeoutConfigByName(name string, uri string) int {
	return 3000
}

func (c *Config) HttpClientMaxIdleConns() int {
	return 5000
}

func (c *Config) HttpClientMaxIdleConnsPerHost() int {
	return 3000
}
