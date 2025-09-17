package config

import "testing"

func TestApolloGetters_EmptyMapsReturnZeroValues(t *testing.T) {
    var c apolloConfig
    if c.GetServiceConfig("none") != nil {
        t.Fatalf("expected nil service config")
    }
    if c.GetEncryptKeyByName("none") != "" || c.GetSignKeyByName("none") != "" {
        t.Fatalf("expected empty keys for unknown names")
    }
    if c.GetTimeoutConfigByName("svc", "/nope") != 0 {
        t.Fatalf("expected 0 timeout for unknown path")
    }
}

