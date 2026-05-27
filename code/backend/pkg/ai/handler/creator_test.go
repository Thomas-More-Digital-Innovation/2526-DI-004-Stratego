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
	aiFafo := AIhandler.CreateAI(models.Fafo, &player)
	assert.NotNil(t, aiFafo)

	// Test fato creation
	aiFato := AIhandler.CreateAI(models.Fato, &player)
	assert.NotNil(t, aiFato)

	// TODO: update with new AI types

	// Test panic for unknown AI
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Contains(t, r.(string), "I don't know that AI!")
	}()
	AIhandler.CreateAI("unknown-type", &player)
}
