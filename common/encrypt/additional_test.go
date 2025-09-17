package encrypt

import (
	"errors"
	"testing"
	"github.com/Kingson4Wu/fast_proxy/common/encrypt/aes"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/agiledragon/gomonkey/v2"
)

// mockLogger is a simple mock for the logger
type mockLogger struct {
	logger.Logger
}

func (m *mockLogger) Errorf(template string, args ...interface{}) {
	// Do nothing
}

// TestEncode_Success tests the successful encoding
func TestEncode_Success(t *testing.T) {
	// Mock the logger
	patchLogger := gomonkey.ApplyFunc(server.GetLogger, func() logger.Logger {
		return &mockLogger{}
	})
	defer patchLogger.Reset()

	text := "hello world"
	key := "ABCDABCDABCD1234"
	result, err := Encode([]byte(text), key)

	if err != nil {
		t.Errorf("Encode() error = %v, want nil", err)
	}

	if result == nil {
		t.Error("Encode() result = nil, want non-nil")
	}

	if string(result) == text {
		t.Error("Encode() result should not equal original text")
	}

	// Test that we can decode the result
	decryptText, err := Decode(result, key)
	if err != nil {
		t.Errorf("Decode() error = %v, want nil", err)
	}

	if string(decryptText) != text {
		t.Errorf("Decode() = %s, want %s", string(decryptText), text)
	}
}

// TestDecode_Success tests the successful decoding
func TestDecode_Success(t *testing.T) {
	// Mock the logger
	patchLogger := gomonkey.ApplyFunc(server.GetLogger, func() logger.Logger {
		return &mockLogger{}
	})
	defer patchLogger.Reset()

	// First encode some data using the real aes.Encrypt function
	text := "hello world"
	key := "ABCDABCDABCD1234"
	encoded, err := aes.Encrypt([]byte(text), []byte(key))
	if err != nil {
		t.Fatalf("Failed to encode test data: %v", err)
	}

	// Then decode it
	result, err := Decode(encoded, key)

	if err != nil {
		t.Errorf("Decode() error = %v, want nil", err)
	}

	if result == nil {
		t.Error("Decode() result = nil, want non-nil")
	}

	if string(result) != text {
		t.Errorf("Decode() = %s, want %s", string(result), text)
	}
}

// TestEncode_AesError tests the error handling in Encode
func TestEncode_AesError(t *testing.T) {
	// Mock the logger
	patchLogger := gomonkey.ApplyFunc(server.GetLogger, func() logger.Logger {
		return &mockLogger{}
	})
	defer patchLogger.Reset()

	// Mock aes.Encrypt to return an error
	expectedErr := errors.New("test error")
	patch := gomonkey.ApplyFunc(aes.Encrypt, func(data []byte, key []byte) ([]byte, error) {
		return nil, expectedErr
	})
	defer patch.Reset()

	text := "hello world"
	key := "ABCDABCDABCD1234"
	result, err := Encode([]byte(text), key)

	if err != expectedErr {
		t.Errorf("Encode() error = %v, want %v", err, expectedErr)
	}

	if result != nil {
		t.Error("Encode() result = non-nil, want nil")
	}
}

// TestDecode_AesError tests the error handling in Decode
func TestDecode_AesError(t *testing.T) {
	// Mock the logger
	patchLogger := gomonkey.ApplyFunc(server.GetLogger, func() logger.Logger {
		return &mockLogger{}
	})
	defer patchLogger.Reset()

	// Mock aes.Decrypt to return an error
	expectedErr := errors.New("test error")
	patch := gomonkey.ApplyFunc(aes.Decrypt, func(data []byte, key []byte) ([]byte, error) {
		return nil, expectedErr
	})
	defer patch.Reset()

	encoded := []byte("test encoded data")
	key := "ABCDABCDABCD1234"
	result, err := Decode(encoded, key)

	if err != expectedErr {
		t.Errorf("Decode() error = %v, want %v", err, expectedErr)
	}

	if result != nil {
		t.Error("Decode() result = non-nil, want nil")
	}
}