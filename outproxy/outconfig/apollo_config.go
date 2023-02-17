package outconfig

import (
	"github.com/Kingson4Wu/fast_proxy/common/config"
)

func LoadApolloConfig(appId string, namespace string, cluster string, address string) Config {
	return &apolloConfig{
		config.LoadApolloConfig(appId, namespace, cluster, address),
	}
}

type apolloConfig struct {
	config.Config
}

func (c *apolloConfig) ForwardAddress() string {
	return config.ApolloConfig().GetString("proxy.forwardAddress")
}
