package config

import (
    "testing"
    capp "github.com/Kingson4Wu/fast_proxy/common/apollo"
    "github.com/Kingson4Wu/fast_proxy/common/logger"
)

func TestGetConfigString_PanicWhenMissing(t *testing.T) {
    // configure a minimal ApolloConfig with default GetString returning empty
    config = &capp.ApolloConfig{Namespace: "application"}
    log = logger.DiscardLogger
    defer func() {
        if r := recover(); r == nil {
            t.Fatalf("expected panic when config string missing")
        }
    }()
    _ = getConfigString("nonexistent.key")
}

