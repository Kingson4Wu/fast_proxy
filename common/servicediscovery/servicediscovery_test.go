package servicediscovery

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestBuilderAndGetters(t *testing.T) {
    sc := Create().
        AddressFunc(func(s string) *Address { return &Address{Ip: "1.2.3.4", Port: 8080} }).
        ClientNameFunc(func(r *http.Request) string { return "svc" }).
        RegisterFunc(func(name, ip string, port int) chan bool { return make(chan bool) }).
        Build()
    if sc == nil || sc.Address("") == nil || sc.Register("", "", 0) == nil {
        t.Fatalf("service center not built correctly")
    }
}

func TestGetRequestDeadTime(t *testing.T) {
    r := httptest.NewRequest("GET", "/", nil)
    if GetRequestDeadTime(r) != 0 {
        t.Fatalf("expected 0 when header missing")
    }
    r.Header.Set("request_dead_time", "abc")
    if GetRequestDeadTime(r) != 0 {
        t.Fatalf("expected 0 for invalid integer")
    }
    r.Header.Set("request_dead_time", "123")
    if GetRequestDeadTime(r) != 123 {
        t.Fatalf("expected 123")
    }
}
