package ai

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"math/rand/v2"
)

// DeterminizeBoard returns a board where all unrevealed opponent pieces are shuffled and randomized from the remaining alive unrevealed pool.
func DeterminizeBoard(board *game.Board, ourPlayer *game.Player, memory *Memory) *game.Board {
	opponent := GetOpponent(board, ourPlayer.GetID())
	if opponent == nil {
		return board.FastClone()
	}

	// 1. Gather all alive opponent pieces
	alivePieces := opponent.GetAlivePieces()
	if len(alivePieces) == 0 {
		return board.FastClone()
	}

	// 2. Identify which are revealed
	revealedPieces := make(map[*game.Piece]bool)
	var unrevealedTypes []models.PieceType

	for _, piece := range alivePieces {
		pos, exists := opponent.GetPiecePosition(piece)
		isKnown := piece.IsRevealed()
		if !isKnown && exists && memory != nil {
			if entry := memory.Recall(pos); entry != nil && entry.Confidence == 1.0 {
				isKnown = true
			}
		}

		if isKnown {
			revealedPieces[piece] = true
		} else {
			unrevealedTypes = append(unrevealedTypes, *piece.GetType())
		}
	}

	// Shuffle the unrevealed types pool
	rand.Shuffle(len(unrevealedTypes), func(i, j int) {
		unrevealedTypes[i], unrevealedTypes[j] = unrevealedTypes[j], unrevealedTypes[i]
	})

	// 3. Clone the board, and replace unrevealed opponent pieces
	determinized := board.FastClone()
	unrevealedIdx := 0

	for y := range 10 {
		for x := range 10 {
			pos := game.NewPosition(x, y)
			piece := determinized.GetPieceAt(pos)
			if piece == nil || piece.GetOwner().GetID() == ourPlayer.GetID() {
				continue
			}

			// It is an opponent piece. Is it revealed?
			if revealedPieces[piece] {
				continue
			}

			// Replace it with one from the unrevealedTypes pool
			if unrevealedIdx < len(unrevealedTypes) {
				sampledType := unrevealedTypes[unrevealedIdx]
				unrevealedIdx++
				newPiece := game.NewPiece(sampledType, opponent)
				determinized.SetPieceAt(pos, newPiece)
			}
		}
	}

	return determinized
}
