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
