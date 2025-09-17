package config

import (
	"encoding/json"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/apollo"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/storage"
)

var config *apollo.ApolloConfig

func ApolloConfig() *apollo.ApolloConfig {
	return config
}

var log logger.Logger

// ConfigReader is a tiny adapter to make Apollo access testable.
// Only methods used in this file are included.
type ConfigReader interface {
    GetString(name string) string
    GetIntValue(name string, def int) int
    GetBoolValue(name string, def bool) bool
}

// cfgReader is consulted by parseAllConfig and getters. Defaults to the real Apollo-backed reader.
var cfgReader ConfigReader

type apReader struct{}

func (apReader) GetString(name string) string            { return config.GetString(name) }
func (apReader) GetIntValue(name string, def int) int    { return config.GetIntValue(name, def) }
func (apReader) GetBoolValue(name string, def bool) bool { return config.GetBoolValue(name, def) }

func LoadApolloConfig(appId string, namespace string, cluster string, address string, logger logger.Logger) Config {

	log = logger
	applicationConfig, applicationClient := initApollo(appId, namespace, cluster, address)

	c2 := &customChangeListener{}
	(*applicationClient).AddChangeListener(c2)

	config = &apollo.ApolloConfig{
		Log:            log,
		Client:         *applicationClient,
		CacheInterface: *applicationConfig,
		Namespace:      namespace,
	}

	// set default reader to real apollo-backed reader
	cfgReader = apReader{}

	parseAllConfig()

	return &apolloConfig{}
}

func initApollo(appId string, namespace string, cluster string, address string) (*agcache.CacheInterface, *agollo.Client) {

	client := apollo.InitApolloClient(appId, namespace, cluster, address)

	nsConfig := (*client).GetConfigCache(namespace)

	if nsConfig == nil {
		log.Error("", "appId", appId, "namespace", namespace, "cluster", cluster, "address", address)
		panic("init apollo failure")
	}

	log.Info("init apollo config successfully", "appId", appId, "namespace", namespace, "cluster", cluster, "address", address)

	return &nsConfig, client

}

type apolloConfig struct {
}

func (c *apolloConfig) ServerPort() int {
	return cfgReader.GetIntValue("application.port", 0)
}

func (c *apolloConfig) ServerName() string {
	return cfgReader.GetString("application.name")
}

func (c *apolloConfig) ServiceRpcHeaderName() string {
	return getConfigString("rpc.serviceHeaderName")
}

func (c *apolloConfig) GetServiceConfig(serviceName string) *ServiceConfig {

	sc, ok := serviceConfigMap[serviceName]
	if ok {
		return &sc
	}

	return nil
}

func (c *apolloConfig) GetSignKey(serviceConfig *ServiceConfig) string {
	v, ok := signKeyMap[serviceConfig.SignKeyName]
	if ok {
		return v
	}

	return ""
}

func (c *apolloConfig) GetSignKeyByName(name string) string {
	v, ok := signKeyMap[name]
	if ok {
		return v
	}

	return ""
}

func (c *apolloConfig) GetEncryptKey(serviceConfig *ServiceConfig) string {
	v, ok := encryptKeyMap[serviceConfig.EncryptKeyName]
	if ok {
		return v
	}

	return ""
}

func (c *apolloConfig) GetEncryptKeyByName(name string) string {
	v, ok := encryptKeyMap[name]
	if ok {
		return v
	}

	return ""
}

func (c *apolloConfig) GetTimeoutConfigByName(name string, uri string) int {
	v, ok := serviceTimeoutConfigMap[name]
	if ok {
		v, ok := v[uri]
		if ok {
			return v.Timeout
		}
	}
	return 0
}

func (c *apolloConfig) HttpClientMaxIdleConns() int {
	return config.GetIntValue("httpClient.MaxIdleConns", 5000)
}

func (c *apolloConfig) HttpClientMaxIdleConnsPerHost() int {
	return cfgReader.GetIntValue("httpClient.MaxIdleConnsPerHost", 3000)
}

func (c *apolloConfig) FastHttpEnable() bool {
	return cfgReader.GetBoolValue("fastHttp.Enable", false)
}

