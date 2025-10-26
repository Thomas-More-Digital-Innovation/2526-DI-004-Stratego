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

func (handler *FafoAIhandler) FafoTurnHandler(board *engine.Board, opponentMove engine.Move, opponent *engine.Player) engine.Move {
	// Analyze opponent's last move to gather intelligence
	if opponentMove.GetFrom().X >= 0 { // Valid move
		handler.ai.AnalyzeMove(opponentMove, opponent)
	}

	// Make AI's move
	return handler.ai.MakeMove(board)
}
