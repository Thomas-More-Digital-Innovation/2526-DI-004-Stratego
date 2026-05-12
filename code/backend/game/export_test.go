package game

import "digital-innovation/gostrategy/engine"

// Export internal fields for testing
func (gs *Session) GetPlayer1Pieces() []*engine.Piece {
	return gs.player1Pieces
}

func (gs *Session) GetPlayer2Pieces() []*engine.Piece {
	return gs.player2Pieces
}
