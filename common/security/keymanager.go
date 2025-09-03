package security

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// KeyPair represents a cryptographic key with metadata
type KeyPair struct {
	ID          string
	PublicKey   []byte
	PrivateKey  []byte
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Active      bool
	Description string
}

// KeyManager handles secure storage and rotation of cryptographic keys
type KeyManager struct {
	mu           sync.RWMutex
	keys         map[string]*KeyPair
	keyRotation  time.Duration
	notification chan struct{}
}

// NewKeyManager creates a new key manager instance
func NewKeyManager() *KeyManager {
	km := &KeyManager{
		keys:         make(map[string]*KeyPair),
		keyRotation:  24 * time.Hour, // Default rotation period
		notification: make(chan struct{}, 1),
	}

	// Generate initial keys
	km.generateInitialKeys()

	// Start key rotation monitoring
	go km.monitorKeyRotation()

	return km
}

// generateInitialKeys generates the initial set of keys
func (km *KeyManager) generateInitialKeys() {
	// Generate a default encryption key
	encryptionKey, _ := km.GenerateKey("encryption-key", "Default encryption key", 32)
	encryptionKey.Active = true
	km.keys[encryptionKey.ID] = encryptionKey

	// Generate a default signature key
	signatureKey, _ := km.GenerateKey("signature-key", "Default signature key", 32)
	signatureKey.Active = true
	km.keys[signatureKey.ID] = signatureKey
}

// GenerateKey generates a new cryptographic key
func (km *KeyManager) GenerateKey(id, description string, size int) (*KeyPair, error) {
	// Generate random key material
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}

	// Create key pair
	keyPair := &KeyPair{
		ID:          id,
		PrivateKey:  key,
		PublicKey:   key[:size/2], // For demonstration, normally would derive public key differently
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(365 * 24 * time.Hour), // 1 year default
		Description: description,
	}

	return keyPair, nil
}

// GetActiveKey returns the active key for a given ID
func (km *KeyManager) GetActiveKey(id string) (*KeyPair, bool) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	key, exists := km.keys[id]
	if !exists || !key.Active || time.Now().After(key.ExpiresAt) {
		return nil, false
	}

	return key, true
}

// RotateKey rotates a key by generating a new one and deactivating the old one
func (km *KeyManager) RotateKey(id, description string, size int) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	// Mark current key as inactive
	if currentKey, exists := km.keys[id]; exists {
		currentKey.Active = false
	}

	// Generate new key
	newKey, err := km.GenerateKey(id, description, size)
	if err != nil {
		return fmt.Errorf("failed to generate new key: %w", err)
	}

	// Mark new key as active
	newKey.Active = true
	km.keys[id] = newKey

	// Notify about key rotation
	select {
	case km.notification <- struct{}{}:
	default:
	}

	return nil
}

// monitorKeyRotation monitors keys for automatic rotation
func (km *KeyManager) monitorKeyRotation() {
	ticker := time.NewTicker(km.keyRotation)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			km.checkAndRotateKeys()
		case <-km.notification:
			km.checkAndRotateKeys()
		}
	}
}

// checkAndRotateKeys checks all keys and rotates expired ones
func (km *KeyManager) checkAndRotateKeys() {
	km.mu.Lock()
	defer km.mu.Unlock()

	now := time.Now()
	for id, key := range km.keys {
		if now.After(key.ExpiresAt) {
			// Rotate expired key
			newKey, err := km.GenerateKey(id, key.Description+" (rotated)", len(key.PrivateKey))
			if err != nil {
				continue // Skip if we can't generate a new key
			}
			newKey.Active = true
			km.keys[id] = newKey
		}
	}
}

// DeriveKey derives a key using PBKDF2
func (km *KeyManager) DeriveKey(password, salt []byte, iterations, keyLen int) []byte {
	return pbkdf2.Key(password, salt, iterations, keyLen, sha256.New)
}

// GenerateSalt generates a random salt
func (km *KeyManager) GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// GetKeyInfo returns information about all keys
func (km *KeyManager) GetKeyInfo() map[string]map[string]interface{} {
	km.mu.RLock()
	defer km.mu.RUnlock()

	info := make(map[string]map[string]interface{})
	for id, key := range km.keys {
		info[id] = map[string]interface{}{
			"created_at":  key.CreatedAt,
			"expires_at":  key.ExpiresAt,
			"active":      key.Active,
			"description": key.Description,
		}
	}
	return info
}