package fafo

import (
	"digital-innovation/stratego/engine"
	"math/rand/v2"
)

type FafoAI struct {
	player *engine.Player
}

func NewFafoAI(player *engine.Player) *FafoAI {
	return &FafoAI{
		player: player,
	}
}

func (ai *FafoAI) GetPlayer() *engine.Player {
	return ai.player
}

func (ai *FafoAI) GetControllerType() engine.ControllerType {
	return engine.AIController
}

func (ai *FafoAI) PickRandomPiece() *engine.Piece {
	pieces := ai.player.GetAlivePieces()
	if len(pieces) == 0 {
		return nil
	}
	random := rand.IntN(len(pieces))
	return pieces[random]
}

func (ai *FafoAI) MakeMove(board *engine.Board) engine.Move {
	return ai.findRandomMove(board)
}

// findRandomMove picks any valid move as last resort
func (ai *FafoAI) findRandomMove(board *engine.Board) engine.Move {
	pieces := ai.player.GetAlivePieces()
	shuffled := make([]*engine.Piece, len(pieces))
	copy(shuffled, pieces)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	for _, piece := range shuffled {
		if !piece.CanMove() {
			continue
		}
		pos, exists := ai.player.GetPiecePosition(piece)
		if !exists {
			continue
		}

		moves, err := board.ListMoves(pos)
		if err != nil || len(moves) == 0 {
			continue
		}
		return moves[rand.IntN(len(moves))]
	}

	// No valid moves available - player has lost (only immobile pieces left)
	// Return empty move to signal defeat
	return engine.Move{}
}
