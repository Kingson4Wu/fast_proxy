package main

import (
	"embed"
	"github.com/Kingson4Wu/fast_proxy/inproxy"
	"github.com/Kingson4Wu/fast_proxy/inproxy/config"
)

//go:embed *
var ConfigFs embed.FS

func main() {

	configBytes, err := ConfigFs.ReadFile("config.yaml")
	if err != nil {
		panic(err.Error())
	}

	inproxy.NewServer(config.LoadYamlConfig(configBytes))

}
