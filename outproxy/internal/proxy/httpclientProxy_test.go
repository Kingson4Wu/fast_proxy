package proxy

import (
    "bytes"
    "context"
    "errors"
    "io"
    "net/http"
    "net/http/httptest"
    "sync"
    "testing"
    "time"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/server"
    "github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
    "github.com/Kingson4Wu/fast_proxy/common/logger/zap"
    "github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
    "github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
    "google.golang.org/protobuf/proto"
)

// mock config with minimal functionality for proxy tests
type proxyMockConfig struct{}

func (c *proxyMockConfig) ServerName() string                            { return "out_proxy" }
func (c *proxyMockConfig) ServerPort() int                                { return 0 }
func (c *proxyMockConfig) ServiceRpcHeaderName() string                   { return "C_ServiceName" }
func (c *proxyMockConfig) GetServiceConfig(serviceName string) *config.ServiceConfig {
    // keep all features off to avoid crypt/sign requirements in pack.EncodeReq
    return &config.ServiceConfig{}
}
func (c *proxyMockConfig) GetSignKey(*config.ServiceConfig) string        { return "" }
func (c *proxyMockConfig) GetSignKeyByName(string) string                 { return "" }
func (c *proxyMockConfig) GetEncryptKey(*config.ServiceConfig) string     { return "" }
func (c *proxyMockConfig) GetEncryptKeyByName(string) string              { return "" }
func (c *proxyMockConfig) GetTimeoutConfigByName(string, string) int      { return 0 }
func (c *proxyMockConfig) HttpClientMaxIdleConns() int                    { return 10 }
func (c *proxyMockConfig) HttpClientMaxIdleConnsPerHost() int             { return 10 }
func (c *proxyMockConfig) FastHttpEnable() bool                           { return false }
func (c *proxyMockConfig) ForwardAddress() string                         { return "http://example.com" }

// cfg type allowing dynamic ForwardAddress override for tests
type forwardCfg struct {
    proxyMockConfig
    url string
}

func (c *forwardCfg) ForwardAddress() string { return c.url }

// timeoutCfg overrides per-route timeout
type timeoutCfg struct{ proxyMockConfig }

func (c *timeoutCfg) GetTimeoutConfigByName(name, uri string) int { return 10 }

// roundTripper that records the outbound request and returns a provided response sequence
type seqRT struct {
    calls   int
    reqs    []*http.Request
    replies []func(*http.Request) (*http.Response, error)
}

func (s *seqRT) RoundTrip(r *http.Request) (*http.Response, error) {
    s.calls++
    s.reqs = append(s.reqs, r.Clone(r.Context()))
    if s.calls-1 < len(s.replies) {
        return s.replies[s.calls-1](r)
    }
    return nil, errors.New("no more replies")
}

// helper to build a simple successful proxy response that pack.DecodeResp can decode directly
func makePlainProxyOK(body []byte) *http.Response {
    pr := &protobuf.ProxyRespData{
        Compress:          false,
        EncryptEnable:     false,
        EncryptKeyName:    "",
        CompressAlgorithm: 0,
        Payload:           body,
    }
    b, _ := proto.Marshal(pr)
    return &http.Response{
        StatusCode: http.StatusOK,
        Header:     http.Header{},
        Body:       io.NopCloser(bytes.NewReader(b)),
    }
}

