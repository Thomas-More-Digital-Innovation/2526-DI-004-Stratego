package utils_test

import (
	"digital-innovation/gostrategy/utils"
	"testing"
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
		err := utils.ValidateUsername(tt.input)
		if tt.expectError {
			if err == nil {
				t.Errorf("ValidateUsername(%q) expected an error but got none", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateUsername(%q) returned an unexpected error: %v", tt.input, err)
			}
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []validatorTest{
		{"Pass123!", false},
		{"Pass 123", false},
		{"<script>alert(1)</script>", true},
	}

	for _, tt := range tests {
		err := utils.ValidatePassword(tt.input)
		if tt.expectError {
			if err == nil {
				t.Errorf("ValidatePassword(%q) expected an error but got none", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ValidatePassword(%q) returned an unexpected error: %v", tt.input, err)
			}
		}
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
		err := utils.ValidateGeneric(tt.input, "Description")
		if tt.expectError {
			if err == nil {
				t.Errorf("ValidateGeneric(%q) expected an error but got none", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateGeneric(%q) returned an unexpected error: %v", tt.input, err)
			}
		}
	}

}
