package encrypt

import (
    "testing"

    c "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
)

type encCfg struct{}

func (encCfg) ServerName() string { return "in" }
func (encCfg) ServerPort() int { return 0 }
func (encCfg) ServiceRpcHeaderName() string { return "C_Service" }
func (encCfg) GetServiceConfig(string) *c.ServiceConfig { return &c.ServiceConfig{} }
func (encCfg) GetSignKey(*c.ServiceConfig) string { return "" }
func (encCfg) GetSignKeyByName(string) string { return "" }
func (encCfg) GetEncryptKey(*c.ServiceConfig) string { return "" }
func (encCfg) GetEncryptKeyByName(name string) string { return "ABCDABCDABCDABCD" }
func (encCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (encCfg) HttpClientMaxIdleConns() int { return 10 }
func (encCfg) HttpClientMaxIdleConnsPerHost() int { return 10 }
func (encCfg) FastHttpEnable() bool { return true }
func (encCfg) ServerContextPath() string { return "/" }
func (encCfg) GetCallTypeConfigByName(string, string) int { return 0 }
func (encCfg) ServiceQps(string, string) int { return 0 }
func (encCfg) ContainsCallPrivilege(string, string) bool { return false }

type emptyEncCfg struct{ encCfg }
func (emptyEncCfg) GetEncryptKeyByName(string) string { return "" }

func TestEncodeDecode_SuccessAndFailure(t *testing.T) {
    inconfig.Read(encCfg{})
    data := []byte("hello")
    enc, err := EncodeResp(data, "k")
    if err != nil { t.Fatalf("encode err: %v", err) }
    dec, err := DecodeReq(enc, "k")
    if err != nil { t.Fatalf("decode err: %v", err) }
    if string(dec) != string(data) { t.Fatalf("roundtrip mismatch") }

    // failure when key empty
    inconfig.Read(emptyEncCfg{})
    if _, err := EncodeResp(data, "k"); err == nil { t.Fatalf("expected encode failure with empty key") }
    if _, err := DecodeReq(enc, "k"); err == nil { t.Fatalf("expected decode failure with empty key") }
}
