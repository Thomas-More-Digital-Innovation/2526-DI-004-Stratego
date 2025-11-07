package fafo

import (
	"digital-innovation/stratego/engine"
)

type FafoAIhandler struct {
	ai FafoAI
}

func NewFafoAIhandler(player *engine.Player) FafoAIhandler {
	ai := NewFafoAI(player)
	return FafoAIhandler{
		ai: *ai,
	}
}

func (handler *FafoAIhandler) FafoTurnHandler(board *engine.Board) engine.Move {
	// Make AI's move
	return handler.ai.MakeMove(board)
}
