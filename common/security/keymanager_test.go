package security

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeyManager(t *testing.T) {
	Convey("Given a KeyManager instance", t, func() {
		km := NewKeyManager()

		Convey("When creating a new KeyManager", func() {
			Convey("Should have initial keys", func() {
				encryptionKey, exists := km.GetActiveKey("encryption-key")
				So(exists, ShouldBeTrue)
				So(encryptionKey, ShouldNotBeNil)
				So(encryptionKey.Active, ShouldBeTrue)

				signatureKey, exists := km.GetActiveKey("signature-key")
				So(exists, ShouldBeTrue)
				So(signatureKey, ShouldNotBeNil)
				So(signatureKey.Active, ShouldBeTrue)
			})
		})

		Convey("When generating a new key", func() {
			key, err := km.GenerateKey("test-key", "Test key", 32)

			Convey("Should generate key without error", func() {
				So(err, ShouldBeNil)
				So(key, ShouldNotBeNil)
				So(key.ID, ShouldEqual, "test-key")
				So(key.Description, ShouldEqual, "Test key")
				So(len(key.PrivateKey), ShouldEqual, 32)
				So(key.Active, ShouldBeFalse) // New keys are not active by default
			})
		})

		Convey("When rotating a key", func() {
			err := km.RotateKey("test-rotate-key", "Test rotation key", 32)

			Convey("Should rotate key without error", func() {
				So(err, ShouldBeNil)

				key, exists := km.GetActiveKey("test-rotate-key")
				So(exists, ShouldBeTrue)
				So(key, ShouldNotBeNil)
				So(key.Active, ShouldBeTrue)
				So(key.Description, ShouldEqual, "Test rotation key")
			})
		})

		Convey("When getting key info", func() {
			info := km.GetKeyInfo()

			Convey("Should return key information", func() {
				So(info, ShouldNotBeNil)
				So(len(info), ShouldBeGreaterThanOrEqualTo, 2) // At least encryption and signature keys
			})
		})

		Convey("When deriving a key", func() {
			password := []byte("test-password")
			salt, err := km.GenerateSalt(16)
			So(err, ShouldBeNil)

			derivedKey := km.DeriveKey(password, salt, 10000, 32)

			Convey("Should derive key successfully", func() {
				So(derivedKey, ShouldNotBeNil)
				So(len(derivedKey), ShouldEqual, 32)
			})
		})

		Convey("When generating salt", func() {
			salt, err := km.GenerateSalt(16)

			Convey("Should generate salt without error", func() {
				So(err, ShouldBeNil)
				So(salt, ShouldNotBeNil)
				So(len(salt), ShouldEqual, 16)
			})
		})
	})
}

func TestKeyExpiration(t *testing.T) {
	Convey("Given a KeyManager with expiring keys", t, func() {
		km := NewKeyManager()

		// Manually create an expired key for testing
		km.mu.Lock()
		expiredKey := &KeyPair{
			ID:         "expired-key",
			PrivateKey: make([]byte, 32),
			CreatedAt:  time.Now().Add(-2 * time.Hour),
			ExpiresAt:  time.Now().Add(-1 * time.Hour), // Expired
			Active:     true,
		}
		km.keys["expired-key"] = expiredKey
		km.mu.Unlock()

		Convey("When checking for expired key", func() {
			_, exists := km.GetActiveKey("expired-key")

			Convey("Should not return expired key", func() {
				So(exists, ShouldBeFalse)
			})
		})
	})
}