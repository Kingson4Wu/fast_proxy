package main

import (
	"embed"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"github.com/Kingson4Wu/fast_proxy/inproxy"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
)

//go:embed *
var ConfigFs embed.FS

func main() {

	configBytes, err := ConfigFs.ReadFile("config.yaml")
	if err != nil {
		panic(err.Error())
	}

	c := inconfig.LoadYamlConfig(configBytes)

	sc := center.GetSC(func() string { return c.ServiceRpcHeaderName() })

	inproxy.NewServer(c, server.WithServiceCenter(sc))

}
