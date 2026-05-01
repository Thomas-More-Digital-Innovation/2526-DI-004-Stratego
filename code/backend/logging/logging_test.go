package logging

import (
	"testing"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"normal string", "normal string"},
		{"  trimmed  ", "trimmed"},
		{"string with\nnewline", "string with newline"},
		{"string with\x00null", "string with null"},
		{"\x1b[31mcolor codes\x1b[0m", "[31mcolor codes [0m"}, // control chars replaced by spaces, then trimmed
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := sanitize(tt.input); got != tt.want {
				t.Errorf("sanitize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatUser(t *testing.T) {
	tests := []struct {
		username string
		userID   int
		want     string
	}{
		{"Alice", 1, "Alice (ID: 1)"},
		{"Bob\nSmith", 2, "Bob Smith (ID: 2)"},
		{"", 0, "Guest"},
		{"", -1, "Guest"},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			if got := FormatUser(tt.username, tt.userID); got != tt.want {
				t.Errorf("FormatUser(%q, %d) = %q, want %q", tt.username, tt.userID, got, tt.want)
			}
		})
	}
}

func TestGetCallerLocation(t *testing.T) {
	// depth 1 should be this function
	loc := getCallerLocation(1)
	if loc == "" || loc == "???:0" {
		t.Errorf("Expected valid location, got %q", loc)
	}
	// It should contain logging_test.go
	if !contains(loc, "logging_test.go") {
		t.Errorf("Expected location to contain logging_test.go, got %q", loc)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
