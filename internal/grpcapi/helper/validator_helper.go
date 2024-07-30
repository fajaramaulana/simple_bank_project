package helper

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	IsValidUsername = "`^[a-z0-9_]+$`"
	IsValidFullName = "`^[a-zA-Z\\s]+$`"
)

// ValidateRequired checks if the field is not empty
func ValidateRequired(value string) error {
	if value == "" {
		return fmt.Errorf("is required")
	}
	return nil
}

// ValidateMin checks if the field meets the minimum length requirement
func ValidateMin(value string, minLength int) error {
	if len(value) < minLength {
		return fmt.Errorf("must be at least %d characters long", minLength)
	}
	return nil
}

// ValidateMax checks if the field does not exceed the maximum length
func ValidateMax(value string, maxLength int) error {
	if len(value) > maxLength {
		return fmt.Errorf("must be at most %d characters long", maxLength)
	}
	return nil
}

// ValidateAlphanum checks if the field contains only alphanumeric characters
func ValidateAlphanum(value string) error {
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", value); !matched {
		return fmt.Errorf("can only contain alphanumeric characters")
	}
	return nil
}

// validateAlphaSpace checks if the field contains only alphabetic characters with space
func ValidateAlphaSpace(value string) error {
	if matched, _ := regexp.MatchString("^[a-zA-Z\\s]+$", value); !matched {
		return fmt.Errorf("can only contain alphabetic characters with space")
	}
	return nil
}

// ValidateEmail checks if the field is a valid email address
func ValidateEmail(value string) error {
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, value); !matched {
		return fmt.Errorf("must be a valid email address")
	}
	return nil
}

// ValidateLen checks if the field length is exactly the specified length
func ValidateLen(value string, length int) error {
	if len(value) != length {
		return fmt.Errorf("must be exactly %d characters long", length)
	}
	return nil
}

// ValidateOneOf checks if the field is one of the specified values
func ValidateOneOf(value string, options []string) error {
	for _, option := range options {
		if value == option {
			return nil
		}
	}
	return fmt.Errorf("must be one of [%s]", strings.Join(options, ", "))
}

// ValidateNumeric checks if the field contains only numeric characters
func ValidateNumeric(value string) error {
	if matched, _ := regexp.MatchString("^[0-9]+$", value); !matched {
		return fmt.Errorf("must be a numeric value")
	}
	return nil
}

// ValidateURL checks if the field is a valid URL
func ValidateURL(value string) error {
	if matched, _ := regexp.MatchString(`^https?://[^\s/$.?#].[^\s]*$`, value); !matched {
		return fmt.Errorf("must be a valid URL")
	}
	return nil
}

// ValidateUUID checks if the field is a valid UUID
func ValidateUUID(value string) error {
	if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, value); !matched {
		return fmt.Errorf("be a valid UUID")
	}
	return nil
}

// ValidateEqField checks if the field is equal to another field
func ValidateEqField(value string, otherValue string, otherFieldName string) error {
	if value != otherValue {
		return fmt.Errorf("be equal to %s", otherFieldName)
	}
	return nil
}

// ValidateNeField checks if the field is not equal to another field
func ValidateNeField(value string, otherValue string, otherFieldName string) error {
	if value == otherValue {
		return fmt.Errorf("not be equal to %s", otherFieldName)
	}
	return nil
}

// ValidateGTE checks if the field is greater than or equal to a value
func ValidateGTE(value string, param string) error {
	if len(value) < len(param) {
		return fmt.Errorf("be greater than or equal to %s", param)
	}
	return nil
}

// ValidateLTE checks if the field is less than or equal to a value
func ValidateLTE(value string, param string) error {
	if len(value) > len(param) {
		return fmt.Errorf("be less than or equal to %s", param)
	}
	return nil
}

// ValidateGT checks if the field is greater than a value
func ValidateGT(value string, param string) error {
	if len(value) <= len(param) {
		return fmt.Errorf("be greater than %s", param)
	}
	return nil
}

// ValidateLT checks if the field is less than a value
func ValidateLT(value string, param string) error {
	if len(value) >= len(param) {
		return fmt.Errorf("be less than %s", param)
	}
	return nil
}

// ValidateContains checks if the field contains a value
func ValidateContains(value string, param string) error {
	if !strings.Contains(value, param) {
		return fmt.Errorf("contain %s", param)
	}
	return nil
}

// ValidateExcludes checks if the field does not contain a value
func ValidateExcludes(value string, param string) error {
	if strings.Contains(value, param) {
		return fmt.Errorf("not contain %s", param)
	}
	return nil
}

// ValidateIP checks if the field is a valid IP address
func ValidateIP(value string) error {
	if matched, _ := regexp.MatchString(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`, value); !matched {
		return fmt.Errorf("be a valid IP address")
	}
	return nil
}

// ValidateRegex checks if the field matches the specified regular expression
func ValidateRegex(value string, regex string) error {
	if matched, _ := regexp.MatchString(regex, value); !matched {
		return fmt.Errorf("match the pattern %s", regex)
	}
	return nil
}

// Validate slice check if the value is in the slice
func ValidateSlice(value string, slice []string) error {
	for _, v := range slice {
		if v == value {
			return nil
		}
	}
	return fmt.Errorf("must be one of %v", slice)
}
