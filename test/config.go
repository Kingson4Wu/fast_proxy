package test

import (
	"embed"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
)

//go:embed *
var ConfigFs embed.FS

func GetOutConfig() outconfig.Config {
	configBytes, err := ConfigFs.ReadFile("config/outconfig.yaml")
	if err != nil {
		panic(err.Error())
	}

	c := outconfig.LoadYamlConfig(configBytes)
	return c
}
