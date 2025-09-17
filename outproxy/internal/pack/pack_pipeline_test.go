package pack

import (
    "testing"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
    "github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
    "github.com/Kingson4Wu/fast_proxy/common/encrypt"
    "github.com/Kingson4Wu/fast_proxy/common/compress"
    "google.golang.org/protobuf/proto"
)

type mockOutCfg struct{}

func (m *mockOutCfg) ServerName() string                      { return "out_proxy" }
func (m *mockOutCfg) ServerPort() int                          { return 0 }
func (m *mockOutCfg) ServiceRpcHeaderName() string             { return "C_ServiceName" }
func (m *mockOutCfg) GetServiceConfig(serviceName string) *config.ServiceConfig {
    switch serviceName {
    case "plain":
        return &config.ServiceConfig{}
    case "enc":
        return &config.ServiceConfig{EncryptEnable: true, EncryptKeyName: "k"}
    case "zip":
        return &config.ServiceConfig{CompressEnable: true, CompressAlgorithm: 0}
    case "enczip":
        return &config.ServiceConfig{EncryptEnable: true, EncryptKeyName: "k", CompressEnable: true, CompressAlgorithm: 0}
    default:
        return nil
    }
}
func (m *mockOutCfg) GetSignKey(*config.ServiceConfig) string  { return "" }
func (m *mockOutCfg) GetSignKeyByName(string) string           { return "" }
func (m *mockOutCfg) GetEncryptKey(*config.ServiceConfig) string { return "ABCDABCDABCDABCD" }
func (m *mockOutCfg) GetEncryptKeyByName(string) string        { return "ABCDABCDABCDABCD" }
func (m *mockOutCfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (m *mockOutCfg) HttpClientMaxIdleConns() int              { return 10 }
func (m *mockOutCfg) HttpClientMaxIdleConnsPerHost() int       { return 10 }
func (m *mockOutCfg) FastHttpEnable() bool                     { return false }
func (m *mockOutCfg) ForwardAddress() string                   { return "http://127.0.0.1" }

func TestEncode_Combinations(t *testing.T) {
    outconfig.Read(&mockOutCfg{})

    cases := []string{"plain", "enc", "zip", "enczip"}
    for _, svc := range cases {
        if _, err := Encode([]byte("hello"), svc); err != nil {
            t.Fatalf("encode failed for %s: %v", svc, err)
        }
    }
}

func TestDecode_Combinations(t *testing.T) {
    outconfig.Read(&mockOutCfg{})
    plaintext := []byte("world")

    // helper to build response payload
    build := func(encEnabled, zipEnabled bool) []byte {
        data := plaintext
        var err error
        if encEnabled {
            data, err = encrypt.Encode(data, "ABCDABCDABCDABCD")
            if err != nil { t.Fatalf("encrypt: %v", err) }
        }
        if zipEnabled {
            data, err = compress.Encode(data, 0)
            if err != nil { t.Fatalf("compress: %v", err) }
        }
        pr := &protobuf.ProxyRespData{EncryptEnable: encEnabled, EncryptKeyName: "k", Compress: zipEnabled, CompressAlgorithm: 0, Payload: data}
        b, err := proto.Marshal(pr)
        if err != nil { t.Fatalf("marshal: %v", err) }
        return b
    }

    for _, tt := range []struct{ enc, zip bool }{{false,false},{true,false},{false,true},{true,true}} {
        body, _ := Decode(build(tt.enc, tt.zip))
        if string(body) != string(plaintext) {
            t.Fatalf("roundtrip mismatch enc=%v zip=%v got=%q", tt.enc, tt.zip, string(body))
        }
    }
}

func TestDecode_BadProtobuf(t *testing.T) {
    if _, err := Decode([]byte("not-protobuf")); err == nil {
        t.Fatal("expected error for bad protobuf input")
    }
}

