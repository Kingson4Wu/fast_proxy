package config

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidatePort(t *testing.T) {
	Convey("Given port validation function", t, func() {
		Convey("When validating valid ports", func() {
			validPorts := []int{1, 80, 443, 8080, 65535}
			
			for _, port := range validPorts {
				err := ValidatePort(port)
				
				Convey(fmt.Sprintf("Port %d should be valid", port), func() {
					So(err, ShouldBeNil)
				})
			}
		})

		Convey("When validating invalid ports", func() {
			invalidPorts := []int{-1, 0, 65536, 100000}
			
			for _, port := range invalidPorts {
				err := ValidatePort(port)
				
				Convey(fmt.Sprintf("Port %d should be invalid", port), func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "port must be between 1 and 65535")
				})
			}
		})
	})
}

func TestValidateHost(t *testing.T) {
	Convey("Given host validation function", t, func() {
		Convey("When validating valid hosts", func() {
			validHosts := []string{
				"localhost",
				"127.0.0.1",
				"::1",
				"example.com",
				"subdomain.example.com",
				"192.168.1.1",
			}
			
			for _, host := range validHosts {
				err := ValidateHost(host)
				
				Convey(fmt.Sprintf("Host '%s' should be valid", host), func() {
					So(err, ShouldBeNil)
				})
			}
		})

		Convey("When validating invalid hosts", func() {
			Convey("Empty host should be invalid", func() {
				err := ValidateHost("")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "host cannot be empty")
			})
		})
	})
}

func TestValidateURL(t *testing.T) {
	Convey("Given URL validation function", t, func() {
		Convey("When validating valid URLs", func() {
			validURLs := []string{
				"http://example.com",
				"https://example.com",
				"http://localhost:8080",
				"https://192.168.1.1:443/path",
			}
			
			for _, url := range validURLs {
				err := ValidateURL(url)
				
				Convey(fmt.Sprintf("URL '%s' should be valid", url), func() {
					So(err, ShouldBeNil)
				})
			}
		})

		Convey("When validating invalid URLs", func() {
			Convey("Empty URL should be invalid", func() {
				err := ValidateURL("")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "URL cannot be empty")
			})

			Convey("URL without scheme should be invalid", func() {
				err := ValidateURL("example.com")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "URL must start with http:// or https://")
			})
		})
	})
}

func TestValidateDuration(t *testing.T) {
	Convey("Given duration validation function", t, func() {
		Convey("When validating positive durations", func() {
			validDurations := []time.Duration{
				1 * time.Second,
				5 * time.Minute,
				1 * time.Hour,
				24 * time.Hour,
			}
			
			for _, duration := range validDurations {
				err := ValidateDuration(duration)
				
				Convey(fmt.Sprintf("Duration %v should be valid", duration), func() {
					So(err, ShouldBeNil)
				})
			}
		})

		Convey("When validating non-positive durations", func() {
			invalidDurations := []time.Duration{
				0,
				-1 * time.Second,
				-5 * time.Minute,
			}
			
			for _, duration := range invalidDurations {
				err := ValidateDuration(duration)
				
				Convey(fmt.Sprintf("Duration %v should be invalid", duration), func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "duration must be positive")
				})
			}
		})
	})
}

func TestValidateNonEmptyString(t *testing.T) {
	Convey("Given non-empty string validation function", t, func() {
		Convey("When validating non-empty strings", func() {
			err := ValidateNonEmptyString("test", "test_field")
			So(err, ShouldBeNil)
		})

		Convey("When validating empty strings", func() {
			err := ValidateNonEmptyString("", "test_field")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "test_field cannot be empty")
		})
	})
}

func TestValidatePositiveInt(t *testing.T) {
	Convey("Given positive integer validation function", t, func() {
		Convey("When validating positive integers", func() {
			err := ValidatePositiveInt(5, "test_field")
			So(err, ShouldBeNil)
		})

		Convey("When validating non-positive integers", func() {
			err := ValidatePositiveInt(0, "test_field")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "test_field must be positive")

			err = ValidatePositiveInt(-1, "test_field")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "test_field must be positive")
		})
	})
}

func TestValidateRangeInt(t *testing.T) {
	Convey("Given range integer validation function", t, func() {
		Convey("When validating integers within range", func() {
			err := ValidateRangeInt(5, 1, 10, "test_field")
			So(err, ShouldBeNil)
		})

		Convey("When validating integers outside range", func() {
			err := ValidateRangeInt(0, 1, 10, "test_field")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "test_field must be between 1 and 10")

			err = ValidateRangeInt(11, 1, 10, "test_field")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "test_field must be between 1 and 10")
		})
	})
}

func TestValidationError(t *testing.T) {
	Convey("Given a ValidationError", t, func() {
		err := &ValidationError{
			Field:   "test_field",
			Value:   "test_value",
			Message: "test message",
		}

		Convey("Should format error message correctly", func() {
			expected := "validation error in field 'test_field' with value 'test_value': test message"
			So(err.Error(), ShouldEqual, expected)
		})
	})
}

func TestValidationErrors(t *testing.T) {
	Convey("Given ValidationErrors", t, func() {
		Convey("With no errors", func() {
			errs := &ValidationErrors{Errors: []error{}}
			So(errs.Error(), ShouldEqual, "no validation errors")
		})

		Convey("With multiple errors", func() {
			err1 := &ValidationError{Field: "field1", Value: "value1", Message: "message1"}
			err2 := &ValidationError{Field: "field2", Value: "value2", Message: "message2"}
			errs := &ValidationErrors{Errors: []error{err1, err2}}

			errStr := errs.Error()
			So(errStr, ShouldContainSubstring, "2 validation errors occurred")
			So(errStr, ShouldContainSubstring, "field1")
			So(errStr, ShouldContainSubstring, "field2")
		})
	})
}