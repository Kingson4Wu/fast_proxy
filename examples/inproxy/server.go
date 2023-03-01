package main

import (
	"embed"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"github.com/Kingson4Wu/fast_proxy/inproxy"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
	"os"
)

//go:embed *
var ConfigFs embed.FS

func main() {

	configBytes, err := ConfigFs.ReadFile("config.yaml")
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	c := inconfig.LoadYamlConfig(configBytes)

	rpcHeaderNameFunc := func() string { return c.ServiceRpcHeaderName() }

	sc := center.GetSC(rpcHeaderNameFunc)

	inproxy.NewServer(c, server.WithServiceCenter(sc))

}
