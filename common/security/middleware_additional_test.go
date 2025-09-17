package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	issuer := "fastproxy-test"
	sm := NewSecurityMiddleware(jwtSecret, issuer)

	// Test that AuthMiddleware can be called without panicking
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := sm.AuthMiddleware(handler)
	if middleware == nil {
		t.Error("AuthMiddleware() = nil, want non-nil handler")
	}

	// Test that the middleware can be called
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	// Note: This test doesn't validate authentication logic as that would require a valid token
	// In a real test, you would set up a valid JWT token in the request
}

func TestRateLimitMiddleware(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	issuer := "fastproxy-test"
	sm := NewSecurityMiddleware(jwtSecret, issuer)

	// Test that RateLimitMiddleware can be called without panicking
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := sm.RateLimitMiddleware(handler)
	if middleware == nil {
		t.Error("RateLimitMiddleware() = nil, want non-nil handler")
	}

	// Test that the middleware can be called
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	// Verify that the handler was called (status should be OK)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCombinedMiddleware(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	issuer := "fastproxy-test"
	sm := NewSecurityMiddleware(jwtSecret, issuer)

	// Test that CombinedMiddleware can be called without panicking
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := sm.CombinedMiddleware(handler)
	if middleware == nil {
		t.Error("CombinedMiddleware() = nil, want non-nil handler")
	}

	// Test that the middleware can be called
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	// Note: We can't verify the status code here because the AuthMiddleware
	// will likely return 401 since we don't have a valid token
	
	// Verify that security headers were added
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header to be set")
	}
}

func TestGetKeyManager(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	issuer := "fastproxy-test"
	sm := NewSecurityMiddleware(jwtSecret, issuer)

	// Test that GetKeyManager returns the key manager instance
	keyManager := sm.GetKeyManager()
	if keyManager == nil {
		t.Error("GetKeyManager() = nil, want non-nil key manager")
	}
}

func TestNewSecurityMiddleware(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	issuer := "fastproxy-test"

	// Test that NewSecurityMiddleware creates a new instance successfully
	sm := NewSecurityMiddleware(jwtSecret, issuer)
	if sm == nil {
		t.Fatal("NewSecurityMiddleware() = nil, want non-nil SecurityMiddleware")
	}

	if sm.jwtAuth == nil {
		t.Error("NewSecurityMiddleware() should initialize jwtAuth")
	}

	if sm.keyManager == nil {
		t.Error("NewSecurityMiddleware() should initialize keyManager")
	}
}

func TestGenerateServiceToken(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	issuer := "fastproxy-test"
	sm := NewSecurityMiddleware(jwtSecret, issuer)

	serviceName := "test-service"
	permissions := []string{"read", "write"}
	duration := time.Hour

	// Test that GenerateServiceToken generates a token without error
	token, err := sm.GenerateServiceToken(serviceName, permissions, duration)
	if err != nil {
		t.Fatalf("GenerateServiceToken() error = %v, want nil", err)
	}

	if token == "" {
		t.Error("GenerateServiceToken() = empty string, want non-empty token")
	}
}

func TestRotateKeyAndRotateKey(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	issuer := "fastproxy-test"
	sm := NewSecurityMiddleware(jwtSecret, issuer)

	// Test that RotateKey can be called without error
	err := sm.RotateKey("test-key", "Test key", 32)
	if err != nil {
		t.Fatalf("RotateKey() error = %v, want nil", err)
	}

	// Test that GetActiveKey returns the key we just created
	key, exists := sm.GetActiveKey("test-key")
	if !exists {
		t.Error("GetActiveKey() exists = false, want true")
	}

	if key == nil {
		t.Error("GetActiveKey() = nil, want non-nil key")
	}
}

func TestGetActiveKey_NonExistentKey(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	issuer := "fastproxy-test"
	sm := NewSecurityMiddleware(jwtSecret, issuer)

	// Test that GetActiveKey returns false for non-existent keys
	key, exists := sm.GetActiveKey("non-existent-key")
	if exists {
		t.Error("GetActiveKey() exists = true, want false")
	}

	if key != nil {
		t.Error("GetActiveKey() = non-nil, want nil")
	}
}