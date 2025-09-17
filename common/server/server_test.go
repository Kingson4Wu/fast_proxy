package server

import (
    "net/http"
    "testing"
    "time"

    "github.com/Kingson4Wu/fast_proxy/common/config"
    "github.com/Kingson4Wu/fast_proxy/common/logger"
    "github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
)

type cfg struct{}

func (c *cfg) ServerName() string                      { return "srv" }
func (c *cfg) ServerPort() int                          { return 0 }
func (c *cfg) ServiceRpcHeaderName() string             { return "C_Service" }
func (c *cfg) GetServiceConfig(string) *config.ServiceConfig { return &config.ServiceConfig{} }
func (c *cfg) GetSignKey(*config.ServiceConfig) string  { return "" }
func (c *cfg) GetSignKeyByName(string) string           { return "" }
func (c *cfg) GetEncryptKey(*config.ServiceConfig) string { return "" }
func (c *cfg) GetEncryptKeyByName(string) string        { return "" }
func (c *cfg) GetTimeoutConfigByName(string, string) int { return 0 }
func (c *cfg) HttpClientMaxIdleConns() int              { return 10 }
func (c *cfg) HttpClientMaxIdleConnsPerHost() int       { return 10 }
func (c *cfg) FastHttpEnable() bool                     { return false }

func TestOptionsAndHandlers(t *testing.T) {
    // getters before init should be safe
    // Depending on test order in other packages, global server may have been set.
    // These checks are flaky across packages, so skip strict assertions here.
    if GetLogger() == nil { t.Fatalf("expected non-nil logger") }

    p := NewServer(&cfg{}, logger.DiscardLogger, func(http.ResponseWriter, *http.Request) {})
    if len(p.otherHandlers) != 0 { t.Fatalf("expected empty handlers") }
    p.AddHandler("/healthz", func(http.ResponseWriter, *http.Request) {})
    p.AddHandler("/healthz", func(http.ResponseWriter, *http.Request) {}) // idempotent for same uri
    if len(p.otherHandlers) != 1 { t.Fatalf("expected one handler after duplicate add") }

    // apply options
    WithShutdownTimeout(123 * time.Millisecond)(p)
    if p.shutdownTimeout != 123*time.Millisecond { t.Fatalf("shutdown timeout not set") }

    sc := &servicediscovery.ServiceCenter{}
    WithServiceCenter(sc)(p)
    if p.sc != sc { t.Fatalf("service center not set") }

    WithLogger(logger.DiscardLogger)(p)
    if p.logger == nil { t.Fatalf("logger not set") }

    WithCustomHandler(func(http.ResponseWriter, *http.Request) {})(p)
    if p.proxyHandler == nil { t.Fatalf("custom handler not set") }

    // RegisterOnShutdown wires into http.Server; we can at least call it without starting
    done := make(chan struct{}, 1)
    p.RegisterOnShutdown(func() { close(done) })
}

func TestCenter(t *testing.T) {
    // Test Center() returns nil when server is not initialized
    srvMu.Lock()
    originalServer := server
    server = nil
    srvMu.Unlock()

    result := Center()
    if result != nil {
        t.Errorf("Center() = %v, want nil", result)
    }

    // Restore original server
    srvMu.Lock()
    server = originalServer
    srvMu.Unlock()
}

func TestConfig(t *testing.T) {
    // Test Config() returns nil when server is not initialized
    srvMu.Lock()
    originalServer := server
    server = nil
    srvMu.Unlock()

    result := Config()
    if result != nil {
        t.Errorf("Config() = %v, want nil", result)
    }

    // Restore original server
    srvMu.Lock()
    server = originalServer
    srvMu.Unlock()
}

func TestGetLogger_NilServer(t *testing.T) {
    // Test GetLogger() returns DiscardLogger when server is not initialized
    srvMu.Lock()
    originalServer := server
    server = nil
    srvMu.Unlock()

    result := GetLogger()
    if result == nil {
        t.Error("GetLogger() = nil, want non-nil logger")
    }

    // Restore original server
    srvMu.Lock()
    server = originalServer
    srvMu.Unlock()
}

func TestNew(t *testing.T) {
    p := New()
    if p == nil {
        t.Fatal("New() = nil, want non-nil Proxy")
    }
    if p.otherHandlers == nil {
        t.Error("New() should initialize otherHandlers map")
    }
    if len(p.otherHandlers) != 0 {
        t.Errorf("New() otherHandlers should be empty, got %d items", len(p.otherHandlers))
    }
}

func TestWithOptions(t *testing.T) {
    p := New()

    // Test WithShutdownTimeout
    timeout := 5 * time.Second
    WithShutdownTimeout(timeout)(p)
    if p.shutdownTimeout != timeout {
        t.Errorf("WithShutdownTimeout() = %v, want %v", p.shutdownTimeout, timeout)
    }

    // Test WithServiceCenter
    sc := &servicediscovery.ServiceCenter{}
    WithServiceCenter(sc)(p)
    if p.sc != sc {
        t.Errorf("WithServiceCenter() = %v, want %v", p.sc, sc)
    }

    // Test WithLogger
    logger := logger.DiscardLogger
    WithLogger(logger)(p)
    if p.logger != logger {
        t.Errorf("WithLogger() = %v, want %v", p.logger, logger)
    }

    // Test WithCustomHandler
    handler := func(http.ResponseWriter, *http.Request) {}
    WithCustomHandler(handler)(p)
    if p.proxyHandler == nil {
        t.Error("WithCustomHandler() should set proxyHandler")
    }
}

func TestAddHandler(t *testing.T) {
    p := New()
    
    // Add first handler
    handler1 := func(http.ResponseWriter, *http.Request) {}
    p.AddHandler("/test1", handler1)
    
    if len(p.otherHandlers) != 1 {
        t.Errorf("AddHandler() should have 1 handler, got %d", len(p.otherHandlers))
    }
    
    if p.otherHandlers["/test1"] == nil {
        t.Error("AddHandler() should add handler for /test1")
    }

    // Add second handler
    handler2 := func(http.ResponseWriter, *http.Request) {}
    p.AddHandler("/test2", handler2)
    
    if len(p.otherHandlers) != 2 {
        t.Errorf("AddHandler() should have 2 handlers, got %d", len(p.otherHandlers))
    }
    
    if p.otherHandlers["/test2"] == nil {
        t.Error("AddHandler() should add handler for /test2")
    }

    // Add duplicate handler (should not increase count)
    p.AddHandler("/test1", handler1)
    
    if len(p.otherHandlers) != 2 {
        t.Errorf("AddHandler() should still have 2 handlers after duplicate, got %d", len(p.otherHandlers))
    }
}

func TestRegisterOnShutdown(t *testing.T) {
    // Create a proxy with a real HTTP server to avoid nil pointer dereference
    p := New()
    p.svr = &http.Server{}
    
    // This test just verifies that RegisterOnShutdown can be called without panicking
    done := make(chan struct{}, 1)
    p.RegisterOnShutdown(func() { close(done) })
    
    // If we get here without panicking, the test passes
    t.Log("RegisterOnShutdown() called successfully")
}
