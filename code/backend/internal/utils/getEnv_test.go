package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		t.Run(tt.key, func(t *testing.T) {
			assert.Equal(t, tt.want, GetEnv(tt.key, tt.fallback))
		})
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
		t.Run(tt.key, func(t *testing.T) {
			got, err := GetEnvOrError(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetEnvOrErrorInProduction(t *testing.T) {
	_ = os.Setenv("TEST_KEY_PROD", "prod_value")
	defer func() { _ = os.Unsetenv("TEST_KEY_PROD") }()

	if !IsProduction() {
		got, err := GetEnvOrErrorInProduction("NON_EXISTENT", "dev_fallback")
		assert.NoError(t, err)
		assert.Equal(t, "dev_fallback", got)
	} else {
		got, err := GetEnvOrErrorInProduction("TEST_KEY_PROD", "fallback")
		assert.NoError(t, err)
		assert.Equal(t, "prod_value", got)

		_, err = GetEnvOrErrorInProduction("NON_EXISTENT", "fallback")
		assert.Error(t, err)
	}
}
