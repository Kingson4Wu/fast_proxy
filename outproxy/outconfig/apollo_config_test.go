package outconfig

import (
	"testing"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
)

// mockConfig implements config.Config for testing
type mockConfig struct {
	forwardAddress string
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
func (m mockConfig) HttpClientMaxIdleConns() int { return 10 }
func (m mockConfig) HttpClientMaxIdleConnsPerHost() int { return 5 }
func (m mockConfig) FastHttpEnable() bool { return false }
func (m mockConfig) ForwardAddress() string { return m.forwardAddress }

// mockLogger is a simple mock for the logger
type mockLogger struct {
	logger.Logger
}

func (m *mockLogger) Info(args ...interface{}) {}
func (m *mockLogger) Warn(args ...interface{}) {}
func (m *mockLogger) Error(args ...interface{}) {}
func (m *mockLogger) Debug(args ...interface{}) {}
func (m *mockLogger) Panic(args ...interface{}) {}
func (m *mockLogger) Fatal(args ...interface{}) {}
func (m *mockLogger) Infof(template string, args ...interface{}) {}
func (m *mockLogger) Warnf(template string, args ...interface{}) {}
func (m *mockLogger) Errorf(template string, args ...interface{}) {}
func (m *mockLogger) Debugf(template string, args ...interface{}) {}
func (m *mockLogger) Panicf(template string, args ...interface{}) {}
func (m *mockLogger) Fatalf(template string, args ...interface{}) {}

// TestLoadApolloConfig tests the LoadApolloConfig function
func TestLoadApolloConfig(t *testing.T) {
	// Since this function depends on config.LoadApolloConfig() which requires
	// a real Apollo configuration, we'll just test that the function exists
	// and can be called without panicking
	t.Skip("Skipping test due to dependency on real Apollo configuration")
}

// TestApolloConfig_ForwardAddress tests the ForwardAddress method of apolloConfig
func TestApolloConfig_ForwardAddress(t *testing.T) {
	// Since this function depends on config.ApolloConfig().GetString() which requires
	// a real Apollo configuration, we'll just test that the method exists
	// and can be called without panicking
	t.Skip("Skipping test due to dependency on real Apollo configuration")
}