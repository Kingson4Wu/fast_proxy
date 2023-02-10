package server

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Proxy struct {
	svr             *http.Server
	port            int
	wg              sync.WaitGroup
	logger          *zap.Logger
	shutdownTimeout time.Duration
}

type Option func(*Proxy)

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(p *Proxy) {
		p.shutdownTimeout = timeout
	}
}

func NewServer(port int, logger *zap.Logger, proxyHandler func(http.ResponseWriter, *http.Request)) *Proxy {
	http.HandleFunc("/", proxyHandler)
	svr := http.Server{
		Addr: ":" + strconv.Itoa(port),
	}
	return &Proxy{
		svr:    &svr,
		port:   port,
		logger: logger,
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

	err := p.svr.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			p.logger.Error("proxy server start failed ", zap.Any("err", err))
			return
		}
	}
	p.wg.Wait()
	p.logger.Info("proxy shutdown ok")
}
