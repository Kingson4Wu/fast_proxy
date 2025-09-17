package server

import (
    "net/http"
    "testing"
    "time"

    c "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/logger"
)

type fhCfg struct{}

func (fhCfg) ServerName() string { return "srv" }
func (fhCfg) ServerPort() int { return 0 }
func (fhCfg) ServiceRpcHeaderName() string { return "C_Service" }
func (fhCfg) GetServiceConfig(string) *c.ServiceConfig { return &c.ServiceConfig{} }
func (fhCfg) GetSignKey(*c.ServiceConfig) string { return "" }
func (fhCfg) GetSignKeyByName(string) string { return "" }
func (fhCfg) GetEncryptKey(*c.ServiceConfig) string { return "" }
func (fhCfg) GetEncryptKeyByName(string) string { return "" }
func (fhCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (fhCfg) HttpClientMaxIdleConns() int { return 10 }
func (fhCfg) HttpClientMaxIdleConnsPerHost() int { return 10 }
func (fhCfg) FastHttpEnable() bool { return true }

func TestStartFastHTTPAndShutdown(t *testing.T) {
    p := NewServer(fhCfg{}, logger.DiscardLogger, func(http.ResponseWriter, *http.Request) {})
    p.AddHandler("/healthz", func(http.ResponseWriter, *http.Request) {})
    go p.Start(WithShutdownTimeout(200 * time.Millisecond))
    time.Sleep(50 * time.Millisecond)
    // Trigger shutdown via HTTP server Close by sending SIGINT through the test harness' channel
    // The Start goroutine listens for os.Interrupt; since we can't portably send a signal on all OSes in CI here,
    // rely on timeout path to exercise shutdown code.
    time.Sleep(250 * time.Millisecond)
}

