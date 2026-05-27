package logging

import (
	"fmt"
	"testing"
)

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
		{"", 0, GuestUser},
		{"", -1, GuestUser},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			if got := FormatUser(tt.username, tt.userID); got != tt.want {
				t.Errorf("FormatUser(%q, %d) = %q, want %q", tt.username, tt.userID, got, tt.want)
			}
		})
	}
}

func TestLoggingFunctions(_ *testing.T) {
	// These tests just call the functions to ensure they don't panic unexpectedly
	// and to get coverage. We can't easily verify the output without mocking log.

	err := fmt.Errorf("test error")

	UserRegistered("testuser", 1)
	UserLoggedIn("testuser", 1)

	GameStarted("game1", "pvp", "testuser", 1)
	GameFinished("game1", "winner", "loser", 10)
	GameAborted("game1", "timeout", "testuser", 1)

	ConnectionError("game1", "testuser", 1, err)
	ConnectionClosed("game1", "testuser", 1)

	Error("message", err)
	ErrorWithUser("message", "testuser", 1, err)
	ErrorWith2Users("message", "user1", 1, "user2", 2, err)
	ErrorWithIP("message", "127.0.0.1", err)

	SecurityWarning("message", "details", "testuser", 1)
	SecurityWarningWithIP("message", "details", "127.0.0.1")

	LogRaw("raw message")
}

func TestFatal(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Fatal did not panic")
		}
	}()
	Fatal("fatal message", fmt.Errorf("fatal error"))
}

func TestFatalf(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Fatalf did not panic")
		}
	}()
	Fatalf("fatal message %s", "formatted")
}

func TestDebug(t *testing.T) {
	t.Setenv("GOSTRATEGY_DEBUG", "true")
	Debug("TAG", "debug message %s", "val")
	DebugWithUser("TAG", "user", 1, "debug message %s", "val")

	t.Setenv("GOSTRATEGY_DEBUG", "false")
	Debug("TAG", "debug message %s", "val")
	DebugWithUser("TAG", "user", 1, "debug message %s", "val")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
