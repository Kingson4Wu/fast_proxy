package inproxy

import (
	"github.com/Kingson4Wu/fast_proxy/common/logger/zap"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/proxy"
	"net/http"
	"time"
)

func requestProxy(res http.ResponseWriter, req *http.Request) {
	proxy.DoProxy(res, req)

}

func NewServer(c inconfig.Config, opts ...server.Option) {
	inconfig.Read(c)
	proxy.BuildClient(c)
	p := server.NewServer(c, zap.DefaultLogger(), requestProxy)
	p.RegisterOnShutdown(func() {
		server.GetLogger().Info("clean resources on shutdown...")
		time.Sleep(2 * time.Second)
		server.GetLogger().Info("clean resources ok")
	})

	var options []server.Option
	options = append(options, server.WithShutdownTimeout(5*time.Second))
	options = append(options, opts...)
	p.Start(options...)
}
