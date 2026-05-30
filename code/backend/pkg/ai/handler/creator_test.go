package AIhandler_test

import (
	AIhandler "digital-innovation/gostrategy/pkg/ai/handler"
	"digital-innovation/gostrategy/pkg/game"
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAI(t *testing.T) {
	player := game.NewPlayer(0, "test-ai", "red")

	// Test fafo creation
	aiFafo, err := AIhandler.CreateAI(models.Fafo, &player)
	assert.NoError(t, err)
	assert.NotNil(t, aiFafo)

	// Test fato creation
	aiFato, err := AIhandler.CreateAI(models.Fato, &player)
	assert.NoError(t, err)
	assert.NotNil(t, aiFato)

	// Test error returned for unknown AI type
	_, err = AIhandler.CreateAI("unknown-type", &player)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "i don't know that AI!")
}
