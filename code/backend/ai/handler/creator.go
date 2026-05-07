// Package AIhandler provides a factory for creating AI instances
package AIhandler

import (
	"digital-innovation/stratego/ai"
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/ai/fato"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
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
