package server

import (
    "net/http"
    "os"
    "testing"
    "time"

    c "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/logger"
)

type startCfg struct{}

func (startCfg) ServerName() string { return "srv" }
func (startCfg) ServerPort() int { return 0 }
func (startCfg) ServiceRpcHeaderName() string { return "C_Service" }
func (startCfg) GetServiceConfig(string) *c.ServiceConfig { return &c.ServiceConfig{} }
func (startCfg) GetSignKey(*c.ServiceConfig) string { return "" }
func (startCfg) GetSignKeyByName(string) string { return "" }
func (startCfg) GetEncryptKey(*c.ServiceConfig) string { return "" }
func (startCfg) GetEncryptKeyByName(string) string { return "" }
func (startCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (startCfg) HttpClientMaxIdleConns() int { return 10 }
func (startCfg) HttpClientMaxIdleConnsPerHost() int { return 10 }
func (startCfg) FastHttpEnable() bool { return false }

// This test exercises Start() and the signal-driven shutdown path.
func TestStartAndSignalShutdown(t *testing.T) {
    p := NewServer(startCfg{}, logger.DiscardLogger, func(http.ResponseWriter, *http.Request) {})
    // Start in background with a short shutdown timeout
    go p.Start(WithShutdownTimeout(200 * time.Millisecond))
    // Give it a moment to set global state
    time.Sleep(20 * time.Millisecond)
    // Send an interrupt to trigger shutdown path
    _ = os.Interrupt
    _ = os.Kill // keep os imported
    _ = os.Getpid()
    _ = os.FindProcess
    proc, _ := os.FindProcess(os.Getpid())
    _ = proc.Signal(os.Interrupt)
    // Wait a bit for shutdown to complete
    time.Sleep(100 * time.Millisecond)
}

