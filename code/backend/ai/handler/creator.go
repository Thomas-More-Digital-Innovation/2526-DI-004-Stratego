package AIhandler

import (
	"digital-innovation/stratego/ai"
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/ai/fato"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
)

func CreateAI(ai string, player *engine.Player) ai.AI {
	switch ai {
	case models.Fafo:
		return fafo.NewFafoAI(player, false)
	case models.Fato:
		return fato.NewFatoAI(player, true)
	default:
		panic("I don't know that AI! " + ai)
	}
}
