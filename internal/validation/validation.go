package validation

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Amount parses s and requires a value greater than zero.
func Amount(s string) (float64, error) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("amount must be a number")
	}
	if v <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}
	return v, nil
}

// Date parses s as YYYY-MM-DD.
func Date(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, fmt.Errorf("please pick a valid date")
	}
	return t, nil
}

// Required checks s is non-empty after trimming; field names the input in the error.
func Required(s, field string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("%s is required", field)
	}
	return s, nil
}

// UUIDField checks s is a valid UUID; field names the input in the error.
func UUIDField(s, field string) (string, error) {
	s = strings.TrimSpace(s)
	if _, err := uuid.Parse(s); err != nil {
		return "", fmt.Errorf("please select a %s", field)
	}
	return s, nil
}

// IntField parses s as an integer; field names the input in the error.
func IntField(s, field string) (int, error) {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, fmt.Errorf("please select a %s", field)
	}
	return v, nil
}
