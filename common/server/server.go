package server

import (
	"context"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/network"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

type Proxy struct {
	svr             *http.Server
	port            int
	wg              sync.WaitGroup
	logger          *zap.Logger
	shutdownTimeout time.Duration
	sc              *servicediscovery.ServiceCenter
	c               config.Config
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

func NewServer(config config.Config, logger *zap.Logger, proxyHandler func(http.ResponseWriter, *http.Request)) *Proxy {

	http.HandleFunc("/", proxyHandler)
	svr := http.Server{
		Addr: ":" + strconv.Itoa(config.ServerPort()),
	}
	return &Proxy{
		svr:    &svr,
		port:   config.ServerPort(),
		logger: logger,
		c:      config,
	}

}

func (p *Proxy) AddHandler(uri string, handler func(http.ResponseWriter, *http.Request)) *Proxy {
	http.HandleFunc(uri, handler)
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
	stop := p.sc.Register(p.c.ServerName(), intranetIp, p.c.ServerPort())

	err := p.svr.ListenAndServe()
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
