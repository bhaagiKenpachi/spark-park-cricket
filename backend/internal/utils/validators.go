package utils

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidateRequired validates that a field is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(fieldName + " is required")
	}
	return nil
}

// ValidateMinLength validates minimum length
func ValidateMinLength(value string, minLength int, fieldName string) error {
	if len(strings.TrimSpace(value)) < minLength {
		return errors.New(fieldName + " must be at least " + string(rune(minLength)) + " characters long")
	}
	return nil
}

// ValidateMaxLength validates maximum length
func ValidateMaxLength(value string, maxLength int, fieldName string) error {
	if len(value) > maxLength {
		return errors.New(fieldName + " must not exceed " + string(rune(maxLength)) + " characters")
	}
	return nil
}

// ValidateRange validates that a number is within a range
func ValidateRange(value, min, max int, fieldName string) error {
	if value < min || value > max {
		return errors.New(fieldName + " must be between " + string(rune(min)) + " and " + string(rune(max)))
	}
	return nil
}

// ValidateStruct validates a struct using reflection (basic implementation)
func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Check for required fields (basic validation)
		if fieldType.Tag.Get("validate") != "" {
			if field.Kind() == reflect.String && field.String() == "" {
				return errors.New(fieldType.Name + " is required")
			}
			if field.Kind() == reflect.Int && field.Int() == 0 {
				return errors.New(fieldType.Name + " is required")
			}
		}
	}

	return nil
}
