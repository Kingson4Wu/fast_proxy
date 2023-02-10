package main

import (
	"github.com/Kingson4Wu/fast_proxy/outproxy"
	"github.com/Kingson4Wu/fast_proxy/outproxy/config"
)

func main() {

	outproxy.NewServer(config.LoadApolloConfig("song_service", "application", "default", "http://192.168.33.174:8080"))
}
