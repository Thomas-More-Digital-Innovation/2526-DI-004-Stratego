package api

import (
	"testing"
)

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		password string
		want     bool
	}{
		{"Short1!", false},       // Too short (< 8)
		{"password", false},      // No number, no upper
		{"PASSWORD", false},      // No number, no lower
		{"12345678", false},      // No upper, no lower
		{"Pass123", false},       // Too short
		{"StrongPass1", true},    // Valid
		{"Another1Valid", true},  // Valid
		{"validButNoNum", false}, // No number
		{"VALID1NONO", false},    // No lower
		{"lower1nono", false},    // No upper
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			if got := isStrongPassword(tt.password); got != tt.want {
				t.Errorf("isStrongPassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}
