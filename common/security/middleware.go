package security

import (
	"net/http"
	"time"

	"github.com/Kingson4Wu/fast_proxy/common/auth"
)

// SecurityMiddleware provides comprehensive security middleware
type SecurityMiddleware struct {
	jwtAuth    *auth.JWTAuth
	keyManager *KeyManager
}

// NewSecurityMiddleware creates a new security middleware instance
func NewSecurityMiddleware(jwtSecret []byte, issuer string) *SecurityMiddleware {
	return &SecurityMiddleware{
		jwtAuth:    auth.NewJWTAuth(jwtSecret, issuer),
		keyManager: NewKeyManager(),
	}
}

// AuthMiddleware provides JWT authentication middleware
func (sm *SecurityMiddleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return sm.jwtAuth.AuthMiddleware(next)
}

// RateLimitMiddleware provides rate limiting middleware
func (sm *SecurityMiddleware) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// For demonstration, we'll implement a simple rate limiter
	// In a production environment, you might want to use a more sophisticated solution
	return func(w http.ResponseWriter, r *http.Request) {
		// Here you would implement rate limiting logic
		// For now, we'll just call the next handler
		next.ServeHTTP(w, r)
	}
}

// SecurityHeadersMiddleware adds security headers to responses
func (sm *SecurityMiddleware) SecurityHeadersMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		
		next.ServeHTTP(w, r)
	}
}

// CORSMiddleware handles CORS headers
func (sm *SecurityMiddleware) CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}

// CombinedMiddleware combines all security middleware
func (sm *SecurityMiddleware) CombinedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Apply middleware in order
	// Security headers first
	handler := sm.SecurityHeadersMiddleware(next)
	
	// Then CORS
	handler = sm.CORSMiddleware(handler)
	
	// Then rate limiting
	handler = sm.RateLimitMiddleware(handler)
	
	// Finally authentication
	handler = sm.AuthMiddleware(handler)
	
	return handler
}

// GenerateServiceToken generates a JWT token for a service
func (sm *SecurityMiddleware) GenerateServiceToken(serviceName string, permissions []string, duration time.Duration) (string, error) {
	expiresAt := time.Now().Add(duration)
	return sm.jwtAuth.GenerateToken(serviceName, permissions, expiresAt)
}

// GetKeyManager returns the key manager instance
func (sm *SecurityMiddleware) GetKeyManager() *KeyManager {
	return sm.keyManager
}

// RotateKey rotates a key in the key manager
func (sm *SecurityMiddleware) RotateKey(id, description string, size int) error {
	return sm.keyManager.RotateKey(id, description, size)
}

// GetActiveKey returns the active key from the key manager
func (sm *SecurityMiddleware) GetActiveKey(id string) (*KeyPair, bool) {
	return sm.keyManager.GetActiveKey(id)
}