package main

import (
	"embed"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"github.com/Kingson4Wu/fast_proxy/outproxy"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
	"os"
)

//go:embed *
var ConfigFs embed.FS

func main() {

	configBytes, err := ConfigFs.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	c := outconfig.LoadYamlConfig(configBytes)

	sc := center.GetSC(func() string { return outconfig.Get().ServiceRpcHeaderName() })

	outproxy.NewServer(c, server.WithServiceCenter(sc))

}
