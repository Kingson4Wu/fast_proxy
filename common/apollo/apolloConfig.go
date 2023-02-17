package apollo

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/env/config"
	"go.uber.org/zap"
)

var c *ApolloConfig

type ApolloConfig struct {
	Log *zap.Logger
	agollo.Client
	agcache.CacheInterface
	Namespace string
}

func (ac *ApolloConfig) GetString(name string) string {
	return ac.GetConfig(ac.Namespace).GetStringValue(name, "")
}

func (ac *ApolloConfig) GetStringDefault(name string, defaultValue string) string {
	return ac.GetConfig(ac.Namespace).GetStringValue(name, defaultValue)
}

func InitApolloClient(appId string, namespace string, cluster string, address string) *agollo.Client {

	c := &config.AppConfig{
		AppID:         appId,
		Cluster:       cluster,
		IP:            address,
		NamespaceName: namespace,
	}

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})

	if err != nil {
		panic(fmt.Sprintf("init apollo client failure, appId:%s, namespace:%s, cluster:%s, address:%s", appId, namespace, cluster, address))
	}

	return &client
}

func init() {

}
