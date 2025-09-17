package pack

import (
    "testing"
    "time"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/logger/zap"
    "github.com/Kingson4Wu/fast_proxy/common/server"
    "github.com/Kingson4Wu/fast_proxy/common/sign"
    "github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
    "github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
    "google.golang.org/protobuf/proto"
    "net/http"
)

type signCfg struct{}

func (signCfg) ServerName() string { return "out" }
func (signCfg) ServerPort() int { return 0 }
func (signCfg) ServiceRpcHeaderName() string { return "C_Service" }
func (signCfg) GetServiceConfig(name string) *config.ServiceConfig {
    if name == "svcSign" {
        return &config.ServiceConfig{SignEnable: true, SignKeyName: "sig"}
    }
    return &config.ServiceConfig{}
}
func (signCfg) GetSignKey(*config.ServiceConfig) string { return "S-KEY" }
func (signCfg) GetSignKeyByName(string) string { return "S-KEY" }
func (signCfg) GetEncryptKey(*config.ServiceConfig) string { return "" }
func (signCfg) GetEncryptKeyByName(string) string { return "" }
func (signCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (signCfg) HttpClientMaxIdleConns() int { return 10 }
func (signCfg) HttpClientMaxIdleConnsPerHost() int { return 10 }
func (signCfg) FastHttpEnable() bool { return false }
func (signCfg) ForwardAddress() string { return "http://127.0.0.1" }

func TestEncode_SignEnabled_SetsSignature(t *testing.T) {
    // Set config for outconfig and server
    cfg := signCfg{}
    outconfig.Read(cfg)
    if server.Config() == nil {
        go func() {
            p := server.NewServer(cfg, zap.DefaultLogger(), func(http.ResponseWriter, *http.Request) {})
            p.Start()
        }()
        // wait for server.Config to be set
        deadline := time.Now().Add(300 * time.Millisecond)
        for time.Now().Before(deadline) {
            if server.Config() != nil { break }
            time.Sleep(5 * time.Millisecond)
        }
    }

    body := []byte("hello-sign")
    pdata, err := Encode(body, "svcSign")
    if err != nil { t.Fatalf("Encode err: %v", err) }
    var pd protobuf.ProxyData
    if err := proto.Unmarshal(pdata, &pd); err != nil { t.Fatalf("unmarshal: %v", err) }
    if pd.Sign == "" || !pd.SignEnable {
        t.Fatalf("expected sign present and enabled; got sign=%q enable=%v", pd.Sign, pd.SignEnable)
    }
    expected, _ := sign.GenerateSign(body, "S-KEY")
    if pd.Sign != expected {
        t.Fatalf("sign mismatch: got %s want %s", pd.Sign, expected)
    }
}
