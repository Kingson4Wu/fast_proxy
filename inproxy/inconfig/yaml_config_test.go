package inconfig

import (
    "testing"
)

const inSample = `
application:
  name: inproxy
  port: 18081
  contextPath: /api
rpc:
  serviceHeaderName: C_Service
serviceCallTypeConfig:
  svc:
    /path:
      callType: 2
      qps: 50
serviceConfig:
  svc:
    encryptKeyName: enc
encryptKeyConfig:
  enc: 0123456789ABCDEF
`

func TestInYamlConfig_Getters(t *testing.T) {
    c := LoadYamlConfig([]byte(inSample)).(*yamlConfig)
    Read(c) // set global
    if Get().ServerContextPath() != "/api" { t.Fatalf("context path") }
    if v := Get().GetCallTypeConfigByName("svc", "/path"); v != 2 { t.Fatalf("callType=%d", v) }
    if v := Get().ServiceQps("svc", "/path"); v != 50 { t.Fatalf("qps=%d", v) }
    if !Get().ContainsCallPrivilege("svc", "/path") { t.Fatalf("ContainsCallPrivilege false") }
    if Get().GetEncryptKeyByName("enc") == "" { t.Fatalf("encrypt key missing") }
}

