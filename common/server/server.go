package server

import (
	"context"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/Kingson4Wu/fast_proxy/common/network"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"github.com/felixge/fgprof"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var server *Proxy

func Center() *servicediscovery.ServiceCenter {
	return server.sc
}

func Config() config.Config {
	return server.c
}

func GetLogger() logger.Logger {
	return server.logger
}

type Proxy struct {
	svr             *http.Server
	port            int
	proxyHandler    func(http.ResponseWriter, *http.Request)
	otherHandlers   map[string]func(http.ResponseWriter, *http.Request)
	wg              sync.WaitGroup
	logger          logger.Logger
	shutdownTimeout time.Duration
	sc              *servicediscovery.ServiceCenter
	c               config.Config
}

func New() *Proxy {
	return &Proxy{otherHandlers: make(map[string]func(http.ResponseWriter, *http.Request))}
}

type Option func(*Proxy)

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(p *Proxy) {
		p.shutdownTimeout = timeout
	}
}

func WithServiceCenter(sc *servicediscovery.ServiceCenter) Option {
	return func(p *Proxy) {
		p.sc = sc
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(p *Proxy) {
		p.logger = logger
	}
}

func WithCustomHandler(handler func(http.ResponseWriter, *http.Request)) Option {
	return func(p *Proxy) {
		p.proxyHandler = handler
	}
}

func NewServer(config config.Config, logger logger.Logger, proxyHandler func(http.ResponseWriter, *http.Request)) *Proxy {

	//http.HandleFunc("/", proxyHandler)

	svr := http.Server{
		Addr: ":" + strconv.Itoa(config.ServerPort()),
	}

	p := New()
	p.svr = &svr
	p.port = config.ServerPort()
	p.proxyHandler = proxyHandler
	p.logger = logger
	p.c = config

	return p
}

func (p *Proxy) AddHandler(uri string, handler func(http.ResponseWriter, *http.Request)) *Proxy {
	//http.HandleFunc(uri, handler)

	if _, ok := p.otherHandlers[uri]; !ok {
		p.otherHandlers[uri] = handler
	}

	return p
}

func (p *Proxy) RegisterOnShutdown(f func()) {
	p.svr.RegisterOnShutdown(func() {
		f()
		p.wg.Done()
	})
	p.wg.Add(1)
}

func (p *Proxy) Start(opts ...Option) {

	server = p

	for _, opt := range opts {
		opt(p)
	}
	p.wg.Add(1)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		<-quit

		var ctx context.Context
		if p.shutdownTimeout > 0 {
			var cf context.CancelFunc
			ctx, cf = context.WithTimeout(context.Background(), p.shutdownTimeout)
			defer cf()
		} else {
			ctx = context.TODO()
		}

		var done = make(chan struct{}, 1)
		go func() {
			if err := p.svr.Shutdown(ctx); err != nil {
				p.logger.Error("proxy server shutdown error ", zap.Any("err", err))
			} else {
				p.logger.Info("proxy server shutdown ok")
			}
			done <- struct{}{}
			p.wg.Done()
		}()
		select {
		case <-ctx.Done():
			if p.shutdownTimeout > 0 {
				p.logger.Info("proxy server shutdown timeout")
			} else {
				p.logger.Info("proxy server shutdown")
			}
		case <-done:
		}
	}()

	pid := os.Getpid()
	p.logger.Info("server start ...", zap.Int("port", p.port), zap.Int("pid", pid))

	intranetIp := network.GetIntranetIp()
	var stop chan bool
	if p.sc != nil {
		stop = p.sc.Register(p.c.ServerName(), intranetIp, p.c.ServerPort())
	}

	var err error

	if p.c.FastHttpEnable() {

		handler := fasthttpadaptor.NewFastHTTPHandlerFunc(p.proxyHandler)

		pprofHandler := fasthttpadaptor.NewFastHTTPHandler(fgprof.Handler())

		otherHandlers := make(map[string]fasthttp.RequestHandler)

		for uri, f := range p.otherHandlers {
			otherHandlers[uri] = fasthttpadaptor.NewFastHTTPHandlerFunc(f)
		}

		otherHandlers["/debug/pprof"] = pprofHandler

		m := func(ctx *fasthttp.RequestCtx) {

			path := string(ctx.Path())
			if strings.HasPrefix(path, "/debug/pprof") {
				pprofHandler(ctx)
				return
			}

			if f, ok := otherHandlers[path]; ok {
				f(ctx)
				return
			}
			handler(ctx)
		}

		err = fasthttp.ListenAndServe(":"+strconv.Itoa(p.port), m)
	} else {
		http.HandleFunc("/", p.proxyHandler)
		for uri, f := range p.otherHandlers {
			http.HandleFunc(uri, f)
		}
		err = p.svr.ListenAndServe()
	}

	if err != nil {
		if err != http.ErrServerClosed {
			close(stop)
			p.logger.Error("proxy server start failed ", zap.Any("err", err))
			return
		}
	}
	p.wg.Wait()
	p.logger.Info("proxy shutdown ok")
}

func init() {
	// print banner
	path, _ := filepath.Abs("resource/banner.txt")
	banner, _ := os.ReadFile(path)
	fmt.Println(string(banner))
}
