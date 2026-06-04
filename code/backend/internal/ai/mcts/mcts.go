// Package mcts implements the Monte Carlo Tree Search strategy.
package mcts

import (
	"digital-innovation/gostrategy/internal/ai"
	ai_const "digital-innovation/gostrategy/internal/ai/const"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"math/rand/v2"
	"runtime"
	"sync"
)

// AI implements the Monte Carlo Tree Search strategy.
type AI struct {
	ai.BaseAI
	params *ai.Parameters
}

// NewAI creates a new MCTS AI instance.
func NewAI(player *game.Player, hasMemory bool) *AI {
	params, _ := ai.Load(models.Mcts, "default")
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
// Candidate moves are evaluated in parallel, bounded by GOMAXPROCS to avoid scheduler thrashing.
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

	//nolint:gosec
	rand.Shuffle(len(moves), func(i, j int) {
		moves[i], moves[j] = moves[j], moves[i]
	})

	scores := make([]float64, len(moves))
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.GOMAXPROCS(0))

	for i, move := range moves {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, m game.Move) {
			defer wg.Done()
			defer func() { <-sem }()
			simulated := ai.SimulateMove(board, m)
			total := 0.0
			for range iterations {
				total += aiObj.rollout(simulated, aiObj.GetPlayer(), opponent)
			}
			scores[idx] = total / float64(iterations)
		}(i, move)
	}

	wg.Wait()

	bestMove := moves[0]
	bestRate := scores[0]
	for i := 1; i < len(scores); i++ {
		if scores[i] > bestRate {
			bestRate = scores[i]
			bestMove = moves[i]
		}
	}
	return bestMove
}

func (aiObj *AI) rollout(board *game.Board, ourPlayer *game.Player, opponent *game.Player) float64 {
	tempBoard := ai.DeterminizeBoard(board, ourPlayer, aiObj.GetMemory())
	currentPlayer := ourPlayer
	nextPlayer := opponent

	ourFlagPos := aiObj.findFlagPosition(tempBoard, ourPlayer)
	oppFlagPos := aiObj.findFlagPosition(tempBoard, opponent)

	// Build active piece index once per rollout — avoids 100-tile board scan on every depth step.
	ourIndex := ai.BuildMobileIndex(tempBoard, ourPlayer)
	var oppIndex []game.Position
	if opponent != nil {
		oppIndex = ai.BuildMobileIndex(tempBoard, opponent)
	}

	captureOccurred := false

	maxRolloutDepth := 10
	for range maxRolloutDepth {
		if aiObj.isFlagCapturedAt(tempBoard, oppFlagPos, opponent) {
			return 1.0
		}
		if aiObj.isFlagCapturedAt(tempBoard, ourFlagPos, ourPlayer) {
			return 0.0
		}

		var currentIndex *[]game.Position
		if currentPlayer.GetID() == ourPlayer.GetID() {
			currentIndex = &ourIndex
		} else {
			currentIndex = &oppIndex
		}

		moves := ai.GetMovesFromIndex(tempBoard, *currentIndex, currentPlayer)
		if len(moves) == 0 {
			if currentPlayer.GetID() == ourPlayer.GetID() {
				return 0.0
			}
			return 1.0
		}

		move := pickRolloutMove(tempBoard, moves)
		captured := aiObj.applySimulatedMoveInPlace(tempBoard, move)
		if captured {
			captureOccurred = true
		}

		// Update index: remove old position, add new if piece survived.
		updateIndex(currentIndex, move.GetFrom(), move.GetTo(), captured)

		currentPlayer, nextPlayer = nextPlayer, currentPlayer
	}

	// Skip EvaluateBoard for non-decisive rollouts — returns indeterminate 0.5 directly.
	if !captureOccurred {
		return 0.5
	}

	eval := ai.EvaluateBoard(tempBoard, ourPlayer, aiObj.GetMemory(), aiObj.params.Weights, aiObj.params.Aggression)
	if eval > 10.0 {
		return 1.0
	} else if eval < -10.0 {
		return 0.0
	}
	return 0.5
}

// pickRolloutMove biases rollout selection toward capture moves (80% preference when available).
func pickRolloutMove(board *game.Board, moves []game.Move) game.Move {
	var captures []game.Move
	for _, m := range moves {
		if board.GetPieceAt(m.GetTo()) != nil {
			captures = append(captures, m)
		}
	}

	//nolint:gosec
	if len(captures) > 0 && rand.Float64() < 0.8 {
		return captures[rand.IntN(len(captures))]
	}
	//nolint:gosec
	return moves[rand.IntN(len(moves))]
}

// updateIndex mutates the position index in-place after a move is applied.
func updateIndex(index *[]game.Position, from, to game.Position, captured bool) {
	s := *index
	for i, pos := range s {
		if pos == from {
			if captured {
				// piece was lost — remove from index
				s[i] = s[len(s)-1]
				*index = s[:len(s)-1]
			} else {
				s[i] = to
			}
			return
		}
	}
}

func (aiObj *AI) findFlagPosition(board *game.Board, player *game.Player) game.Position {
	if player == nil {
		return game.NewPosition(-1, -1)
	}
	for y := range 10 {
		for x := range 10 {
			pos := game.NewPosition(x, y)
			piece := board.GetPieceAt(pos)
			if piece != nil && piece.GetType().GetName() == ai_const.Flag && piece.GetOwner().GetID() == player.GetID() {
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
	return piece == nil || piece.GetType().GetName() != ai_const.Flag || piece.GetOwner().GetID() != player.GetID()
}

// applySimulatedMoveInPlace applies a move directly to the board without cloning.
// Returns true if the attacker was captured (lost the fight) or if a piece was removed from that position.
func (aiObj *AI) applySimulatedMoveInPlace(b *game.Board, move game.Move) bool {
	attacker := b.GetPieceAt(move.GetFrom())
	if attacker == nil {
		return false
	}

	target := b.GetPieceAt(move.GetTo())
	if target == nil {
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
		return false
	}

	attackerRank := attacker.GetRank()
	defenderRank := target.GetRank()

	if defenderRank == models.Flag.GetRank() {
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
		return false
	}

	if attackerRank == models.Spy.GetRank() && defenderRank == models.Marshal.GetRank() {
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
		return false
	}

	if defenderRank == models.Bomb.GetRank() {
		if attacker.GetType().GetName() == ai_const.Miner {
			b.SetPieceAt(move.GetFrom(), nil)
			b.SetPieceAt(move.GetTo(), attacker)
			return false
		}
		b.SetPieceAt(move.GetFrom(), nil)
		return true // attacker destroyed
	}

	switch {
	case attackerRank > defenderRank:
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
		return false
	case attackerRank < defenderRank:
		b.SetPieceAt(move.GetFrom(), nil)
		return true // attacker destroyed
	default:
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), nil)
		return true // both destroyed — attacker gone from index
	}
}
