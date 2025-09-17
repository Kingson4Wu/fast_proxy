package servicediscovery

import (
    "net/http"
    "testing"
    "time"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/logger/zap"
    "github.com/Kingson4Wu/fast_proxy/common/server"
    "github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
    "github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
)

type cfg struct{}
func (cfg) ServerName() string { return "in" }
func (cfg) ServerPort() int { return 0 }
func (cfg) ServiceRpcHeaderName() string { return "C_Service" }
func (cfg) GetServiceConfig(string) *config.ServiceConfig { return &config.ServiceConfig{} }
func (cfg) GetSignKey(*config.ServiceConfig) string { return "" }
func (cfg) GetSignKeyByName(string) string { return "" }
func (cfg) GetEncryptKey(*config.ServiceConfig) string { return "" }
func (cfg) GetEncryptKeyByName(string) string { return "" }
func (cfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (cfg) HttpClientMaxIdleConns() int { return 10 }
func (cfg) HttpClientMaxIdleConnsPerHost() int { return 10 }
func (cfg) FastHttpEnable() bool { return false }
func (cfg) ServerContextPath() string { return "/api" }
func (cfg) GetCallTypeConfigByName(string, string) int { return 0 }
func (cfg) ServiceQps(string, string) int { return 0 }
func (cfg) ContainsCallPrivilege(string, string) bool { return false }

func TestRealRequestUriAndForward(t *testing.T) {
    inconfig.Read(cfg{})
    if got := RealRequestUri("/api/svc/echo"); got != "/svc/echo" {
        t.Fatalf("RealRequestUri=%q", got)
    }
    // Setup a service center that maps svc->127.0.0.1:1234
    sc := servicediscovery.Create().
        AddressFunc(func(serviceName string) *servicediscovery.Address { return &servicediscovery.Address{Ip: "127.0.0.1", Port: 1234} }).
        ClientNameFunc(func(r *http.Request) string { return "svc" }).
        RegisterFunc(func(string, string, int) chan bool { return make(chan bool) }).
        Build()
    go func() {
        p := server.NewServer(cfg{}, zap.DefaultLogger(), func(http.ResponseWriter, *http.Request) {})
        p.Start(server.WithServiceCenter(sc))
    }()
    // wait brief for global server to be set
    deadline := time.Now().Add(300 * time.Millisecond)
    for time.Now().Before(deadline) {
        if server.Center() != nil { break }
        time.Sleep(5 * time.Millisecond)
    }
    req := &http.Request{Method: http.MethodGet, RequestURI: "/api/svc/echo"}
    url, handler := Forward(req)
    if url == "" {
        t.Fatalf("Forward unexpected empty url")
    }
    _ = handler // may be nil per current implementation
}
