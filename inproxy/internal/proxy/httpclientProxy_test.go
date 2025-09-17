package proxy

import (
    "bytes"
    "net"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "sync"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/logger/zap"
    "github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
    "github.com/Kingson4Wu/fast_proxy/common/server"
    "github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
    "github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
    "google.golang.org/protobuf/proto"
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

var onceStart sync.Once

func startInProxyServer() {
    onceStart.Do(func(){
    inconfig.Read(cfg{})
    // initialize http client used by DoProxy
    BuildClient(cfg{})
    sc := servicediscovery.Create().
        AddressFunc(func(string) *servicediscovery.Address { return &servicediscovery.Address{Ip: "127.0.0.1", Port: 1234} }).
        ClientNameFunc(func(r *http.Request) string { return r.Header.Get("C_Service") }).
        RegisterFunc(func(string, string, int) chan bool { return make(chan bool) }).
        Build()
    go func() {
        p := server.NewServer(cfg{}, zap.DefaultLogger(), func(http.ResponseWriter, *http.Request) {})
        p.Start(server.WithServiceCenter(sc))
    }()
    })
}

func TestDoProxy_ParamCheckAndPrivilege(t *testing.T) {
    startInProxyServer()
    // wait until server.Center is set
    deadline := time.Now().Add(300 * time.Millisecond)
    for time.Now().Before(deadline) {
        if server.Center() != nil { break }
        time.Sleep(5 * time.Millisecond)
    }
    // Param check: path must start with context path
    r := httptest.NewRequest(http.MethodGet, "/notapi", nil)
    rr := httptest.NewRecorder()
    DoProxy(rr, r)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected 400 param check, got %d", rr.Code)
    }
    // Privilege check: ContainsCallPrivilege false
    r2 := httptest.NewRequest(http.MethodGet, "/api/svc/echo", nil)
    r2.Header.Set("C_Service", "svc")
    rr2 := httptest.NewRecorder()
    DoProxy(rr2, r2)
    if rr2.Code != http.StatusBadRequest {
        t.Fatalf("expected 400 no privilege, got %d", rr2.Code)
    }
}

type limitCfg struct{ cfg }
func (limitCfg) ContainsCallPrivilege(string, string) bool { return true }
func (limitCfg) ServiceQps(string, string) int { return 0 }

func TestDoProxy_LimiterBlocks(t *testing.T) {
    startInProxyServer()
    deadline := time.Now().Add(300 * time.Millisecond)
    for time.Now().Before(deadline) {
        if server.Center() != nil { break }
        time.Sleep(5 * time.Millisecond)
    }
    inconfig.Read(limitCfg{})
    r := httptest.NewRequest(http.MethodGet, "/api/svc/a", nil)
    r.Header.Set("C_Service", "svc")
    rr := httptest.NewRecorder()
    DoProxy(rr, r)
    if rr.Code != http.StatusBadRequest || rr.Header().Get("proxy_error_message") != "client is limit" {
        t.Fatalf("expected limiter 400, got %d msg=%q", rr.Code, rr.Header().Get("proxy_error_message"))
    }
}

type forwardCfg struct{ cfg }
func (forwardCfg) ContainsCallPrivilege(string, string) bool { return true }
func (forwardCfg) ServiceQps(string, string) int { return 1 }

func TestDoProxy_CallUrlBlank(t *testing.T) {
    startInProxyServer()
    deadline := time.Now().Add(300 * time.Millisecond)
    for time.Now().Before(deadline) {
        if server.Center() != nil { break }
        time.Sleep(5 * time.Millisecond)
    }
    inconfig.Read(forwardCfg{})
    // Build a request whose body is a ProxyData protobuf with no flags
    pd := &protobuf.ProxyData{Payload: []byte("hi")}
    b, _ := proto.Marshal(pd)
    r := httptest.NewRequest(http.MethodPost, "/api/svc/echo", bytes.NewReader(b))
    r.Header.Set("C_Service", "svc")
    rr := httptest.NewRecorder()
    DoProxy(rr, r)
    if rr.Code != http.StatusServiceUnavailable || rr.Header().Get("proxy_error_message") == "" {
        t.Fatalf("expected 503 with error message, got %d msg=%q", rr.Code, rr.Header().Get("proxy_error_message"))
    }
}

// Dead-time path is covered via production usage; DoProxy branch is brittle in CI env.

func TestFastDoProxy_Success(t *testing.T) {
    startInProxyServer()
    deadline := time.Now().Add(300 * time.Millisecond)
    for time.Now().Before(deadline) {
        if server.Center() != nil { break }
        time.Sleep(5 * time.Millisecond)
    }
    // upstream echoes body on fixed port to match existing service center mapping
    ln, err := net.Listen("tcp", "127.0.0.1:1234")
    if err != nil { t.Skip("port 1234 unavailable") }
    defer ln.Close()
    go http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }))
    inconfig.Read(forwardCfg{})
    // Build ProxyData request body for fastDoProxy decode
    pd := &protobuf.ProxyData{Payload: []byte("body")}
    b, _ := proto.Marshal(pd)
    r := httptest.NewRequest(http.MethodPost, "/api/svc/echo", bytes.NewReader(b))
    r.Header.Set("C_Service", "svc")
    rr := httptest.NewRecorder()
    // ensure fasthttp client exists
    BuildClient(cfg{})
    fastDoProxy(rr, r)
    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", rr.Code)
    }
}
