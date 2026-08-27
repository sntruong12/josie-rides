package validator

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Define a new Validator type which contains a map of validation errors for our
// form fields.
type Validator struct {
	FieldErrors map[string]string
}

// Valid() returns true if the FieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't already
	// initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt() returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func ValidInt(value string, min int, max int) bool {
	ival, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return ival >= min && ival <= max
}

func ValidFloat(value string, min float64, max float64) bool {
	ival, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	return ival >= min && ival <= max
}

func ValidTimeAndTimezone(value string, timezone string) bool {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return false
	}
	_, err = time.ParseInLocation("2006-01-02T15:04", value, loc)
	return err == nil
}
