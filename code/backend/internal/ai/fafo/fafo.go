// Package fafo implements the FAFO (Fuck Around and Find Out) AI strategy
// This AI prioritizes random moves and aggression over deep planning
package fafo

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/pkg/game"
	"math/rand/v2"
)

// AI implements a random movement strategy
type AI struct {
	ai.BaseAI
}

// NewAI creates a new AI instance
func NewAI(player *game.Player, hasMemory bool) *AI {
	return &AI{
		*ai.NewBaseAI(player, hasMemory),
	}
}

// PickRandomPiece selects a random piece from the player's alive pieces
func (ai *AI) PickRandomPiece() *game.Piece {
	pieces := ai.GetPlayer().GetAlivePieces()
	if len(pieces) == 0 {
		return nil
	}
	// #nosec G404 - weak random is sufficient for AI piece selection
	random := rand.IntN(len(pieces))
	return pieces[random]
}

// MakeMove implements the PlayerController interface
func (ai *AI) MakeMove(board *game.Board) game.Move {
	return ai.FindRandomMove(board)
}

// FindRandomMove picks any valid move for a random piece
func (ai *AI) FindRandomMove(board *game.Board) game.Move {
	pieces := ai.GetPlayer().GetAlivePieces()
	shuffled := make([]*game.Piece, len(pieces))
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
		if err != nil || len(moves) == 0 {
			continue
		}

		// #nosec G404 - weak random is sufficient for AI move selection
		chosen := moves[rand.IntN(len(moves))]
		return game.NewMove(chosen.GetFrom(), chosen.GetTo(), ai.GetPlayer())
	}

	// No valid moves available - player has lost (only immobile pieces left)
	// Return empty move to signal defeat
	return game.Move{}
}
