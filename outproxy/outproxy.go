package outproxy

import (
	"embed"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/logger/zap"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
	"net/http"
	_ "net/http/pprof"
	"time"
)

//go:embed *
var ConfigFs embed.FS

func requestProxy(res http.ResponseWriter, req *http.Request) {

	proxyType := proxy.FORWARD

	p := proxy.GetProxy(proxyType)
	if p != nil {
		p.Handle(res, req)
	}

}

func NewServer(c outconfig.Config, opts ...server.Option) {
    outconfig.Read(c)
    proxy.BuildClient(c)
    p := server.NewServer(c, zap.DefaultLogger(), requestProxy)
    p.RegisterOnShutdown(func() {
        server.GetLogger().Info("clean resources on shutdown...")
        proxy.Close()
    })

	var options []server.Option
	options = append(options, server.WithShutdownTimeout(5*time.Second))
	options = append(options, opts...)
	p.Start(options...)
}

func init() {
	configBytes, err := ConfigFs.ReadFile("banner.txt")
	if err == nil {
		fmt.Println(string(configBytes))
	}
}
