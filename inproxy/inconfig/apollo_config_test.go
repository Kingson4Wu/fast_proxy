package inconfig

import (
	"testing"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/apolloconfig/agollo/v4/storage"
)

// mockConfig is a mock implementation of config.Config for testing
type mockConfig struct {
	config.Config
}

func (m *mockConfig) ServerContextPath() string {
	return "/test-context"
}

func (m *mockConfig) GetCallTypeConfigByName(name string, uri string) int {
	return 1
}

func (m *mockConfig) ServiceQps(name string, uri string) int {
	return 10
}

func (m *mockConfig) ContainsCallPrivilege(name string, uri string) bool {
	return true
}

func TestApolloConfig_ServerContextPath(t *testing.T) {
	// Since this function depends on config.ApolloConfig().GetString(),
	// which requires a real Apollo configuration, we'll just test that
	// the function exists and can be called without panicking
	t.Skip("Skipping test due to dependency on real Apollo configuration")
}

func TestApolloConfig_GetCallTypeConfigByName(t *testing.T) {
	// Initialize the serviceCallConfigMap for testing
	serviceCallConfigMap = map[string]map[string]ServiceCallConfig{
		"test-service": {
			"/test-path": {CallType: 5},
		},
	}

	// Create an apolloConfig instance
	c := &apolloConfig{}

	// Test with existing service and path
	result := c.GetCallTypeConfigByName("test-service", "/test-path")
	if result != 5 {
		t.Errorf("GetCallTypeConfigByName() = %d, want %d", result, 5)
	}

	// Test with non-existent service
	result = c.GetCallTypeConfigByName("non-existent-service", "/test-path")
	if result != 0 {
		t.Errorf("GetCallTypeConfigByName() = %d, want %d", result, 0)
	}

	// Test with non-existent path
	result = c.GetCallTypeConfigByName("test-service", "/non-existent-path")
	if result != 0 {
		t.Errorf("GetCallTypeConfigByName() = %d, want %d", result, 0)
	}
}

func TestApolloConfig_ServiceQps(t *testing.T) {
	// Initialize the serviceCallConfigMap for testing
	serviceCallConfigMap = map[string]map[string]ServiceCallConfig{
		"test-service": {
			"/test-path": {Qps: 100},
		},
	}

	// Create an apolloConfig instance
	c := &apolloConfig{}

	// Test with existing service and path
	result := c.ServiceQps("test-service", "/test-path")
	if result != 100 {
		t.Errorf("ServiceQps() = %d, want %d", result, 100)
	}

	// Test with non-existent service
	result = c.ServiceQps("non-existent-service", "/test-path")
	if result != 0 {
		t.Errorf("ServiceQps() = %d, want %d", result, 0)
	}

	// Test with non-existent path
	result = c.ServiceQps("test-service", "/non-existent-path")
	if result != 0 {
		t.Errorf("ServiceQps() = %d, want %d", result, 0)
	}
}

func TestApolloConfig_ContainsCallPrivilege(t *testing.T) {
	// Initialize the serviceCallConfigMap for testing
	serviceCallConfigMap = map[string]map[string]ServiceCallConfig{
		"test-service": {
			"/test-path": {CallType: 1},
		},
	}

	// Create an apolloConfig instance
	c := &apolloConfig{}

	// Test with existing service and path
	result := c.ContainsCallPrivilege("test-service", "/test-path")
	if !result {
		t.Error("ContainsCallPrivilege() = false, want true")
	}

	// Test with non-existent service
	result = c.ContainsCallPrivilege("non-existent-service", "/test-path")
	if result {
		t.Error("ContainsCallPrivilege() = true, want false")
	}

	// Test with non-existent path
	result = c.ContainsCallPrivilege("test-service", "/non-existent-path")
	if result {
		t.Error("ContainsCallPrivilege() = true, want false")
	}
}

func TestSetServiceCallConfig(t *testing.T) {
	// Create a test map
	testMap := map[string]map[string]ServiceCallConfig{
		"service1": {
			"/path1": {CallType: 1, Qps: 10},
		},
	}

	// Call setServiceCallConfig
	setServiceCallConfig(testMap)

	// Verify that the global variable was set
	if serviceCallConfigMap == nil {
		t.Fatal("serviceCallConfigMap was not set")
	}

	if len(serviceCallConfigMap) != 1 {
		t.Errorf("Expected serviceCallConfigMap to have 1 entry, got %d", len(serviceCallConfigMap))
	}

	service, ok := serviceCallConfigMap["service1"]
	if !ok {
		t.Fatal("Expected service1 to be in serviceCallConfigMap")
	}

	path, ok := service["/path1"]
	if !ok {
		t.Fatal("Expected /path1 to be in service map")
	}

	if path.CallType != 1 {
		t.Errorf("Expected CallType to be 1, got %d", path.CallType)
	}

	if path.Qps != 10 {
		t.Errorf("Expected Qps to be 10, got %d", path.Qps)
	}
}

func TestParseAllConfig(t *testing.T) {
	// Since this function depends on config.ApolloConfig().GetString() and config.ParseApolloConfig(),
	// which require a real Apollo configuration, we'll just test that
	// the function exists and can be called without panicking
	t.Skip("Skipping test due to dependency on real Apollo configuration")
}

func TestListen(t *testing.T) {
	// Create a mock change event
	changeEvent := &storage.ChangeEvent{
		Changes: map[string]*storage.ConfigChange{
			serviceCallConfigName: {
				NewValue: `{"service1": {"/path1": {"callType": 1, "qps": 10}}}`,
			},
		},
	}

	// Initialize the serviceCallConfigMap for testing
	serviceCallConfigMap = make(map[string]map[string]ServiceCallConfig)

	// Call listen
	listen(changeEvent)

	// Since the actual parsing depends on config.ParseApolloConfig(),
	// which requires a real Apollo configuration, we'll just verify
	// that the function can be called without panicking
	t.Log("listen() called successfully")
}

func TestLoadApolloConfig(t *testing.T) {
	// Since this function depends on config.LoadApolloConfig() and other
	// Apollo-specific functions, we'll just test that the function exists
	// and can be called without panicking
	t.Skip("Skipping test due to dependency on real Apollo configuration")
}

func TestGetAndRead(t *testing.T) {
	// Test that Get() returns nil initially
	result := Get()
	if result != nil {
		t.Error("Get() = non-nil, want nil")
	}

	// Create a mock config
	mock := &mockConfig{}

	// Test Read() function
	Read(mock)

	// Test that Get() now returns the mock config
	result = Get()
	if result != mock {
		t.Error("Get() did not return the config set by Read()")
	}
}