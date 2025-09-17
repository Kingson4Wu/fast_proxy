package proxy

import (
    "bytes"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
    "github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
    "google.golang.org/protobuf/proto"
)

func TestWriteErrorMessage(t *testing.T) {
    rr := httptest.NewRecorder()
    writeErrorMessage(rr, http.StatusTeapot, "boom")
    if rr.Code != http.StatusTeapot || rr.Header().Get("proxy_error_message") != "boom" {
        t.Fatalf("unexpected error write: %d %q", rr.Code, rr.Header().Get("proxy_error_message"))
    }
}

// Use Delegate against an upstream that returns a valid ProxyRespData
func TestDelegate_OK(t *testing.T) {
    startServerForTest()
    waitServerReady(t)

    upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // respond with ProxyRespData containing "ok"
        pr := &protobuf.ProxyRespData{Payload: []byte("ok")}
        b, _ := proto.Marshal(pr)
        w.WriteHeader(http.StatusOK)
        io.Copy(w, bytes.NewReader(b))
    }))
    defer upstream.Close()

    // Override forward address to our upstream
    type cfg struct{ forwardCfg }
    c := &cfg{forwardCfg{url: upstream.URL}}
    outconfig.Read(c)

    r := httptest.NewRequest(http.MethodGet, "/p", nil)
    r.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()

    Delegate(rr, r)
    if rr.Code != http.StatusOK || rr.Body.String() != "ok" {
        t.Fatalf("unexpected delegate resp: %d %q", rr.Code, rr.Body.String())
    }
}

func TestDelegate_DecodeError(t *testing.T) {
    startServerForTest()
    waitServerReady(t)
    upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        io.WriteString(w, "bad")
    }))
    defer upstream.Close()
    type cfg struct{ forwardCfg }
    c := &cfg{forwardCfg{url: upstream.URL}}
    outconfig.Read(c)
    r := httptest.NewRequest(http.MethodGet, "/p", nil)
    r.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()
    Delegate(rr, r)
    if rr.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", rr.Code)
    }
    if rr.Header().Get("proxy_error_message") == "" {
        t.Fatalf("expected proxy_error_message header")
    }
}