type timeoutErr struct{}
func (timeoutErr) Error() string   { return "timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

// Build a minimal ServiceCenter matching C_ServiceName header
func testSC() *servicediscovery.ServiceCenter {
    return servicediscovery.Create().
        AddressFunc(func(serviceName string) *servicediscovery.Address { return &servicediscovery.Address{Ip: "127.0.0.1", Port: 0} }).
        ClientNameFunc(func(r *http.Request) string { return r.Header.Get("C_ServiceName") }).
        RegisterFunc(func(string, string, int) chan bool { return make(chan bool) }).
        Build()
}

var startOnce sync.Once

func startServerForTest() {
    cfg := &proxyMockConfig{}
    outconfig.Read(cfg)
    startOnce.Do(func() {
        go func() {
            p := server.NewServer(cfg, zap.DefaultLogger(), func(http.ResponseWriter, *http.Request) {})
            p.Start(server.WithServiceCenter(testSC()))
        }()
    })
}

func waitServerReady(t *testing.T) {
    deadline := time.Now().Add(500 * time.Millisecond)
    for time.Now().Before(deadline) {
        if server.Center() != nil && server.Config() != nil {
            return
        }
        time.Sleep(5 * time.Millisecond)
    }
    t.Fatal("server not ready for test")
}

func TestDoProxy_HeaderFiltering(t *testing.T) {
    startServerForTest()
    waitServerReady(t)

    // stub client with one successful response; include hop-by-hop headers in response to ensure filtered
    rt := &seqRT{replies: []func(*http.Request) (*http.Response, error){
        func(r *http.Request) (*http.Response, error) {
            resp := makePlainProxyOK([]byte("ok"))
            resp.Header.Set("X-Upstream", "abc")
            resp.Header.Set("Connection", "close")
            resp.Header.Set("Transfer-Encoding", "chunked")
            return resp, nil
        },
    }}
    client = &http.Client{Transport: rt}

    // craft request with hop-by-hop headers
    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Proxy-Connection", "keep-alive")
    req.Header.Set("Keep-Alive", "timeout=5")
    req.Header.Set("X-Test", "1")

    rr := httptest.NewRecorder()
    DoProxy(rr, req)

    // request filtering
    if h := rt.reqs[0].Header; h.Get("X-Test") != "1" || h.Get("Connection") != "" || h.Get("Proxy-Connection") != "" {
        t.Fatalf("request headers not filtered correctly: %+v", h)
    }

    // response filtering
    if rr.Header().Get("X-Upstream") != "abc" {
        t.Fatalf("expected X-Upstream header to pass through")
    }
    if rr.Header().Get("Connection") != "" || rr.Header().Get("Transfer-Encoding") != "" {
        t.Fatalf("hop-by-hop headers must not be forwarded in response: %+v", rr.Header())
    }
    if got := rr.Body.String(); got != "ok" {
        t.Fatalf("unexpected body: %q", got)
    }
}

func TestDoProxy_RetrySequence(t *testing.T) {
    startServerForTest()
    waitServerReady(t)

    // set very small backoff via env handled in BuildClient already run in outproxy.NewServer; here we just stub client
    rt := &seqRT{replies: []func(*http.Request) (*http.Response, error){
        func(r *http.Request) (*http.Response, error) { return nil, timeoutErr{} },
        func(r *http.Request) (*http.Response, error) { return &http.Response{StatusCode: http.StatusServiceUnavailable, Body: io.NopCloser(bytes.NewReader(nil))}, nil },
        func(r *http.Request) (*http.Response, error) { return makePlainProxyOK([]byte("ok")), nil },
    }}
    client = &http.Client{Transport: rt}

    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()

    DoProxy(rr, req)

    if rt.calls != 3 {
        t.Fatalf("expected 3 attempts, got %d", rt.calls)
    }
    if rr.Code != http.StatusOK || rr.Body.String() != "ok" {
        t.Fatalf("unexpected response: code=%d body=%q", rr.Code, rr.Body.String())
    }
}

func TestDoProxy_NoRetryOnDeadlineExceeded(t *testing.T) {
    startServerForTest()
    waitServerReady(t)

    rt := &seqRT{replies: []func(*http.Request) (*http.Response, error){
        func(r *http.Request) (*http.Response, error) { return nil, context.DeadlineExceeded },
    }}
    client = &http.Client{Transport: rt}

    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()
    DoProxy(rr, req)

    if rt.calls != 1 {
        t.Fatalf("expected 1 attempt, got %d", rt.calls)
    }
    if rr.Code != http.StatusGatewayTimeout {
        t.Fatalf("expected 504, got %d", rr.Code)
    }
}

func TestDoProxy_DeadlineHeader_Past(t *testing.T) {
    startServerForTest()
    waitServerReady(t)

    rt := &seqRT{replies: []func(*http.Request) (*http.Response, error){
        func(r *http.Request) (*http.Response, error) { return makePlainProxyOK([]byte("ok")), nil },
    }}
    client = &http.Client{Transport: rt}

    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    // set past deadline
    req.Header.Set("request_dead_time", "1")
    rr := httptest.NewRecorder()
    DoProxy(rr, req)

    if rt.calls != 0 {
        t.Fatalf("expected 0 attempts when deadline already passed; got %d", rt.calls)
    }
    if rr.Code != http.StatusGatewayTimeout {
        t.Fatalf("expected 504, got %d", rr.Code)
    }
}

func TestDoProxy_DecodeErrorSetsStatus(t *testing.T) {
    startServerForTest()
    waitServerReady(t)
    rt := &seqRT{replies: []func(*http.Request) (*http.Response, error){
        func(r *http.Request) (*http.Response, error) {
            // invalid body for pack.DecodeResp
            return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString("bad")), Header: http.Header{}}, nil
        },
    }}
    client = &http.Client{Transport: rt}
    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()
    DoProxy(rr, req)
    if rr.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500 when decode fails, got %d", rr.Code)
    }
    if rr.Header().Get("proxy_error_message") == "" {
        t.Fatalf("expected proxy_error_message header")
    }
}

