package outconfig

import (
    "testing"
)

const sample = `
application:
  name: outproxy
  port: 18080
rpc:
  serviceHeaderName: C_Service
proxy:
  forwardAddress: http://127.0.0.1:9000/api
serviceConfig:
  svc:
    encryptKeyName: enc
encryptKeyConfig:
  enc: 0123456789ABCDEF
`

func TestLoadYamlConfig_ForwardAddress(t *testing.T) {
    c := LoadYamlConfig([]byte(sample)).(*yamlConfig)
    Read(c)
    if Get().ForwardAddress() != "http://127.0.0.1:9000/api" {
        t.Fatalf("forward address mismatch: %s", Get().ForwardAddress())
    }
    if Get().ServerName() != "outproxy" || Get().ServiceRpcHeaderName() != "C_Service" {
        t.Fatalf("unexpected base config getters")
    }
}

func TestReadAndGet(t *testing.T) {
    Read(nil)
    if Get() != nil {
        t.Fatalf("expected nil after Read(nil)")
    }
}

