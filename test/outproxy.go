package test

import (
	"net/http"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
)

func GetOutSC() *servicediscovery.ServiceCenter {
	sc := servicediscovery.Create().
		AddressFunc(func(serviceName string) *servicediscovery.Address {
			return &servicediscovery.Address{
				Ip:   "127.0.0.1",
				Port: 9988,
			}
		}).ClientNameFunc(func(req *http.Request) string {
		return req.Header.Get("C_ServiceName")
	}).RegisterFunc(func(name string, ip string, port int) chan bool {
		return make(chan bool)
	}).Build()

	return sc
}

func GetMockOutConfig() outconfig.Config {
	// This would return a mock configuration for outproxy tests
	// For now, we'll create a simple implementation
	return &mockOutConfig{}
}

type mockOutConfig struct{}

func (c *mockOutConfig) ServerName() string {
	return "out_proxy"
}

func (c *mockOutConfig) ServerPort() int {
    // use 0 to let OS choose a free port to avoid collisions in CI
    return 0
}

func (c *mockOutConfig) ServiceRpcHeaderName() string {
	return "C_ServiceName"
}

func (c *mockOutConfig) GetServiceConfig(serviceName string) *config.ServiceConfig {
	return &config.ServiceConfig{
		EncryptKeyName: "encrypt.key.room.v2",
		SignKeyName:    "sign.key.room.v1",
		EncryptEnable:  true,
		SignEnable:     true,
		CompressEnable: true,
	}
}

func (c *mockOutConfig) GetSignKey(serviceConfig *config.ServiceConfig) string {
	return "abcd"
}

func (c *mockOutConfig) GetSignKeyByName(name string) string {
	return "abcd"
}

func (c *mockOutConfig) GetEncryptKey(serviceConfig *config.ServiceConfig) string {
	return "ABCDABCDABCDABCD"
}

func (c *mockOutConfig) GetEncryptKeyByName(name string) string {
	return "ABCDABCDABCDABCD"
}

func (c *mockOutConfig) GetTimeoutConfigByName(name string, uri string) int {
	return 5000
}

func (c *mockOutConfig) HttpClientMaxIdleConns() int {
	return 100
}

func (c *mockOutConfig) HttpClientMaxIdleConnsPerHost() int {
	return 10
}

func (c *mockOutConfig) FastHttpEnable() bool {
    return true
}

func (c *mockOutConfig) ForwardAddress() string {
	return "http://127.0.0.1:9833/outProxy"
}
