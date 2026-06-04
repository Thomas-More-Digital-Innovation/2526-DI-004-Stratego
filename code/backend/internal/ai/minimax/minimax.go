// Package minimax implements the Minimax search strategy.
package minimax

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"slices"
)

const (
	ttExact = 0
	ttAlpha = 1
	ttBeta  = 2
)

type ttEntry struct {
	score float64
	depth int
	flag  int
}

// AI implements the minimax search controller.
type AI struct {
	ai.BaseAI
	params *ai.Parameters
	tt     map[string]ttEntry
}

// NewAI creates a new Minimax AI instance.
func NewAI(player *game.Player, hasMemory bool) *AI {
	params, _ := ai.Load(models.Minimax, "default")
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

	aiObj.tt = make(map[string]ttEntry)

	detBoard := ai.DeterminizeBoard(board, aiObj.GetPlayer(), aiObj.GetMemory())
	opponent := ai.GetOpponent(detBoard, aiObj.GetPlayer().GetID())
	moves := ai.GetMoves(detBoard, aiObj.GetPlayer())
	if len(moves) == 0 {
		return game.Move{}
	}

	orderMoves(detBoard, moves)

	var bestMove game.Move
	bestScore := -1e9

	for _, move := range moves {
		simulated := ai.SimulateMove(detBoard, move)
		score := aiObj.minimax(simulated, depth-1, bestScore, 1e9, false, opponent)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return bestMove
}

func (aiObj *AI) minimax(board *game.Board, depth int, alpha, beta float64, maximizingPlayer bool, opponent *game.Player) float64 {
	originalAlpha := alpha

	key := getBoardStateKey(board, maximizingPlayer)
	if entry, ok := aiObj.tt[key]; ok && entry.depth >= depth {
		switch entry.flag {
		case ttExact:
			return entry.score
		case ttAlpha:
			if entry.score < beta {
				beta = entry.score
			}
		case ttBeta:
			if entry.score > alpha {
				alpha = entry.score
			}
		}
		if alpha >= beta {
			return entry.score
		}
	}

	if depth == 0 {
		val := ai.EvaluateBoard(board, aiObj.GetPlayer(), aiObj.GetMemory(), aiObj.params.Weights, aiObj.params.Aggression)
		aiObj.tt[key] = ttEntry{
			score: val,
			depth: depth,
			flag:  ttExact,
		}
		return val
	}

	if maximizingPlayer {
		maxEval := -1e9
		moves := ai.GetMoves(board, aiObj.GetPlayer())
		if len(moves) == 0 {
			return -1e9
		}

		orderMoves(board, moves)

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

		flag := ttExact
		if maxEval <= originalAlpha {
			flag = ttAlpha
		} else if maxEval >= beta {
			flag = ttBeta
		}
		aiObj.tt[key] = ttEntry{
			score: maxEval,
			depth: depth,
			flag:  flag,
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

	orderMoves(board, moves)

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

	flag := ttExact
	if minEval <= originalAlpha {
		flag = ttAlpha
	} else if minEval >= beta {
		flag = ttBeta
	}
	aiObj.tt[key] = ttEntry{
		score: minEval,
		depth: depth,
		flag:  flag,
	}
	return minEval
}

func getBoardStateKey(board *game.Board, maximizingPlayer bool) string {
	var key [101]byte
	idx := 0
	field := board.GetField()
	for y := range 10 {
		for x := range 10 {
			piece := field[y][x]
			if piece == nil {
				key[idx] = 0
			} else {
				ownerID := piece.GetOwner().GetID()
				rankVal := piece.GetRank()
				key[idx] = byte(ownerID&0xF)<<4 | byte(rankVal&0xF)
			}
			idx++
		}
	}
	if maximizingPlayer {
		key[100] = 1
	} else {
		key[100] = 0
	}
	return string(key[:])
}

func orderMoves(board *game.Board, moves []game.Move) {
	slices.SortFunc(moves, func(a, b game.Move) int {
		scoreA := scoreMove(board, a)
		scoreB := scoreMove(board, b)
		if scoreA > scoreB {
			return -1
		}
		if scoreA < scoreB {
			return 1
		}
		return 0
	})
}

func scoreMove(board *game.Board, move game.Move) float64 {
	score := 0.0
	target := board.GetPieceAt(move.GetTo())
	if target != nil {
		score += 100.0
		score += float64(target.GetRank()) * 10.0
	}
	return score
}
