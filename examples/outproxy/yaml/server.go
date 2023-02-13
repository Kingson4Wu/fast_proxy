package main

import (
	"embed"
	"github.com/Kingson4Wu/fast_proxy/common/network"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"github.com/Kingson4Wu/fast_proxy/outproxy"
	"github.com/Kingson4Wu/fast_proxy/outproxy/config"
)

//go:embed *
var ConfigFs embed.FS

func main() {

	configBytes, err := ConfigFs.ReadFile("config.yaml")
	if err != nil {
		panic(err.Error())
	}

	intranetIp := network.GetIntranetIp()
	c := config.LoadYamlConfig(configBytes)

	//close(stop) todo
	//stop := center.RegisterAsync("token_service", intranetIp, c.ServerPort())
	center.RegisterAsync("out_proxy", intranetIp, c.ServerPort())

	outproxy.NewServer(c)

}
