// Package heuristic implements the heuristic evaluation AI.
package heuristic

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"math/rand/v2"
)

// AI implements the heuristic evaluation strategy.
type AI struct {
	ai.BaseAI
	params *ai.Parameters
}

// NewAI creates a new Heuristic AI instance.
func NewAI(player *game.Player, hasMemory bool) *AI {
	params, _ := ai.Load(models.Heuristic, "default")
	return &AI{
		BaseAI: *ai.NewBaseAI(player, hasMemory),
		params: params,
	}
}

// NewAIWithParams creates a new Heuristic AI instance with predefined params.
func NewAIWithParams(player *game.Player, hasMemory bool, params *ai.Parameters) *AI {
	return &AI{
		BaseAI: *ai.NewBaseAI(player, hasMemory),
		params: params,
	}
}

// MakeMove implements the player controller interface choosing the highest heuristic move.
func (aiObj *AI) MakeMove(board *game.Board) game.Move {
	detBoard := ai.DeterminizeBoard(board, aiObj.GetPlayer(), aiObj.GetMemory())
	pieces := aiObj.GetPlayer().GetAlivePieces()
	var allMoves []game.Move

	for _, piece := range pieces {
		if !piece.CanMove() {
			continue
		}
		pos, exists := aiObj.GetPlayer().GetPiecePosition(piece)
		if !exists {
			continue
		}
		moves, err := detBoard.ListMoves(pos)
		if err != nil {
			continue
		}
		for _, m := range moves {
			allMoves = append(allMoves, game.NewMove(m.GetFrom(), m.GetTo(), aiObj.GetPlayer()))
		}
	}

	if len(allMoves) == 0 {
		return game.Move{}
	}

	rand.Shuffle(len(allMoves), func(i, j int) {
		allMoves[i], allMoves[j] = allMoves[j], allMoves[i]
	})

	var bestMove game.Move
	bestScore := -1e9

	for _, move := range allMoves {
		simulated := ai.SimulateMove(detBoard, move)
		score := ai.EvaluateBoard(simulated, aiObj.GetPlayer(), aiObj.GetMemory(), aiObj.params.Weights, aiObj.params.Aggression)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return bestMove
}
