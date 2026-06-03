package utils_test

import (
	"digital-innovation/gostrategy/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type validatorTest struct {
	input       string
	expectError bool
}

func TestValidateUsername(t *testing.T) {
	tests := []validatorTest{
		{"user123", false},
		{"user_name", false},
		{"user!name", true},
		{"<script>", true},
		{"   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := utils.ValidateUsername(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []validatorTest{
		{"Pass123!", false},
		{"Pass 123", false},
		{" Pass123!", true},  // leading space
		{"Pass123! ", true},  // trailing space
		{"  Pass123!", true}, // multiple leading spaces
		{"Pass123!  ", true}, // multiple trailing spaces
		{"   ", true},        // only spaces
		{"<script>alert(1)</script>", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := utils.ValidatePassword(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGeneric(t *testing.T) {
	tests := []validatorTest{
		{"user123", false},
		{"user_name", false},
		{"user!name", false},
		{"<script>", true},
		{"   ", false},
		{"Name with ;", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := utils.ValidateGeneric(tt.input, "Description")
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
