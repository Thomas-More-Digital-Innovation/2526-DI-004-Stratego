package ai

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameterStore(t *testing.T) {
	tempFile := "temp_test_ai_parameters.json"
	SetFallbackFile(tempFile)
	defer func() { _ = os.Remove(tempFile) }()

	params, err := Load("heuristic", "non-existent")
	assert.NoError(t, err)
	assert.Equal(t, "heuristic", params.AIType)
	assert.Equal(t, 0.5, params.Aggression)
	assert.NotEmpty(t, params.Weights)

	params.Aggression = 0.9
	params.Name = "custom_test"
	params.Weights["Flag"] = 20000.0
	err = Save(params)
	assert.NoError(t, err)

	loaded, err := Load("heuristic", "custom_test")
	assert.NoError(t, err)
	assert.Equal(t, 0.9, loaded.Aggression)
	assert.Equal(t, 20000.0, loaded.Weights["Flag"])
}
