package fato

import (
	ai "digital-innovation/stratego/ai"
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"math"
	"math/rand/v2"
)

// TODO: choose aggression on frontend

type FatoAI struct {
	fafo.FafoAI
	aggression float64
}

func NewFatoAI(player *engine.Player, hasMemory bool) *FatoAI {
	return NewFatoAIWithAggression(player, hasMemory, 0.5)
}

func NewFatoAIWithAggression(player *engine.Player, hasMemory bool, aggression float64) *FatoAI {
	fafoAI := fafo.NewFafoAI(player, hasMemory)

	if aggression < 0.0 {
		aggression = 0.0
	}
	if aggression > 1.0 {
		aggression = 1.0
	}

	return &FatoAI{
		FafoAI:     *fafoAI,
		aggression: aggression,
	}
}

// SetAggression sets the aggression level (0.0 = passive, 1.0 = aggressive)
func (ai *FatoAI) SetAggression(aggression float64) {
	if aggression < 0.0 {
		aggression = 0.0
	}
	if aggression > 1.0 {
		aggression = 1.0
	}
	ai.aggression = aggression
}

// GetAggression returns the current aggression level
func (ai *FatoAI) GetAggression() float64 {
	return ai.aggression
}

func (ai *FatoAI) MakeMove(board *engine.Board) engine.Move {
	// Not so random huh? :-)

	// 1. Try to attack a known enemy piece
	if move, found := ai.findAttackMove(board); found {
		return move
	}

	// 2. Try to explore toward enemy territory
	if move, found := ai.findExplorationMove(board); found {
		return move
	}

	// 3. Fallback: random valid move
	return ai.FindRandomMove(board)
}

// findAttackMove looks for moves that attack known/visible enemy pieces
func (ai *FatoAI) findAttackMove(board *engine.Board) (engine.Move, bool) {
	memory := ai.GetMemory()
	pieces := ai.GetPlayer().GetAlivePieces()

	// Shuffle pieces to add variety
	shuffled := make([]*engine.Piece, len(pieces))
	copy(shuffled, pieces)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	var bestAttack *engine.Move
	bestScore := -1000.0

	for _, piece := range shuffled {
		if !piece.CanMove() {
			continue
		}
		pos, exists := ai.GetPlayer().GetPiecePosition(piece)
		if !exists {
			continue
		}

		moves, err := board.ListMoves(pos)
		if err != nil {
			continue
		}

		for _, move := range moves {
			target := board.GetPieceAt(move.GetTo())
			if target != nil && target.GetOwner() != ai.GetPlayer() {
				score := ai.evaluateAttack(piece, target, move.GetTo(), memory)

				// Aggression determines minimum acceptable score
				// Passive (0.0): only take very favorable trades (> 0)
				// Moderate (0.5): take slightly unfavorable trades (> -50)
				// Aggressive (1.0): take any trade except terrible ones (> -100)
				minScore := -100.0*ai.aggression + 0.0*(1.0-ai.aggression)

				if score > bestScore && score > minScore {
					bestScore = score
					m := engine.NewMove(move.GetFrom(), move.GetTo(), ai.GetPlayer())
					bestAttack = &m
				}
			}
		}
	}

	if bestAttack != nil {
		return *bestAttack, true
	}
	return engine.Move{}, false
}

// evaluateAttack scores an attack opportunity
func (ai *FatoAI) evaluateAttack(attacker *engine.Piece, target *engine.Piece, targetPos engine.Position, memory *ai.AIMemory) float64 {
	score := 0.0

	// Check memory for target
	remembered := memory.Recall(targetPos)

	switch {
	case target.IsRevealed():
		rankDiff := float64(attacker.GetRank() - target.GetRank())
		score = rankDiff * 10

		if target.GetType().GetName() == "Flag" {
			score += 10000
		}
		if target.GetType().GetName() == "Bomb" && attacker.GetType().GetName() != "Miner" {
			score -= 1000
		}
	case remembered != nil && remembered.Confidence > 0.5:
		rankDiff := float64(attacker.GetRank() - remembered.Piece.GetRank())
		score = rankDiff * 10 * remembered.Confidence

		if remembered.Confidence < 0.8 {
			score *= 0.7
		}
	default:
		switch {
		case attacker.GetRank() >= 7:
			score = 20.0
		case attacker.GetRank() >= 5:
			score = 10.0
		case attacker.GetRank() >= 3:
			score = 5.0
		}
		score += rand.Float64()*10 - 5
	}
	return score
}

// findExplorationMove moves toward enemy side
func (ai *FatoAI) findExplorationMove(board *engine.Board) (engine.Move, bool) {
	enemyY := 0
	if ai.GetPlayer().GetID() == 1 {
		enemyY = 9
	}

	pieces := ai.GetPlayer().GetAlivePieces()
	shuffled := make([]*engine.Piece, len(pieces))
	copy(shuffled, pieces)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	for _, piece := range shuffled {
		if !piece.CanMove() {
			continue
		}
		pos, exists := ai.GetPlayer().GetPiecePosition(piece)
		if !exists {
			continue
		}

		moves, err := board.ListMoves(pos)
		if err != nil {
			continue
		}

		var bestMove *engine.Move
		bestDist := 100
		for _, move := range moves {
			if board.GetPieceAt(move.GetTo()) == nil {
				dist := int(math.Abs(float64(move.GetTo().Y - enemyY)))
				if dist < bestDist {
					bestDist = dist
					m := engine.NewMove(move.GetFrom(), move.GetTo(), ai.GetPlayer())
					bestMove = &m
				}
			}
		}
		if bestMove != nil {
			return *bestMove, true
		}
	}
	return engine.Move{}, false
}

// AnalyzeMove observes opponent moves and updates memory
// Overrides BaseAI to add scout detection
func (ai *FatoAI) AnalyzeMove(opponentMove engine.Move, opponent *engine.Player, round int) {
	memory := ai.GetMemory()
	from := opponentMove.GetFrom()
	to := opponentMove.GetTo()

	// First, apply default memory updates (move tracking)
	if memory.Recall(from) != nil {
		memory.MovePiece(from, to)
	}

	// Detect scout moves (moving >1 square in straight line)
	deltaX := int(math.Abs(float64(from.X - to.X)))
	deltaY := int(math.Abs(float64(from.Y - to.Y)))

	if deltaX > 1 || deltaY > 1 {
		scoutPiece := engine.NewPiece(models.Scout, opponent)
		memory.Remember(to, scoutPiece, 1.0, round)
	}
}