var (
	serviceConfigMap        map[string]ServiceConfig
	encryptKeyMap           EncryptKeyConfig
	signKeyMap              SignKeyConfig
	serviceTimeoutConfigMap map[string]map[string]ServiceTimeoutConfig
)

func setServiceConfig(m map[string]ServiceConfig) {
	serviceConfigMap = m
}
func setEncryptKey(m map[string]string) {
	encryptKeyMap = m
}
func setSignKey(m map[string]string) {
	signKeyMap = m
}
func setServiceTimeoutConfig(m map[string]map[string]ServiceTimeoutConfig) {
	serviceTimeoutConfigMap = m
}

const (
	serviceConfigName        = "proxy.service.config"
	signKeyConfigName        = "proxy.sign.key.config"
	encryptKeyConfigName     = "proxy.encrypt.key.config"
	serviceTimeoutConfigName = "proxy.service.timeout.config"
)

func parseAllConfig() {
	serviceConfigJson := getConfigString(serviceConfigName)
	if !ParseApolloConfig(serviceConfigJson, setServiceConfig, serviceConfigName) {

		log.Error("Unmarshal failure ... ", "name", serviceConfigName)
		panic(fmt.Sprintf("Unmarshal %s failure", serviceConfigName))
	}

	encryptKeyConfigJson := getConfigString(encryptKeyConfigName)

	if !ParseApolloConfig(encryptKeyConfigJson, setEncryptKey, encryptKeyConfigName) {

		log.Error("Unmarshal failure ... ", "name", encryptKeyConfigName)
		panic(fmt.Sprintf("Unmarshal %s failure", encryptKeyConfigName))
	}

	signKeyConfigJson := getConfigString(signKeyConfigName)

	if !ParseApolloConfig(signKeyConfigJson, setSignKey, signKeyConfigName) {

		log.Error("Unmarshal failure ... ", "name", signKeyConfigName)
		panic(fmt.Sprintf("Unmarshal %s failure", signKeyConfigName))
	}

	serviceTimeoutConfig := getConfigString(serviceTimeoutConfigName)

	if !ParseApolloConfig(serviceTimeoutConfig, setServiceTimeoutConfig, serviceTimeoutConfigName) {

		log.Error("Unmarshal failure ... ", "name", serviceTimeoutConfigName)
		panic(fmt.Sprintf("Unmarshal %s failure", serviceTimeoutConfigName))
	}
}

func getConfigString(configName string) string {
	c := cfgReader.GetString(configName)
	if c == "" {
		log.Error("get failure ... ", "name", configName)
		panic(fmt.Sprintf("get %s failure", configName))
	}
	return c
}

func ParseApolloConfig[K string, V any](configJson string, p func(map[K]V), configName string) bool {
	b := []byte(configJson)
	m := make(map[K]V)

	err := json.Unmarshal(b, &m)
	if err != nil {
		log.Error("get failure ... ", "name", configName)
		return false
	}
	p(m)

	return true
}

type customChangeListener struct {
	//wg sync.WaitGroup
}

func (c *customChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	//write your code here
	//fmt.Println(changeEvent.Changes)
	for key, value := range changeEvent.Changes {

		log.Info("change key ", key, value)

		switch key {
		case serviceConfigName:
			ParseApolloConfig(value.NewValue.(string), setServiceConfig, serviceConfigName)
		case signKeyConfigName:
			ParseApolloConfig(value.NewValue.(string), setSignKey, signKeyConfigName)
		case encryptKeyConfigName:
			ParseApolloConfig(value.NewValue.(string), setEncryptKey, encryptKeyConfigName)
		case serviceTimeoutConfigName:
			ParseApolloConfig(value.NewValue.(string), setServiceTimeoutConfig, serviceTimeoutConfigName)
		default:
		}
	}
	for _, listener := range listeners {
		listener(changeEvent)
	}

}

func (c *customChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
}

var listeners []func(*storage.ChangeEvent)

func RegisterApolloListener(f func(*storage.ChangeEvent)) {
	listeners = append(listeners, f)
}
