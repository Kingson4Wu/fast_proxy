package outconfig

import (
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
)

func LoadApolloConfig(appId string, namespace string, cluster string, address string, logger logger.Logger) Config {
	return &apolloConfig{
		config.LoadApolloConfig(appId, namespace, cluster, address, logger),
	}
}

type apolloConfig struct {
	config.Config
}

func (c *apolloConfig) ForwardAddress() string {
	return config.ApolloConfig().GetString("proxy.forwardAddress")
}
