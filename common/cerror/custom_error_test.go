package cerror

import (
	"testing"
)

func TestErr_Error(t *testing.T) {
	err := &Err{
		Code: 404,
		Msg:  "Not Found",
	}

	expected := `{"Code":404,"Msg":"Not Found"}`
	result := err.Error()

	if result != expected {
		t.Errorf("Error() = %v, want %v", result, expected)
	}
}

func TestNewError(t *testing.T) {
	code := 500
	msg := "Internal Server Error"

	err := NewError(code, msg)

	if err.Code != code {
		t.Errorf("NewError().Code = %v, want %v", err.Code, code)
	}

	if err.Msg != msg {
		t.Errorf("NewError().Msg = %v, want %v", err.Msg, msg)
	}
}

func TestNewError_NilValues(t *testing.T) {
	err := NewError(0, "")

	if err.Code != 0 {
		t.Errorf("NewError().Code = %v, want %v", err.Code, 0)
	}

	if err.Msg != "" {
		t.Errorf("NewError().Msg = %v, want %v", err.Msg, "")
	}
}