package config

import (
    "testing"
    "github.com/Kingson4Wu/fast_proxy/common/logger"
)

// fake reader implementing ConfigReader for success path
type fakeReader struct {
    s map[string]string
    i map[string]int
    b map[string]bool
}

func (f fakeReader) GetString(name string) string { return f.s[name] }
func (f fakeReader) GetIntValue(name string, def int) int {
    if v, ok := f.i[name]; ok { return v }
    return def
}
func (f fakeReader) GetBoolValue(name string, def bool) bool {
    if v, ok := f.b[name]; ok { return v }
    return def
}

func TestParseAllConfig_WithFakeReader_Success(t *testing.T) {
    // Seed JSON configs and base fields required by getters
    svcJSON := `{"svc":{"encrypt.key.name":"enc","sign.key.name":"sig","encrypt.enable":true,"sign.enable":true,"compress.enable":true,"compress.algorithm":0}}`
    signJSON := `{"sig":"S-KEY"}`
    encJSON := `{"enc":"E-KEY"}`
    toutJSON := `{"svc":{"/x":{"timeout":321}}}`

    cfgReader = fakeReader{
        s: map[string]string{
            serviceConfigName:        svcJSON,
            signKeyConfigName:        signJSON,
            encryptKeyConfigName:     encJSON,
            serviceTimeoutConfigName: toutJSON,
            "application.name":       "fastproxy",
            "rpc.serviceHeaderName":  "C_Service",
        },
        i: map[string]int{
            "application.port":              18080,
            "httpClient.MaxIdleConnsPerHost": 42,
        },
        b: map[string]bool{
            "fastHttp.Enable": true,
        },
    }
    log = logger.DiscardLogger

    // parse config into internal maps
    parseAllConfig()

    // Validate getters
    var c apolloConfig
    if c.ServerName() != "fastproxy" || c.ServerPort() != 18080 {
        t.Fatalf("unexpected app fields: name=%s port=%d", c.ServerName(), c.ServerPort())
    }
    if c.ServiceRpcHeaderName() != "C_Service" {
        t.Fatalf("rpc header name mismatch: %s", c.ServiceRpcHeaderName())
    }
    if !c.FastHttpEnable() || c.HttpClientMaxIdleConnsPerHost() != 42 {
        t.Fatalf("http settings mismatch")
    }
    sc := c.GetServiceConfig("svc")
    if sc == nil || !sc.EncryptEnable || !sc.SignEnable || !sc.CompressEnable {
        t.Fatalf("service config not parsed: %+v", sc)
    }
    if c.GetEncryptKeyByName("enc") != "E-KEY" || c.GetSignKeyByName("sig") != "S-KEY" {
        t.Fatalf("keys not parsed correctly")
    }
    if c.GetTimeoutConfigByName("svc", "/x") != 321 {
        t.Fatalf("timeout not parsed")
    }
}

