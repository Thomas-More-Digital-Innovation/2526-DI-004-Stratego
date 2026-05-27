package utils

import (
	"os"
	"testing"
)

const fallbackStr = "fallback"
const testValueStr = "test_value"
const testKeyStr = "TEST_KEY"

func TestGetEnv(t *testing.T) {
	_ = os.Setenv(testKeyStr, testValueStr)
	defer func() { _ = os.Unsetenv(testKeyStr) }()

	tests := []struct {
		key      string
		fallback string
		want     string
	}{
		{testKeyStr, fallbackStr, testValueStr},
		{"NON_EXISTENT", fallbackStr, fallbackStr},
	}

	for _, tt := range tests {
		if got := GetEnv(tt.key, tt.fallback); got != tt.want {
			t.Errorf("GetEnv(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.want)
		}
	}
}

func TestGetEnvOrError(t *testing.T) {
	_ = os.Setenv(testKeyStr, testValueStr)
	defer func() { _ = os.Unsetenv(testKeyStr) }()

	tests := []struct {
		key     string
		want    string
		wantErr bool
	}{
		{testKeyStr, testValueStr, false},
		{"NON_EXISTENT", "", true},
	}

	for _, tt := range tests {
		got, err := GetEnvOrError(tt.key)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetEnvOrError(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("GetEnvOrError(%q) = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestGetEnvOrErrorInProduction(t *testing.T) {
	// Since isProd is set in init(), it's hard to toggle it per test without global state mutation.
	// But we can test the current state.

	_ = os.Setenv("TEST_KEY_PROD", "prod_value")
	defer func() { _ = os.Unsetenv("TEST_KEY_PROD") }()

	// Test dev behavior (assuming isProd is false in test environment)
	if !IsProduction() {
		got, err := GetEnvOrErrorInProduction("NON_EXISTENT", "dev_fallback")
		if err != nil {
			t.Errorf("GetEnvOrErrorInProduction(dev) error = %v", err)
		}
		if got != "dev_fallback" {
			t.Errorf("GetEnvOrErrorInProduction(dev) = %q, want %q", got, "dev_fallback")
		}
	} else {
		// Test prod behavior
		got, err := GetEnvOrErrorInProduction("TEST_KEY_PROD", "fallback")
		if err != nil {
			t.Errorf("GetEnvOrErrorInProduction(prod) error = %v", err)
		}
		if got != "prod_value" {
			t.Errorf("GetEnvOrErrorInProduction(prod) = %q, want %q", got, "prod_value")
		}

		_, err = GetEnvOrErrorInProduction("NON_EXISTENT", "fallback")
		if err == nil {
			t.Error("GetEnvOrErrorInProduction(prod) expected error for non-existent key")
		}
	}
}
