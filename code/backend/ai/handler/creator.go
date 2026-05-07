// Package AIhandler provides a factory for creating AI instances
package AIhandler

import (
	"digital-innovation/gostrategy/ai"
	"digital-innovation/gostrategy/ai/fafo"
	"digital-innovation/gostrategy/ai/fato"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/models"
)

// CreateAI is a factory function that returns an AI instance based on the given type
func CreateAI(ai string, player *engine.Player) ai.AI {
	switch ai {
	case models.Fafo:
		return fafo.NewAI(player, false)
	case models.Fato:
		return fato.NewAI(player, true)
	default:
		panic("I don't know that AI! " + ai)
	}
}
