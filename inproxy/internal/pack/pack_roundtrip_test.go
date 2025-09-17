package pack

import (
    "bytes"
    "io"
    "net/http"
    "testing"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/compress"
    "github.com/Kingson4Wu/fast_proxy/common/encrypt"
    "github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
    "github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
    "google.golang.org/protobuf/proto"
)

// minimal mock config for inproxy encrypt key lookup
type mockInCfg struct{}

// inconfig.Config requires embedding of common/config.Config; implement full set
func (m *mockInCfg) ServerName() string                       { return "in_proxy" }
func (m *mockInCfg) ServerPort() int                           { return 0 }
func (m *mockInCfg) ServiceRpcHeaderName() string              { return "C_Service" }
func (m *mockInCfg) GetServiceConfig(string) *config.ServiceConfig {
    return &config.ServiceConfig{}
}
func (m *mockInCfg) GetSignKey(*config.ServiceConfig) string   { return "" }
func (m *mockInCfg) GetSignKeyByName(string) string            { return "" }
func (m *mockInCfg) GetEncryptKey(*config.ServiceConfig) string { return "" }
func (m *mockInCfg) GetEncryptKeyByName(string) string         { return "ABCDABCDABCDABCD" }
func (m *mockInCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (m *mockInCfg) HttpClientMaxIdleConns() int               { return 10 }
func (m *mockInCfg) HttpClientMaxIdleConnsPerHost() int        { return 10 }
func (m *mockInCfg) FastHttpEnable() bool                      { return true }
// additional inproxy.Config methods
func (m *mockInCfg) ServerContextPath() string                 { return "/" }
func (m *mockInCfg) GetCallTypeConfigByName(string, string) int { return 0 }
func (m *mockInCfg) ServiceQps(string, string) int             { return 0 }
func (m *mockInCfg) ContainsCallPrivilege(string, string) bool { return false }

// emptyKeyCfg returns empty encrypt key to force failures
type emptyKeyCfg struct{}

func (m *emptyKeyCfg) ServerName() string                       { return "in_proxy" }
func (m *emptyKeyCfg) ServerPort() int                           { return 0 }
func (m *emptyKeyCfg) ServiceRpcHeaderName() string              { return "C_Service" }
func (m *emptyKeyCfg) GetServiceConfig(string) *config.ServiceConfig { return &config.ServiceConfig{} }
func (m *emptyKeyCfg) GetSignKey(*config.ServiceConfig) string   { return "" }
func (m *emptyKeyCfg) GetSignKeyByName(string) string            { return "" }
func (m *emptyKeyCfg) GetEncryptKey(*config.ServiceConfig) string { return "" }
func (m *emptyKeyCfg) GetEncryptKeyByName(string) string         { return "" }
func (m *emptyKeyCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (m *emptyKeyCfg) HttpClientMaxIdleConns() int               { return 10 }
func (m *emptyKeyCfg) HttpClientMaxIdleConnsPerHost() int        { return 10 }
func (m *emptyKeyCfg) FastHttpEnable() bool                      { return true }
func (m *emptyKeyCfg) ServerContextPath() string                 { return "/" }
func (m *emptyKeyCfg) GetCallTypeConfigByName(string, string) int { return 0 }
func (m *emptyKeyCfg) ServiceQps(string, string) int             { return 0 }
func (m *emptyKeyCfg) ContainsCallPrivilege(string, string) bool { return false }

// helper to construct ProxyData (request) for Decode tests
func makeProxyDataPayload(body []byte, enc, zip bool) []byte {
    data := body
    var err error
    if enc {
        data, err = encrypt.Encode(data, "ABCDABCDABCDABCD")
        if err != nil { panic(err) }
    }
    if zip {
        data, err = compress.Encode(data, 0)
        if err != nil { panic(err) }
    }
    pd := &protobuf.ProxyData{
        Payload:           data,
        EncryptEnable:     enc,
        EncryptKeyName:    "k",
        Compress:          zip,
        CompressAlgorithm: 0,
        SignEnable:        false,
    }
    b, _ := proto.Marshal(pd)
    return b
}

func TestDecode_Combinations_Inproxy(t *testing.T) {
    inconfig.Read(&mockInCfg{})
    plain := []byte("inproxy")
    for _, c := range []struct{ enc, zip bool }{{false,false},{true,false},{false,true},{true,true}} {
        got, pd, err := Decode(makeProxyDataPayload(plain, c.enc, c.zip))
        if err != nil || pd == nil { t.Fatalf("decode err=%v pd=%v", err, pd) }
        if string(got) != string(plain) { t.Fatalf("mismatch: %q", string(got)) }
    }
}

func TestEncodeResp_And_DecodeReq(t *testing.T) {
    inconfig.Read(&mockInCfg{})
    // EncodeResp with encryption enabled (ensures payload present)
    resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString("ok"))}
    reData := &protobuf.ProxyData{EncryptEnable: true, EncryptKeyName: "k", Compress: false}
    b, err := EncodeResp(resp, reData)
    if err != nil { t.Fatalf("EncodeResp err: %v", err) }
    pr := new(protobuf.ProxyRespData)
    if err := proto.Unmarshal(b, pr); err != nil { t.Fatalf("unmarshal: %v", err) }
    if !pr.EncryptEnable { t.Fatalf("expected encrypted resp data") }
    // decrypt and check
    dec, derr := encrypt.Decode(pr.Payload, "ABCDABCDABCDABCD")
    if derr != nil { t.Fatalf("decrypt: %v", derr) }
    if string(dec) != "ok" { t.Fatalf("unexpected decoded: %q", string(dec)) }

    // Now verify DecodeReq path using ProxyData body
    reqBody := makeProxyDataPayload([]byte("hi"), true, true)
    r := &http.Request{Body: io.NopCloser(bytes.NewReader(reqBody))}
    bb, _, e := DecodeReq(r)
    if e != nil { t.Fatalf("DecodeReq err: %v", e) }
    if string(bb) != "hi" { t.Fatalf("unexpected body: %q", string(bb)) }
}

func TestEncode_EncryptFailure(t *testing.T) {
    // config with empty key forces Encrypt failure
    inconfig.Read(&emptyKeyCfg{})
    re := &protobuf.ProxyData{EncryptEnable: true, EncryptKeyName: "k"}
    if _, err := Encode([]byte("x"), re); err == nil {
        t.Fatalf("expected encrypt failure when key empty")
    }
}

func TestDecode_EncryptFailure(t *testing.T) {
    inconfig.Read(&mockInCfg{})
    // build ProxyData with encrypt, but our config returns empty key name -> failure
    b := makeProxyDataPayload([]byte("x"), true, false)
    // now force inconfig to return empty key name by using a cfg with GetEncryptKeyByName empty
    inconfig.Read(&emptyKeyCfg{})
    if _, _, err := Decode(b); err == nil {
        t.Fatalf("expected encrypt decode failure")
    }
}
