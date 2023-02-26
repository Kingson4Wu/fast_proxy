package main

import (
	"github.com/Kingson4Wu/fast_proxy/common/logger/zap"
	"github.com/Kingson4Wu/fast_proxy/outproxy"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
)

func main() {

	outproxy.NewServer(outconfig.LoadApolloConfig("song_service", "application", "default", "http://192.168.33.174:8080", zap.DefaultLogger()))
}
