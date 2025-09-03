package config

import (
	"fmt"
	"net"
	"time"
)

// Validator interface for configuration validation
type Validator interface {
	Validate() error
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s' with value '%v': %s", e.Field, e.Value, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []error
}

func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "no validation errors"
	}
	
	result := fmt.Sprintf("%d validation errors occurred:", len(e.Errors))
	for _, err := range e.Errors {
		result += fmt.Sprintf("\n  - %s", err.Error())
	}
	return result
}

// ValidatePort validates that a port number is in the valid range
func ValidatePort(port int) error {
	if port <= 0 || port > 65535 {
		return &ValidationError{
			Field:   "port",
			Value:   port,
			Message: "port must be between 1 and 65535",
		}
	}
	return nil
}

// ValidateHost validates that a host is a valid IP address or hostname
func ValidateHost(host string) error {
	if host == "" {
		return &ValidationError{
			Field:   "host",
			Value:   host,
			Message: "host cannot be empty",
		}
	}
	
	// Check if it's a valid IP address
	if ip := net.ParseIP(host); ip != nil {
		return nil
	}
	
	// Check if it's a valid hostname
	// This is a simple check, in production you might want to be more thorough
	if len(host) > 253 {
		return &ValidationError{
			Field:   "host",
			Value:   host,
			Message: "hostname too long",
		}
	}
	
	return nil
}

// ValidateURL validates that a URL is properly formatted
func ValidateURL(url string) error {
	if url == "" {
		return &ValidationError{
			Field:   "url",
			Value:   url,
			Message: "URL cannot be empty",
		}
	}
	
	// Simple URL validation - in production you might want to use net/url package
	if len(url) < 4 || (url[:4] != "http" && url[:5] != "https") {
		return &ValidationError{
			Field:   "url",
			Value:   url,
			Message: "URL must start with http:// or https://",
		}
	}
	
	return nil
}

// ValidateDuration validates that a duration is positive
func ValidateDuration(duration time.Duration) error {
	if duration <= 0 {
		return &ValidationError{
			Field:   "duration",
			Value:   duration,
			Message: "duration must be positive",
		}
	}
	return nil
}

// ValidateNonEmptyString validates that a string is not empty
func ValidateNonEmptyString(value, fieldName string) error {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Message: fmt.Sprintf("%s cannot be empty", fieldName),
		}
	}
	return nil
}

// ValidatePositiveInt validates that an integer is positive
func ValidatePositiveInt(value int, fieldName string) error {
	if value <= 0 {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Message: fmt.Sprintf("%s must be positive", fieldName),
		}
	}
	return nil
}

// ValidateRangeInt validates that an integer is within a range
func ValidateRangeInt(value, min, max int, fieldName string) error {
	if value < min || value > max {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Message: fmt.Sprintf("%s must be between %d and %d", fieldName, min, max),
		}
	}
	return nil
}

// ValidateBool validates a boolean value (always returns nil, but included for consistency)
func ValidateBool(value bool, fieldName string) error {
	return nil
}