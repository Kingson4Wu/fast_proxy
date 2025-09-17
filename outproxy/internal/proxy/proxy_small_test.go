package proxy

import (
    "net/http"
    "testing"
)

func TestGetProxyAndHandle(t *testing.T) {
    if GetProxy(FORWARD) == nil || GetProxy(REVERSE) == nil {
        t.Fatalf("expected non-nil proxy funcs")
    }
    called := 0
    f := Func(func(http.ResponseWriter, *http.Request) { called++ })
    f.Handle(nil, nil)
    if called != 1 {
        t.Fatalf("Func.Handle did not call underlying function")
    }
}

// Note: wrappers ReverseProxy/ForwardProxy defer to real proxying paths which are covered elsewhere.
