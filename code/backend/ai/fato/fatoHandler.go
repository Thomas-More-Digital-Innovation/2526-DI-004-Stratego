package fato

import (
	"digital-innovation/stratego/engine"
)

type FatoAIhandler struct {
	ai FatoAI
}

func NewFatoAIhandler(player *engine.Player) FatoAIhandler {
	ai := NewFatoAI(player)
	return FatoAIhandler{
		ai: *ai,
	}
}

func (handler *FatoAIhandler) FafoTurnHandler(board *engine.Board, opponentMove engine.Move, opponent *engine.Player) engine.Move {
	// Analyze opponent's last move to gather intelligence
	if opponentMove.GetFrom().X >= 0 { // Valid move
		handler.ai.AnalyzeMove(opponentMove, opponent)
	}

	// Make AI's move
	return handler.ai.MakeMove(board)
}
