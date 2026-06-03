// Package mcts implements the Monte Carlo Tree Search strategy.
package mcts

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"math/rand/v2"
)

// AI implements the Monte Carlo Tree Search strategy.
type AI struct {
	ai.BaseAI
	params *ai.Parameters
}

// NewAI creates a new MCTS AI instance.
func NewAI(player *game.Player, hasMemory bool) *AI {
	params, _ := ai.Load("mcts", "default")
	return &AI{
		BaseAI: *ai.NewBaseAI(player, hasMemory),
		params: params,
	}
}

// NewAIWithParams creates a new MCTS AI instance with predefined params.
func NewAIWithParams(player *game.Player, hasMemory bool, params *ai.Parameters) *AI {
	return &AI{
		BaseAI: *ai.NewBaseAI(player, hasMemory),
		params: params,
	}
}

// MakeMove implements the player controller interface selecting the best move using MCTS rollouts.
func (aiObj *AI) MakeMove(board *game.Board) game.Move {
	opponent := ai.GetOpponent(board, aiObj.GetPlayer().GetID())
	moves := ai.GetMoves(board, aiObj.GetPlayer())
	if len(moves) == 0 {
		return game.Move{}
	}

	iterationsVal, ok := aiObj.params.Config["iterations"].(float64)
	iterations := 50
	if ok {
		iterations = int(iterationsVal)
	}

	rand.Shuffle(len(moves), func(i, j int) {
		moves[i], moves[j] = moves[j], moves[i]
	})

	var bestMove game.Move
	bestWinRate := -1.0

	for _, move := range moves {
		simulated := ai.SimulateMove(board, move)
		totalScore := 0.0

		for range iterations {
			totalScore += aiObj.rollout(simulated, aiObj.GetPlayer(), opponent)
		}

		winRate := totalScore / float64(iterations)
		if winRate > bestWinRate {
			bestWinRate = winRate
			bestMove = move
		}
	}

	return bestMove
}

func (aiObj *AI) rollout(board *game.Board, ourPlayer *game.Player, opponent *game.Player) float64 {
	tempBoard := board.Clone()
	currentPlayer := ourPlayer
	nextPlayer := opponent

	maxRolloutDepth := 10
	for range maxRolloutDepth {
		if aiObj.isFlagCaptured(tempBoard, opponent) {
			return 1.0
		}
		if aiObj.isFlagCaptured(tempBoard, ourPlayer) {
			return 0.0
		}

		moves := ai.GetMoves(tempBoard, currentPlayer)
		if len(moves) == 0 {
			if currentPlayer.GetID() == ourPlayer.GetID() {
				return 0.0
			}
			return 1.0
		}

		//nolint:gosec
		move := moves[rand.IntN(len(moves))]
		tempBoard = ai.SimulateMove(tempBoard, move)

		currentPlayer, nextPlayer = nextPlayer, currentPlayer
	}

	eval := ai.EvaluateBoard(tempBoard, ourPlayer, aiObj.GetMemory(), aiObj.params.Weights, aiObj.params.Aggression)
	if eval > 10.0 {
		return 1.0
	} else if eval < -10.0 {
		return 0.0
	}
	return 0.5
}

func (aiObj *AI) isFlagCaptured(board *game.Board, player *game.Player) bool {
	if player == nil {
		return false
	}
	for y := range 10 {
		for x := range 10 {
			piece := board.GetPieceAt(game.NewPosition(x, y))
			if piece != nil && piece.GetType().GetName() == "Flag" && piece.GetOwner().GetID() == player.GetID() {
				return false
			}
		}
	}
	return true
}
