package AIhandler

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAI(t *testing.T) {
	player := game.NewPlayer(0, "Piet", "red")

	types := []string{models.Fafo, models.Fato, models.Heuristic, models.Minimax, models.Mcts}
	for _, typ := range types {
		t.Run(typ, func(t *testing.T) {
			instance, err := CreateAI(typ, &player)
			assert.NoError(t, err)
			assert.NotNil(t, instance)
		})
	}

	_, err := CreateAI("invalid", &player)
	assert.Error(t, err)
}

func TestCreateAIWithOptions(t *testing.T) {
	player := game.NewPlayer(0, "Piet", "red")

	t.Run("Fato options", func(t *testing.T) {
		opts := map[string]any{"aggression": 0.8}
		instance, err := CreateAIWithOptions(models.Fato, &player, opts)
		assert.NoError(t, err)
		assert.NotNil(t, instance)

		// Without option
		instance2, err := CreateAIWithOptions(models.Fato, &player, nil)
		assert.NoError(t, err)
		assert.NotNil(t, instance2)
	})

	t.Run("Heuristic options", func(t *testing.T) {
		opts := map[string]any{
			"name":       "custom",
			"aggression": 0.35,
		}
		instance, err := CreateAIWithOptions(models.Heuristic, &player, opts)
		assert.NoError(t, err)
		assert.NotNil(t, instance)

		// Without options
		instance2, err := CreateAIWithOptions(models.Heuristic, &player, nil)
		assert.NoError(t, err)
		assert.NotNil(t, instance2)
	})

	t.Run("Minimax options", func(t *testing.T) {
		opts := map[string]any{
			"name":       "custom_minimax",
			"aggression": 0.9,
			"depth":      float64(4),
		}
		instance, err := CreateAIWithOptions(models.Minimax, &player, opts)
		assert.NoError(t, err)
		assert.NotNil(t, instance)

		// Without options
		instance2, err := CreateAIWithOptions(models.Minimax, &player, nil)
		assert.NoError(t, err)
		assert.NotNil(t, instance2)
	})

	t.Run("Mcts options", func(t *testing.T) {
		opts := map[string]any{
			"name":       "custom_mcts",
			"aggression": 0.2,
			"iterations": float64(10),
		}
		instance, err := CreateAIWithOptions(models.Mcts, &player, opts)
		assert.NoError(t, err)
		assert.NotNil(t, instance)

		// Without options
		instance2, err := CreateAIWithOptions(models.Mcts, &player, nil)
		assert.NoError(t, err)
		assert.NotNil(t, instance2)
	})
}
