package config

import (
	"encoding/json"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/apollo"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/logger"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/storage"
	"go.uber.org/zap"
)

var config *apollo.ApolloConfig

func LoadApolloConfig(appId string, namespace string, cluster string, address string) Config {

	applicationConfig, applicationClient := initApollo(appId, namespace, cluster, address)

	c2 := &CustomChangeListener{}
	(*applicationClient).AddChangeListener(c2)

	config = &apollo.ApolloConfig{
		Log:            logger.GetLogger(),
		Client:         *applicationClient,
		CacheInterface: *applicationConfig,
		Namespace:      namespace,
	}

	parseAllConfig()

	return &apolloConfig{}
}

func initApollo(appId string, namespace string, cluster string, address string) (*agcache.CacheInterface, *agollo.Client) {

	client := apollo.InitApolloClient(appId, namespace, cluster, address)

	nsConfig := (*client).GetConfigCache(namespace)

	if nsConfig == nil {
		logger.GetLogger().Error("", zap.String("appId", appId), zap.String("namespace", namespace), zap.String("cluster", cluster), zap.String("address", address))
		panic("init apollo failure")
	}

	logger.GetLogger().Info("初始化Apollo配置成功", zap.String("appId", appId), zap.String("namespace", namespace), zap.String("cluster", cluster), zap.String("address", address))

	return &nsConfig, client

}

type apolloConfig struct {
}

func (c *apolloConfig) ForwardAddress() string {
	return getConfigString("proxy.forwardAddress")
}

func (c *apolloConfig) ServerPort() int {
	return config.GetIntValue("application.port", 0)
}

func (c *apolloConfig) ServerContextPath() string {
	return config.GetString("application.contextPath")
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

func (c *apolloConfig) GetSignKeyByName(signName string) string {

	v, ok := signKeyMap[signName]
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

func (c *apolloConfig) GetCallTypeConfigByName(name string, uri string) int {

	v, ok := serviceCallConfigMap[name]
	if ok {
		v, ok := v[uri]
		if ok {
			return v.CallType
		}
	}
	return 0
}

func (c *apolloConfig) ContainsCallPrivilege(name string, uri string) bool {
	v, ok := serviceCallConfigMap[name]
	if ok {
		_, ok := v[uri]
		return ok
	}
	return false
}

var (
	serviceConfigMap        map[string]ServiceConfig
	encryptKeyMap           EncryptKeyConfig
	signKeyMap              SignKeyConfig
	serviceTimeoutConfigMap map[string]map[string]ServiceTimeoutConfig
	serviceCallConfigMap    map[string]map[string]ServiceCallConfig
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
func setServiceCallConfig(m map[string]map[string]ServiceCallConfig) {
	serviceCallConfigMap = m
}

const (
	serviceConfigName        = "proxy.service.config"
	signKeyConfigName        = "proxy.sign.key.config"
	encryptKeyConfigName     = "proxy.encrypt.key.config"
	serviceTimeoutConfigName = "proxy.service.timeout.config"
	serviceCallConfigName    = "proxy.service.call.config"
)

func parseAllConfig() {
	serviceConfigJson := getConfigString(serviceConfigName)
	if !ParseConfig(serviceConfigJson, setServiceConfig, serviceConfigName) {

		logger.GetLogger().Error("Unmarshal failure ... ", zap.String("name", serviceConfigName))
		panic(fmt.Sprintf("Unmarshal %s failure", serviceConfigName))
	}

	encryptKeyConfigJson := getConfigString(encryptKeyConfigName)

	if !ParseConfig(encryptKeyConfigJson, setEncryptKey, encryptKeyConfigName) {

		logger.GetLogger().Error("Unmarshal failure ... ", zap.String("name", encryptKeyConfigName))
		panic(fmt.Sprintf("Unmarshal %s failure", encryptKeyConfigName))
	}

	signKeyConfigJson := getConfigString(signKeyConfigName)

	if !ParseConfig(signKeyConfigJson, setSignKey, signKeyConfigName) {

		logger.GetLogger().Error("Unmarshal failure ... ", zap.String("name", signKeyConfigName))
		panic(fmt.Sprintf("Unmarshal %s failure", signKeyConfigName))
	}

	serviceTimeoutConfig := getConfigString(serviceTimeoutConfigName)

	if !ParseConfig(serviceTimeoutConfig, setServiceTimeoutConfig, serviceTimeoutConfigName) {

		logger.GetLogger().Error("Unmarshal failure ... ", zap.String("name", serviceTimeoutConfigName))
		panic(fmt.Sprintf("Unmarshal %s failure", serviceTimeoutConfigName))
	}

	serviceCallConfig := getConfigString(serviceCallConfigName)

	if !ParseConfig(serviceCallConfig, setServiceCallConfig, serviceCallConfigName) {

		logger.GetLogger().Error("Unmarshal failure ... ", zap.String("name", serviceCallConfigName))
		panic(fmt.Sprintf("Unmarshal %s failure", serviceCallConfigName))
	}
}

func getConfigString(configName string) string {
	c := config.GetString(configName)
	if c == "" {
		logger.GetLogger().Error("get failure ... ", zap.String("name", configName))
		panic(fmt.Sprintf("get %s failure", configName))
	}
	return c
}

func ParseConfig[K string, V any](configJson string, p func(map[K]V), configName string) bool {
	b := []byte(configJson)
	m := make(map[K]V)

	err := json.Unmarshal(b, &m)
	if err != nil {
		logger.GetLogger().Error("get failure ... ", zap.String("name", configName))
		return false
	}
	p(m)

	return true
}

type CustomChangeListener struct {
	//wg sync.WaitGroup
}

func (c *CustomChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	//write your code here
	//fmt.Println(changeEvent.Changes)
	for key, value := range changeEvent.Changes {

		switch key {
		case serviceConfigName:
			ParseConfig(value.NewValue.(string), setServiceConfig, serviceConfigName)
		case signKeyConfigName:
			ParseConfig(value.NewValue.(string), setSignKey, signKeyConfigName)
		case encryptKeyConfigName:
			ParseConfig(value.NewValue.(string), setEncryptKey, encryptKeyConfigName)
		case serviceTimeoutConfigName:
			ParseConfig(value.NewValue.(string), setServiceTimeoutConfig, serviceTimeoutConfigName)
		case serviceCallConfigName:
			ParseConfig(value.NewValue.(string), setServiceCallConfig, serviceCallConfigName)
		default:
		}

		logger.GetLogger().Info("change key ", zap.Any(key, value))

	}
	//fmt.Println(changeEvent.Namespace)
	//c.wg.Done()
}

func (c *CustomChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
	//write your code here
}
