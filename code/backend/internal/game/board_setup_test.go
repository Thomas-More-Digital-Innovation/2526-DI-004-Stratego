package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeBoardSetup(t *testing.T) {
	setup := []string{
		"9876543210",
		"BBBBBBBBBB",
		"1111111111",
		"MMMMMMMMMM",
	}

	encoded, err := EncodeBoardSetup(setup)
	require.NoError(t, err)

	decoded, err := DecodeBoardSetup(encoded)
	require.NoError(t, err)
	assert.Equal(t, setup, decoded)
}

func TestEncodeBoardSetupSingleString(t *testing.T) {
	setup := "9876543210BBBBBBBBBB1111111111MMMMMMMMMM"
	encoded, err := EncodeBoardSetup([]string{setup})
	require.NoError(t, err)

	decoded, err := DecodeBoardSetup(encoded)
	require.NoError(t, err)

	conc := ""
	for _, r := range decoded {
		conc += r
	}
	assert.Equal(t, setup, conc)
}

func TestEncodeBoardSetupErrors(t *testing.T) {
	tests := []struct {
		name    string
		rows    []string
		wantErr bool
	}{
		{"no rows", []string{}, true},
		{"wrong row count", []string{"1", "2"}, true},
		{"wrong row length", []string{"1234567890", "1234567890", "1234567890", "123456789"}, true},
		{"unknown char", []string{"X234567890", "1234567890", "1234567890", "1234567890"}, true},
		{"unknown char in lo", []string{"1X34567890", "1234567890", "1234567890", "1234567890"}, true},
		{"wrong conc length", []string{"1234567890123456789012345678901234567890EXTRA"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncodeBoardSetup(tt.rows)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDecodeBoardSetupErrors(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
		wantErr bool
	}{
		{"invalid base64", "!!!", true},
		{"wrong length", "abcd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeBoardSetup(tt.encoded)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseBoardSetupSmart(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    int // number of rows
		wantErr bool
	}{
		{"concatenated", "9876543210BBBBBBBBBB1111111111MMMMMMMMMM", 4, false},
		{"newlines", "9876543210\nBBBBBBBBBB\n1111111111\nMMMMMMMMMM", 4, false},
		{"commas", "9876543210,BBBBBBBBBB,1111111111,MMMMMMMMMM", 4, false},
		{"invalid format", "too short", 0, true},
		{"wrong row length", "9876543210\nBBB\n1111111111\nMMMMMMMMMM", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBoardSetupSmart(tt.raw)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, got, tt.want)
			}
		})
	}
}
