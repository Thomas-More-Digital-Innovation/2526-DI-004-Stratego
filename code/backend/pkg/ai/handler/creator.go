// Package AIhandler provides a factory for creating AI instances
package AIhandler

import (
	"digital-innovation/gostrategy/pkg/ai"
	"digital-innovation/gostrategy/pkg/ai/fafo"
	"digital-innovation/gostrategy/pkg/ai/fato"
	"digital-innovation/gostrategy/pkg/game"
	"digital-innovation/gostrategy/pkg/game/models"
)

// CreateAI is a factory function that returns an AI instance based on the given type
func CreateAI(ai string, player *game.Player) ai.AI {
	switch ai {
	case models.Fafo:
		return fafo.NewAI(player, false)
	case models.Fato:
		return fato.NewAI(player, true)
	default:
		panic("I don't know that AI! " + ai)
	}
}
