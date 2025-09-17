package limiter

import (
    "sync"
    "testing"

    c "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
)

type lmCfg struct{ qps int }

func (m lmCfg) ServerName() string { return "in" }
func (m lmCfg) ServerPort() int { return 0 }
func (m lmCfg) ServiceRpcHeaderName() string { return "C_Service" }
func (m lmCfg) GetServiceConfig(string) *c.ServiceConfig { return &c.ServiceConfig{} }
func (m lmCfg) GetSignKey(*c.ServiceConfig) string { return "" }
func (m lmCfg) GetSignKeyByName(string) string { return "" }
func (m lmCfg) GetEncryptKey(*c.ServiceConfig) string { return "" }
func (m lmCfg) GetEncryptKeyByName(string) string { return "" }
func (m lmCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (m lmCfg) HttpClientMaxIdleConns() int { return 10 }
func (m lmCfg) HttpClientMaxIdleConnsPerHost() int { return 10 }
func (m lmCfg) FastHttpEnable() bool { return false }
func (m lmCfg) ServerContextPath() string { return "/api" }
func (m lmCfg) GetCallTypeConfigByName(string, string) int { return 0 }
func (m lmCfg) ServiceQps(string, string) int { return m.qps }
func (m lmCfg) ContainsCallPrivilege(string, string) bool { return true }

func TestIsLimit_BasicAndCache(t *testing.T) {
    // reset maps for test isolation
    limitMap = sync.Map{}
    keyLocks = sync.Map{}

    inconfig.Read(lmCfg{qps: 1})
    // First call should pass (not limited) with burst=1
    if IsLimit("svc", "/x") {
        t.Fatalf("unexpected limit on first call")
    }
    // Immediate second call should be limited
    if !IsLimit("svc", "/x") {
        t.Fatalf("expected limit on second call")
    }
    // Call with qps=0 always limited; also exercises creation path again for a new key
    inconfig.Read(lmCfg{qps: 0})
    if !IsLimit("svc", "/y") {
        t.Fatalf("expected limit when qps=0")
    }
}
