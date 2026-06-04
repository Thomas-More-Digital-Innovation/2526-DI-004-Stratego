package ai

import (
	ai_const "digital-innovation/gostrategy/internal/ai/const"
	"digital-innovation/gostrategy/internal/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefault(t *testing.T) {
	tests := []struct {
		aiType         string
		name           string
		wantAggression float64
		hasWeights     bool
		hasConfig      bool
	}{
		{"fato", "def", 0.5, false, false},
		{"heuristic", "def", 0.5, true, false},
		{"minimax", "def", 0.5, true, true},
		{"mcts", "def", 0.5, false, true},
		{"unknown", "def", 0.5, false, false},
	}

	for _, tc := range tests {
		t.Run(tc.aiType, func(t *testing.T) {
			params := GetDefault(tc.aiType, tc.name)
			assert.Equal(t, tc.aiType, params.AIType)
			assert.Equal(t, tc.name, params.Name)
			assert.Equal(t, tc.wantAggression, params.Aggression)
			if tc.hasWeights {
				assert.NotEmpty(t, params.Weights)
			} else {
				assert.Empty(t, params.Weights)
			}
			if tc.hasConfig {
				assert.NotEmpty(t, params.Config)
			} else {
				assert.Empty(t, params.Config)
			}
		})
	}
}

func TestParameterStore_FallbackFile(t *testing.T) {
	tempFile := "temp_test_ai_parameters.json"
	SetFallbackFile(tempFile)
	defer func() { _ = os.Remove(tempFile) }()

	params, err := Load("heuristic", "non-existent")
	assert.NoError(t, err)
	assert.Equal(t, "heuristic", params.AIType)
	assert.Equal(t, 0.5, params.Aggression)

	params.Aggression = 0.8
	params.Name = "custom_test"
	params.Weights = map[string]float64{ai_const.Flag: 20000.0}
	params.Config = map[string]any{"depth": 3.0}
	err = Save(params)
	assert.NoError(t, err)

	loaded, err := Load("heuristic", "custom_test")
	assert.NoError(t, err)
	assert.Equal(t, 0.8, loaded.Aggression)
	assert.Equal(t, 20000.0, loaded.Weights[ai_const.Flag])
	assert.Equal(t, 3.0, loaded.Config["depth"])

	params.Aggression = 0.95
	err = Save(params)
	assert.NoError(t, err)
	loaded, err = Load("heuristic", "custom_test")
	assert.NoError(t, err)
	assert.Equal(t, 0.95, loaded.Aggression)

	err = Save(nil)
	assert.Error(t, err)

	err = os.WriteFile(tempFile, []byte("invalid json"), 0600)
	assert.NoError(t, err)
	fallbackParams, err := Load("heuristic", "custom_test")
	assert.NoError(t, err)
	assert.Equal(t, 0.5, fallbackParams.Aggression)
}

func TestParameterStore_Database(t *testing.T) {
	db.SetupDBTest(t)

	tempFile := "temp_db_test_ai_parameters.json"
	SetFallbackFile(tempFile)
	defer func() { _ = os.Remove(tempFile) }()

	params := &Parameters{
		AIType:     "minimax",
		Name:       "db_test",
		Aggression: 0.77,
		Weights:    map[string]float64{ai_const.Marshal: 120.0},
		Config:     map[string]any{"depth": 4.0},
	}

	err := Save(params)
	assert.NoError(t, err)

	loaded, err := Load("minimax", "db_test")
	assert.NoError(t, err)
	assert.Equal(t, 0.77, loaded.Aggression)
	assert.Equal(t, 120.0, loaded.Weights[ai_const.Marshal])
	assert.Equal(t, 4.0, loaded.Config["depth"])

	params.Aggression = 0.99
	err = Save(params)
	assert.NoError(t, err)

	loaded2, err := Load("minimax", "db_test")
	assert.NoError(t, err)
	assert.Equal(t, 0.99, loaded2.Aggression)
}
