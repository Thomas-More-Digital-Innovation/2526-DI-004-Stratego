// Package minimax implements the Minimax search strategy.
package minimax

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"math/rand/v2"
)

// AI implements the minimax search controller.
type AI struct {
	ai.BaseAI
	params *ai.Parameters
}

// NewAI creates a new Minimax AI instance.
func NewAI(player *game.Player, hasMemory bool) *AI {
	params, _ := ai.Load("minimax", "default")
	return &AI{
		BaseAI: *ai.NewBaseAI(player, hasMemory),
		params: params,
	}
}

// NewAIWithParams creates a new Minimax AI instance with predefined params.
func NewAIWithParams(player *game.Player, hasMemory bool, params *ai.Parameters) *AI {
	return &AI{
		BaseAI: *ai.NewBaseAI(player, hasMemory),
		params: params,
	}
}

// MakeMove implements the player controller interface selecting the best minimax search move.
func (aiObj *AI) MakeMove(board *game.Board) game.Move {
	depthVal, ok := aiObj.params.Config["depth"].(float64)
	depth := 2
	if ok {
		depth = int(depthVal)
	}

	detBoard := ai.DeterminizeBoard(board, aiObj.GetPlayer(), aiObj.GetMemory())
	opponent := ai.GetOpponent(detBoard, aiObj.GetPlayer().GetID())
	moves := ai.GetMoves(detBoard, aiObj.GetPlayer())
	if len(moves) == 0 {
		return game.Move{}
	}

	rand.Shuffle(len(moves), func(i, j int) {
		moves[i], moves[j] = moves[j], moves[i]
	})

	var bestMove game.Move
	bestScore := -1e9

	for _, move := range moves {
		simulated := ai.SimulateMove(detBoard, move)
		score := aiObj.minimax(simulated, depth-1, -1e9, 1e9, false, opponent)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return bestMove
}

func (aiObj *AI) minimax(board *game.Board, depth int, alpha, beta float64, maximizingPlayer bool, opponent *game.Player) float64 {
	if depth == 0 {
		return ai.EvaluateBoard(board, aiObj.GetPlayer(), aiObj.GetMemory(), aiObj.params.Weights, aiObj.params.Aggression)
	}

	if maximizingPlayer {
		maxEval := -1e9
		moves := ai.GetMoves(board, aiObj.GetPlayer())
		if len(moves) == 0 {
			return -1e9
		}
		for _, move := range moves {
			simulated := ai.SimulateMove(board, move)
			eval := aiObj.minimax(simulated, depth-1, alpha, beta, false, opponent)
			if eval > maxEval {
				maxEval = eval
			}
			if eval > alpha {
				alpha = eval
			}
			if beta <= alpha {
				break
			}
		}
		return maxEval
	}

	minEval := 1e9
	if opponent == nil {
		return ai.EvaluateBoard(board, aiObj.GetPlayer(), aiObj.GetMemory(), aiObj.params.Weights, aiObj.params.Aggression)
	}
	moves := ai.GetMoves(board, opponent)
	if len(moves) == 0 {
		return 1e9
	}
	for _, move := range moves {
		simulated := ai.SimulateMove(board, move)
		eval := aiObj.minimax(simulated, depth-1, alpha, beta, true, opponent)
		if eval < minEval {
			minEval = eval
		}
		if eval < beta {
			beta = eval
		}
		if beta <= alpha {
			break
		}
	}
	return minEval
}
