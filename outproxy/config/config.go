package config

type ServiceConfig struct {
	EncryptKeyName string `json:"encrypt.key.name"`
	SignKeyName    string `json:"sign.key.name"`
	EncryptEnable  bool   `json:"encrypt.enable"`
	SignEnable     bool   `json:"sign.enable"`
	CompressEnable bool   `json:"compress.enable"`
}

type EncryptKeyConfig map[string]string
type SignKeyConfig map[string]string

type ServiceTimeoutConfig struct {
	Timeout int `json:"timeout"`
}

type Config interface {
	ForwardAddress() string
	ServerPort() int
	ServiceRpcHeaderName() string
	GetServiceConfig(serviceName string) *ServiceConfig
	GetSignKey(serviceConfig *ServiceConfig) string
	GetEncryptKey(serviceConfig *ServiceConfig) string
	GetEncryptKeyByName(name string) string
	GetTimeoutConfigByName(name string, uri string) int
}

var Configuration Config

func Get() Config {
	return Configuration
}
