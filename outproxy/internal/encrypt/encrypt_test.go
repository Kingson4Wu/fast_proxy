package encrypt

import (
    "testing"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
)

type cfg struct{}
func (cfg) ServerName() string { return "out" }
func (cfg) ServerPort() int { return 0 }
func (cfg) ServiceRpcHeaderName() string { return "C_Service" }
func (cfg) GetServiceConfig(string) *config.ServiceConfig { return &config.ServiceConfig{EncryptKeyName: "enc"} }
func (cfg) GetSignKey(*config.ServiceConfig) string { return "" }
func (cfg) GetSignKeyByName(string) string { return "" }
func (cfg) GetEncryptKey(*config.ServiceConfig) string { return "0123456789ABCDEF" }
func (cfg) GetEncryptKeyByName(string) string { return "0123456789ABCDEF" }
func (cfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (cfg) HttpClientMaxIdleConns() int { return 10 }
func (cfg) HttpClientMaxIdleConnsPerHost() int { return 10 }
func (cfg) FastHttpEnable() bool { return false }
func (cfg) ForwardAddress() string { return "http://localhost" }

type badcfg struct{ cfg }
func (badcfg) GetEncryptKey(*config.ServiceConfig) string { return "" }
func (badcfg) GetEncryptKeyByName(string) string { return "" }

func TestEncodeReqDecodeResp(t *testing.T) {
    outconfig.Read(cfg{})
    data := []byte("hello")
    sc := &config.ServiceConfig{EncryptKeyName: "enc", EncryptEnable: true}
    enc, err := EncodeReq(data, sc)
    if err != nil { t.Fatalf("EncodeReq err: %v", err) }
    dec, err := DecodeResp(enc, "enc")
    if err != nil { t.Fatalf("DecodeResp err: %v", err) }
    if string(dec) != string(data) { t.Fatalf("mismatch") }

    outconfig.Read(badcfg{})
    if _, err := EncodeReq(data, sc); err == nil { t.Fatalf("expected failure with empty key") }
    if _, err := DecodeResp(enc, "enc"); err == nil { t.Fatalf("expected failure with empty key name") }
}

