package pack

import (
    "errors"
    "testing"

    "github.com/Kingson4Wu/fast_proxy/common/config"
)

func TestApplyMiddlewares_Error(t *testing.T) {
    sc := &config.ServiceConfig{}
    m1 := func(b []byte, _ *config.ServiceConfig) ([]byte, error) { return append(b, 'x'), nil }
    m2 := func(_ []byte, _ *config.ServiceConfig) ([]byte, error) { return nil, errors.New("boom") }
    if _, err := ApplyMiddlewares([]byte("a"), sc, m1, m2); err == nil {
        t.Fatalf("expected error from middleware chain")
    }
}

