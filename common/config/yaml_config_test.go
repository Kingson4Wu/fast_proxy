package config

import (
    "testing"
)

const sampleYAML = `
application:
  name: fastproxy
  port: 18080
rpc:
  serviceHeaderName: C_Service
serviceConfig:
  svc1:
    encryptKeyName: enc.k1
    signKeyName: sign.k1
    encryptEnable: true
    signEnable: true
    compressEnable: true
    compressAlgorithm: 0
signKeyConfig:
  sign.k1: ABCD
encryptKeyConfig:
  enc.k1: 0123456789ABCDEF
serviceTimeoutConfig:
  svc1:
    /path: 1500
httpClient:
  MaxIdleConns: 100
  MaxIdleConnsPerHost: 10
fastHttp:
  enable: true
`

func TestLoadYamlConfigAndGetters(t *testing.T) {
    c := LoadYamlConfig([]byte(sampleYAML))

    if got := c.ServerName(); got != "fastproxy" {
        t.Fatalf("ServerName=%s", got)
    }
    if got := c.ServerPort(); got != 18080 {
        t.Fatalf("ServerPort=%d", got)
    }
    if got := c.ServiceRpcHeaderName(); got != "C_Service" {
        t.Fatalf("ServiceRpcHeaderName=%s", got)
    }

    sc := c.GetServiceConfig("svc1")
    if sc == nil || !sc.EncryptEnable || !sc.SignEnable || !sc.CompressEnable || sc.CompressAlgorithm != 0 {
        t.Fatalf("unexpected service config: %+v", sc)
    }
    if c.GetSignKey(sc) != "ABCD" || c.GetSignKeyByName("sign.k1") != "ABCD" {
        t.Fatalf("sign keys mismatch")
    }
    if c.GetEncryptKey(sc) != "0123456789ABCDEF" || c.GetEncryptKeyByName("enc.k1") != "0123456789ABCDEF" {
        t.Fatalf("encrypt keys mismatch")
    }
    if c.GetTimeoutConfigByName("svc1", "/path") != 1500 {
        t.Fatalf("timeout mismatch")
    }
    if !c.FastHttpEnable() {
        t.Fatalf("expected FastHttpEnable true")
    }
    if c.HttpClientMaxIdleConns() != 100 || c.HttpClientMaxIdleConnsPerHost() != 10 {
        t.Fatalf("http client settings mismatch")
    }
}

func TestLoadYamlConfig_Error(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Fatalf("expected panic on invalid yaml")
        }
    }()
    _ = LoadYamlConfig([]byte(":::not-yaml"))
}