func TestDoProxy_ForwardNonOK(t *testing.T) {
    startServerForTest()
    waitServerReady(t)
    rt := &seqRT{replies: []func(*http.Request) (*http.Response, error){
        func(r *http.Request) (*http.Response, error) {
            resp := &http.Response{StatusCode: 418, Header: http.Header{}}
            resp.Body = io.NopCloser(bytes.NewBufferString("oops"))
            return resp, nil
        },
    }}
    client = &http.Client{Transport: rt}
    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()
    DoProxy(rr, req)
    if rr.Body.String() != "oops" {
        t.Fatalf("expected body passthrough for non-OK, got %q", rr.Body.String())
    }
}

func TestFastDoProxy_ForwardNonOK(t *testing.T) {
    startServerForTest()
    waitServerReady(t)
    upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(418)
        io.WriteString(w, "oops")
    }))
    defer upstream.Close()
    c := &forwardCfg{url: upstream.URL}
    c.proxyMockConfig = proxyMockConfig{}
    outconfig.Read(c)
    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()
    if fastHttpClient == nil { BuildClient(c) }
    fastDoProxy(rr, req)
    if rr.Body.String() != "oops" {
        t.Fatalf("expected body passthrough for non-OK, got %q", rr.Body.String())
    }
}

func TestFastDoProxy_HeaderFilteringAndBody(t *testing.T) {
    startServerForTest()
    waitServerReady(t)

    // upstream server returns a ProxyRespData wrapping "ok"
    upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        resp := makePlainProxyOK([]byte("ok"))
        for k, v := range resp.Header {
            w.Header()[k] = v
        }
        w.WriteHeader(resp.StatusCode)
        io.Copy(w, resp.Body)
    }))
    defer upstream.Close()

    // set forward address to upstream server
    c := &forwardCfg{url: upstream.URL}
    c.proxyMockConfig = proxyMockConfig{}
    outconfig.Read(c)

    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()

    // ensure we have a client
    if fastHttpClient == nil {
        BuildClient(c)
    }

    // call fast path directly with valid body
    fastDoProxy(rr, req)

    if rr.Code != http.StatusOK || rr.Body.String() != "ok" {
        t.Fatalf("unexpected response: %d %q", rr.Code, rr.Body.String())
    }
}

func TestFastDoProxy_DecodeError(t *testing.T) {
    startServerForTest()
    waitServerReady(t)

    // upstream returns garbage body to trigger decode error
    upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        io.WriteString(w, "bad")
    }))
    defer upstream.Close()

    c := &forwardCfg{url: upstream.URL}
    c.proxyMockConfig = proxyMockConfig{}
    outconfig.Read(c)

    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    req.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()
    if fastHttpClient == nil { BuildClient(c) }
    fastDoProxy(rr, req)
    if rr.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500 for decode error, got %d", rr.Code)
    }
}

func TestRetryEnvAndClose(t *testing.T) {
    t.Setenv("FAST_PROXY_RETRY_ENABLED", "true")
    t.Setenv("FAST_PROXY_RETRY_MAX", "3")
    t.Setenv("FAST_PROXY_RETRY_BASE_MS", "5")
    t.Setenv("FAST_PROXY_RETRY_MAX_MS", "10")
    // Build client to initialize retry config
    BuildClient(&proxyMockConfig{})
    if !enableRetries || maxRetries != 3 {
        t.Fatalf("retry env not applied: enabled=%v max=%d", enableRetries, maxRetries)
    }
    Close() // ensure no panic
}

func TestDoProxy_PerRouteTimeoutAndNonTimeoutError(t *testing.T) {
    startServerForTest()
    waitServerReady(t)

    // Config with per-route timeout
    outconfig.Read(&timeoutCfg{})

    // RoundTripper that checks deadline and returns a non-timeout error
    sawDeadline := false
    rt := &seqRT{replies: []func(*http.Request) (*http.Response, error){
        func(r *http.Request) (*http.Response, error) {
            if dl, ok := r.Context().Deadline(); ok && time.Until(dl) > 0 {
                sawDeadline = true
            }
            return nil, errors.New("boom")
        },
    }}
    client = &http.Client{Transport: rt}

    req := httptest.NewRequest(http.MethodGet, "/will-timeout", nil)
    req.Header.Set("C_ServiceName", "svc")
    rr := httptest.NewRecorder()
    DoProxy(rr, req)
    if !sawDeadline {
        t.Fatalf("expected per-route deadline to be set on request context")
    }
    if rt.calls != 1 || rr.Code != http.StatusServiceUnavailable {
        t.Fatalf("expected 1 attempt and 503 for non-timeout error; calls=%d code=%d", rt.calls, rr.Code)
    }
}
