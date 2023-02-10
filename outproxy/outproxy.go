package outproxy

import (
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/outproxy/config"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/logger"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func init() {

	/* ... */
	//todo

	//go build -gcflags '-m -m -l' main.go 测试环境开启逃逸分析
	//go build命令传入-x -v命令行选项来输出详细的构建日志信息
	//查看调度器状态 GOMAXPROCS=1 GODEBUG=schedtrace=1000
	//设置环境变量GODEBUG='gctrace=1'让位于Go程序中的运行时在每次GC执行时输出此次GC相关的跟踪信息。

	/* if err := WaitForServer(url); err != nil {
			fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
			os.Exit(1)

		}
			if err := WaitForServer(url); err != nil {
	    log.Fatalf("Site is down: %v\n", err)
	}
	*/

}

func requestProxy(res http.ResponseWriter, req *http.Request) {

	proxyType := proxy.FORWARD

	p := proxy.GetProxy(proxyType)
	if p != nil {
		p.Handle(res, req)
	}

}

func NewServer(c config.Config) {
	config.Configuration = c
	proxy := server.NewServer(c.ServerPort(), logger.GetLogger(), requestProxy)
	proxy.RegisterOnShutdown(func() {
		logger.GetLogger().Info("clean resources on shutdown...")
		time.Sleep(2 * time.Second)
		logger.GetLogger().Info("clean resources ok")
	})

	proxy.Start(server.WithShutdownTimeout(5 * time.Second))
}
