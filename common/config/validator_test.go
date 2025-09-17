package config

import (
	"testing"
	"time"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "testField",
		Value:   "testValue",
		Message: "test message",
	}

	expected := "validation error in field 'testField' with value 'testValue': test message"
	result := err.Error()

	if result != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", result, expected)
	}
}

func TestValidationErrors_Error(t *testing.T) {
	errs := &ValidationErrors{
		Errors: []error{
			&ValidationError{Field: "field1", Value: "value1", Message: "error1"},
			&ValidationError{Field: "field2", Value: "value2", Message: "error2"},
		},
	}

	result := errs.Error()
	expected := "2 validation errors occurred:\n  - validation error in field 'field1' with value 'value1': error1\n  - validation error in field 'field2' with value 'value2': error2"

	if result != expected {
		t.Errorf("ValidationErrors.Error() = %v, want %v", result, expected)
	}
}

func TestValidationErrors_Error_Empty(t *testing.T) {
	errs := &ValidationErrors{
		Errors: []error{},
	}

	result := errs.Error()
	expected := "no validation errors"

	if result != expected {
		t.Errorf("ValidationErrors.Error() = %v, want %v", result, expected)
	}
}

func TestValidatePort_Valid(t *testing.T) {
	tests := []struct {
		port int
	}{
		{1},
		{80},
		{443},
		{8080},
		{65535},
	}

	for _, tt := range tests {
		err := ValidatePort(tt.port)
		if err != nil {
			t.Errorf("ValidatePort(%d) = %v, want nil", tt.port, err)
		}
	}
}

func TestValidatePort_Invalid(t *testing.T) {
	tests := []struct {
		port    int
		wantErr bool
	}{
		{0, true},
		{-1, true},
		{65536, true},
		{100000, true},
	}

	for _, tt := range tests {
		err := ValidatePort(tt.port)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePort(%d) error = %v, wantErr %v", tt.port, err, tt.wantErr)
		}
	}
}

func TestValidateHost_Valid(t *testing.T) {
	tests := []struct {
		host string
	}{
		{"localhost"},
		{"example.com"},
		{"192.168.1.1"},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		{"a"},
		{"a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z"},
	}

	for _, tt := range tests {
		err := ValidateHost(tt.host)
		if err != nil {
			t.Errorf("ValidateHost(%s) = %v, want nil", tt.host, err)
		}
	}
}

func TestValidateHost_Invalid(t *testing.T) {
	err := ValidateHost("")
	if err == nil {
		t.Error("ValidateHost(\"\") = nil, want error")
	}
}

func TestValidateURL_Valid(t *testing.T) {
	tests := []struct {
		url string
	}{
		{"http://example.com"},
		{"https://example.com"},
		{"http://localhost:8080"},
		{"https://api.example.com/v1/resource"},
	}

	for _, tt := range tests {
		err := ValidateURL(tt.url)
		if err != nil {
			t.Errorf("ValidateURL(%s) = %v, want nil", tt.url, err)
		}
	}
}

func TestValidateURL_Invalid(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"", true},
		{"ftp://example.com", true},
		{"example.com", true},
		{"//example.com", true},
	}

	for _, tt := range tests {
		err := ValidateURL(tt.url)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateURL(%s) error = %v, wantErr %v", tt.url, err, tt.wantErr)
		}
	}
}

func TestValidateDuration_Valid(t *testing.T) {
	tests := []struct {
		duration time.Duration
	}{
		{1 * time.Nanosecond},
		{1 * time.Microsecond},
		{1 * time.Millisecond},
		{1 * time.Second},
		{1 * time.Minute},
		{1 * time.Hour},
	}

	for _, tt := range tests {
		err := ValidateDuration(tt.duration)
		if err != nil {
			t.Errorf("ValidateDuration(%v) = %v, want nil", tt.duration, err)
		}
	}
}

func TestValidateDuration_Invalid(t *testing.T) {
	tests := []struct {
		duration time.Duration
		wantErr  bool
	}{
		{0, true},
		{-1 * time.Nanosecond, true},
		{-1 * time.Second, true},
	}

	for _, tt := range tests {
		err := ValidateDuration(tt.duration)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateDuration(%v) error = %v, wantErr %v", tt.duration, err, tt.wantErr)
		}
	}
}

func TestValidateNonEmptyString_Valid(t *testing.T) {
	err := ValidateNonEmptyString("test", "fieldName")
	if err != nil {
		t.Errorf("ValidateNonEmptyString(\"test\", \"fieldName\") = %v, want nil", err)
	}
}

func TestValidateNonEmptyString_Invalid(t *testing.T) {
	err := ValidateNonEmptyString("", "fieldName")
	if err == nil {
		t.Error("ValidateNonEmptyString(\"\", \"fieldName\") = nil, want error")
	}
}

func TestValidatePositiveInt_Valid(t *testing.T) {
	err := ValidatePositiveInt(1, "fieldName")
	if err != nil {
		t.Errorf("ValidatePositiveInt(1, \"fieldName\") = %v, want nil", err)
	}
}

func TestValidatePositiveInt_Invalid(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{0, true},
		{-1, true},
		{-100, true},
	}

	for _, tt := range tests {
		err := ValidatePositiveInt(tt.value, "fieldName")
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePositiveInt(%d, \"fieldName\") error = %v, wantErr %v", tt.value, err, tt.wantErr)
		}
	}
}

func TestValidateRangeInt_Valid(t *testing.T) {
	err := ValidateRangeInt(5, 1, 10, "fieldName")
	if err != nil {
		t.Errorf("ValidateRangeInt(5, 1, 10, \"fieldName\") = %v, want nil", err)
	}
}

func TestValidateRangeInt_Invalid(t *testing.T) {
	tests := []struct {
		value   int
		min     int
		max     int
		wantErr bool
	}{
		{0, 1, 10, true},
		{11, 1, 10, true},
		{-1, 1, 10, true},
		{15, 1, 10, true},
	}

	for _, tt := range tests {
		err := ValidateRangeInt(tt.value, tt.min, tt.max, "fieldName")
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateRangeInt(%d, %d, %d, \"fieldName\") error = %v, wantErr %v", tt.value, tt.min, tt.max, err, tt.wantErr)
		}
	}
}

func TestValidateBool(t *testing.T) {
	// ValidateBool always returns nil, so we just test that it doesn't panic
	err := ValidateBool(true, "fieldName")
	if err != nil {
		t.Errorf("ValidateBool(true, \"fieldName\") = %v, want nil", err)
	}

	err = ValidateBool(false, "fieldName")
	if err != nil {
		t.Errorf("ValidateBool(false, \"fieldName\") = %v, want nil", err)
	}
}