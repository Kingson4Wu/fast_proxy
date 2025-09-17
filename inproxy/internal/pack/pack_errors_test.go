package pack

import (
    "net/http"
    "testing"
    "time"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/logger/zap"
    "github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
    "github.com/Kingson4Wu/fast_proxy/common/server"
    "google.golang.org/protobuf/proto"
)

// setup a minimal server for sign key lookup
type signCfg struct{}
func (signCfg) ServerName() string { return "in" }
func (signCfg) ServerPort() int { return 0 }
func (signCfg) ServiceRpcHeaderName() string { return "C_Service" }
func (signCfg) GetServiceConfig(string) *config.ServiceConfig { return &config.ServiceConfig{} }
func (signCfg) GetSignKey(*config.ServiceConfig) string { return "" }
func (signCfg) GetSignKeyByName(string) string { return "sek" }
func (signCfg) GetEncryptKey(*config.ServiceConfig) string { return "" }
func (signCfg) GetEncryptKeyByName(string) string { return "" }
func (signCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (signCfg) HttpClientMaxIdleConns() int { return 10 }
func (signCfg) HttpClientMaxIdleConnsPerHost() int { return 10 }
func (signCfg) FastHttpEnable() bool { return true }

func TestDecode_SignMismatchAndCompressError(t *testing.T) {
    // Start a server so server.Config is available for sign lookups
    go func() {
        p := server.NewServer(signCfg{}, zap.DefaultLogger(), func(http.ResponseWriter, *http.Request) {})
        p.Start()
    }()
    deadline := time.Now().Add(300 * time.Millisecond)
    for time.Now().Before(deadline) {
        if server.Config() != nil { break }
        time.Sleep(5 * time.Millisecond)
    }
    // 1) Sign mismatch
    pd := &protobuf.ProxyData{Payload: []byte("hi"), SignEnable: true, SignKeyName: "sig", Sign: "wrong"}
    b, _ := proto.Marshal(pd)
    if _, _, err := Decode(b); err == nil {
        t.Fatalf("expected sign mismatch error")
    }

    // 2) Compress error: bogus payload with Compress=true
    pd2 := &protobuf.ProxyData{Payload: []byte("bogus"), Compress: true, CompressAlgorithm: 0}
    b2, _ := proto.Marshal(pd2)
    if _, _, err := Decode(b2); err == nil {
        t.Fatalf("expected compress decode failure")
    }
}
