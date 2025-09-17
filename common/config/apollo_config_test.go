package config

import "testing"

func TestParseApolloConfig_Success(t *testing.T) {
    json := `{"a":{"b":1}}`
    called := false
    ok := ParseApolloConfig[string, map[string]int](json, func(m map[string]map[string]int) {
        called = true
        if m["a"]["b"] != 1 {
            t.Fatalf("unexpected map: %+v", m)
        }
    }, "name")
    if !ok || !called {
        t.Fatalf("expected success and callback")
    }
}

func TestApolloConfig_GettersFromMaps(t *testing.T) {
    // seed internal maps
    setServiceConfig(map[string]ServiceConfig{
        "svc": {EncryptKeyName: "enc", SignKeyName: "sig", EncryptEnable: true, SignEnable: true, CompressEnable: true, CompressAlgorithm: 0},
    })
    setEncryptKey(map[string]string{"enc": "E"})
    setSignKey(map[string]string{"sig": "S"})
    setServiceTimeoutConfig(map[string]map[string]ServiceTimeoutConfig{
        "svc": {"/x": {Timeout: 123}},
    })

    var c apolloConfig
    if sc := c.GetServiceConfig("svc"); sc == nil || !sc.EncryptEnable || !sc.SignEnable || !sc.CompressEnable {
        t.Fatalf("service config not returned correctly: %+v", sc)
    }
    if c.GetEncryptKey(&ServiceConfig{EncryptKeyName: "enc"}) != "E" || c.GetEncryptKeyByName("enc") != "E" {
        t.Fatalf("encrypt key getter failed")
    }
    if c.GetSignKey(&ServiceConfig{SignKeyName: "sig"}) != "S" || c.GetSignKeyByName("sig") != "S" {
        t.Fatalf("sign key getter failed")
    }
    if c.GetTimeoutConfigByName("svc", "/x") != 123 {
        t.Fatalf("timeout getter failed")
    }
}

