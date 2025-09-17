package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/Kingson4Wu/fast_proxy/common/config"
)

// mockConfig implements config.Config for testing
type mockConfig struct {
	httpClientMaxIdleConns        int
	httpClientMaxIdleConnsPerHost int
}

func (m mockConfig) ServerName() string { return "test" }
func (m mockConfig) ServerPort() int { return 8080 }
func (m mockConfig) ServiceRpcHeaderName() string { return "C_Service" }
func (m mockConfig) GetServiceConfig(string) *config.ServiceConfig { return &config.ServiceConfig{} }
func (m mockConfig) GetSignKey(*config.ServiceConfig) string { return "test" }
func (m mockConfig) GetSignKeyByName(string) string { return "test" }
func (m mockConfig) GetEncryptKey(*config.ServiceConfig) string { return "test" }
func (m mockConfig) GetEncryptKeyByName(string) string { return "test" }
func (m mockConfig) GetTimeoutConfigByName(string, string) int { return 1000 }
func (m mockConfig) HttpClientMaxIdleConns() int { return m.httpClientMaxIdleConns }
func (m mockConfig) HttpClientMaxIdleConnsPerHost() int { return m.httpClientMaxIdleConnsPerHost }
func (m mockConfig) FastHttpEnable() bool { return false }

// TestBuildClient tests the BuildClient function
func TestBuildClient(t *testing.T) {
	// Create a mock config
	cfg := mockConfig{
		httpClientMaxIdleConns:        10,
		httpClientMaxIdleConnsPerHost: 5,
	}

	BuildClient(cfg)

	if client == nil {
		t.Error("BuildClient() should initialize http client")
	}

	if fastHttpClient == nil {
		t.Error("BuildClient() should initialize fasthttp client")
	}
}

// TestWriteErrorMessage tests the writeErrorMessage function
func TestWriteErrorMessage(t *testing.T) {
	rr := httptest.NewRecorder()
	statusCode := http.StatusForbidden
	errorHeader := "test error message"

	writeErrorMessage(rr, statusCode, errorHeader)

	if rr.Code != statusCode {
		t.Errorf("Expected status code %d, got %d", statusCode, rr.Code)
	}

	if rr.Header().Get("proxy_error_message") != errorHeader {
		t.Errorf("Expected proxy_error_message header %s, got %s", errorHeader, rr.Header().Get("proxy_error_message"))
	}

	if rr.Header().Get("proxy_name") != "in_proxy" {
		t.Errorf("Expected proxy_name header 'in_proxy', got %s", rr.Header().Get("proxy_name"))
	}
}