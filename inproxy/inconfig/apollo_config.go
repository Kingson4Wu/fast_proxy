package inconfig

import (
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/apolloconfig/agollo/v4/storage"
)

func LoadApolloConfig(appId string, namespace string, cluster string, address string, logger logger.Logger) Config {

	c := &apolloConfig{
		config.LoadApolloConfig(appId, namespace, cluster, address, logger),
	}
	parseAllConfig()

	config.RegisterApolloListener(listen)

	return c
}

type apolloConfig struct {
	config.Config
}

func (c *apolloConfig) ServerContextPath() string {
	return config.ApolloConfig().GetString("application.contextPath")
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
	serviceCallConfigMap map[string]map[string]ServiceCallConfig
)

func setServiceCallConfig(m map[string]map[string]ServiceCallConfig) {
	serviceCallConfigMap = m
}

const (
	serviceCallConfigName = "proxy.service.call.config"
)

func parseAllConfig() {
	serviceCallConfig := config.ApolloConfig().GetString(serviceCallConfigName)

	if !config.ParseApolloConfig(serviceCallConfig, setServiceCallConfig, serviceCallConfigName) {

		config.ApolloConfig().Log.Error("Unmarshal failure ... ", "name", serviceCallConfigName)
		panic(fmt.Sprintf("Unmarshal %s failure", serviceCallConfigName))
	}
}

func listen(changeEvent *storage.ChangeEvent) {
	for key, value := range changeEvent.Changes {
		switch key {
		case serviceCallConfigName:
			config.ParseApolloConfig(value.NewValue.(string), setServiceCallConfig, serviceCallConfigName)
		default:
		}
	}
}
