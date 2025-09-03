package error

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	err := New(50001, CategoryInternal, "test error")
	expected := "[internal:50001] test error"
	if err.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, err.Error())
	}
}

func TestError_WithContext(t *testing.T) {
	err := New(50001, CategoryInternal, "test error").
		WithContext("key1", "value1").
		WithContext("key2", 123)

	if err.Context["key1"] != "value1" {
		t.Errorf("Expected context key1 to be 'value1', got %v", err.Context["key1"])
	}

	if err.Context["key2"] != 123 {
		t.Errorf("Expected context key2 to be 123, got %v", err.Context["key2"])
	}
}

func TestError_Wrap(t *testing.T) {
	cause := errors.New("original error")
	err := Wrap(cause, 50001, CategoryInternal, "wrapped error")

	if !errors.Is(err, cause) {
		t.Error("Expected err to wrap cause")
	}

	if err.Cause != cause {
		t.Error("Expected err.Cause to be cause")
	}
}

func TestError_Is(t *testing.T) {
	err1 := New(50001, CategoryInternal, "test error")
	err2 := New(50001, CategoryInternal, "test error")
	err3 := New(50002, CategoryInternal, "different error")

	if !err1.Is(err2) {
		t.Error("Expected err1 to be equal to err2")
	}

	if err1.Is(err3) {
		t.Error("Expected err1 to not be equal to err3")
	}
}

func TestError_MarshalJSON(t *testing.T) {
	err := New(50001, CategoryInternal, "test error").
		WithContext("key", "value")

	jsonData, jsonErr := err.MarshalJSON()
	if jsonErr != nil {
		t.Errorf("Expected nil error, got %v", jsonErr)
	}

	expectedSubstrings := []string{
		`"code":50001`,
		`"message":"test error"`,
		`"category":"internal"`,
		`"key":"value"`,
		`"error":"[internal:50001] test error"`,
	}

	for _, substring := range expectedSubstrings {
		if !containsString(string(jsonData), substring) {
			t.Errorf("Expected JSON to contain %s, got %s", substring, string(jsonData))
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
			(len(s) > len(substr) && 
				(s[:len(substr)] == substr || 
					s[len(s)-len(substr):] == substr ||
					findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}