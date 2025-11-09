package fafo

import (
	"digital-innovation/stratego/ai"
	"digital-innovation/stratego/engine"
	"math/rand/v2"
)

type FafoAI struct {
	ai.BaseAI
}

func NewFafoAI(player *engine.Player) *FafoAI {
	return &FafoAI{
		*ai.NewBaseAI(player),
	}
}

func (ai *FafoAI) PickRandomPiece() *engine.Piece {
	pieces := ai.GetPlayer().GetAlivePieces()
	if len(pieces) == 0 {
		return nil
	}
	random := rand.IntN(len(pieces))
	return pieces[random]
}

func (ai *FafoAI) MakeMove(board *engine.Board) engine.Move {
	return ai.FindRandomMove(board)
}

// findRandomMove picks any valid move as last resort
func (ai *FafoAI) FindRandomMove(board *engine.Board) engine.Move {
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
		if err != nil || len(moves) == 0 {
			continue
		}

		chosen := moves[rand.IntN(len(moves))]
		return engine.NewMove(chosen.GetFrom(), chosen.GetTo(), ai.GetPlayer())
	}

	// No valid moves available - player has lost (only immobile pieces left)
	// Return empty move to signal defeat
	return engine.Move{}
}
