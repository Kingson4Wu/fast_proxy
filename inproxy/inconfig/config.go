package inconfig

import "github.com/Kingson4Wu/fast_proxy/common/config"

type ServiceCallConfig struct {
	CallType int `json:"callType"`
	Timeout  int `json:"timeout"`
	Limit    int `json:"limit"`
}

type Config interface {
	ServerContextPath() string
	GetCallTypeConfigByName(name string, uri string) int
	ContainsCallPrivilege(name string, uri string) bool
	config.Config
}

var Configuration Config

func Get() Config {
	return Configuration
}
