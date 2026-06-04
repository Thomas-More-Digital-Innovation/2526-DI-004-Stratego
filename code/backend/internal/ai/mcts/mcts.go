// Package mcts implements the Monte Carlo Tree Search strategy.
package mcts

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
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
	tempBoard := board.FastClone()
	currentPlayer := ourPlayer
	nextPlayer := opponent

	ourFlagPos := aiObj.findFlagPosition(tempBoard, ourPlayer)
	oppFlagPos := aiObj.findFlagPosition(tempBoard, opponent)

	maxRolloutDepth := 10
	for range maxRolloutDepth {
		if aiObj.isFlagCapturedAt(tempBoard, oppFlagPos, opponent) {
			return 1.0
		}
		if aiObj.isFlagCapturedAt(tempBoard, ourFlagPos, ourPlayer) {
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
		aiObj.applySimulatedMoveInPlace(tempBoard, move)

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

func (aiObj *AI) findFlagPosition(board *game.Board, player *game.Player) game.Position {
	if player == nil {
		return game.NewPosition(-1, -1)
	}
	for y := range 10 {
		for x := range 10 {
			pos := game.NewPosition(x, y)
			piece := board.GetPieceAt(pos)
			if piece != nil && piece.GetType().GetName() == "Flag" && piece.GetOwner().GetID() == player.GetID() {
				return pos
			}
		}
	}
	return game.NewPosition(-1, -1)
}

func (aiObj *AI) isFlagCapturedAt(board *game.Board, flagPos game.Position, player *game.Player) bool {
	if flagPos.X == -1 {
		return true
	}
	piece := board.GetPieceAt(flagPos)
	return piece == nil || piece.GetType().GetName() != "Flag" || piece.GetOwner().GetID() != player.GetID()
}

func (aiObj *AI) applySimulatedMoveInPlace(b *game.Board, move game.Move) {
	attacker := b.GetPieceAt(move.GetFrom())
	if attacker == nil {
		return
	}

	target := b.GetPieceAt(move.GetTo())
	if target != nil {
		attackerRank := attacker.GetRank()
		defenderRank := target.GetRank()

		if defenderRank == models.Flag.GetRank() {
			b.SetPieceAt(move.GetFrom(), nil)
			b.SetPieceAt(move.GetTo(), attacker)
			return
		}

		if attackerRank == models.Spy.GetRank() && defenderRank == models.Marshal.GetRank() {
			b.SetPieceAt(move.GetFrom(), nil)
			b.SetPieceAt(move.GetTo(), attacker)
			return
		}

		if defenderRank == models.Bomb.GetRank() {
			if attacker.GetType().GetName() == "Miner" {
				b.SetPieceAt(move.GetFrom(), nil)
				b.SetPieceAt(move.GetTo(), attacker)
			} else {
				b.SetPieceAt(move.GetFrom(), nil)
			}
			return
		}

		switch {
		case attackerRank > defenderRank:
			b.SetPieceAt(move.GetFrom(), nil)
			b.SetPieceAt(move.GetTo(), attacker)
		case attackerRank < defenderRank:
			b.SetPieceAt(move.GetFrom(), nil)
		default:
			b.SetPieceAt(move.GetFrom(), nil)
			b.SetPieceAt(move.GetTo(), nil)
		}
	} else {
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
	}
}
