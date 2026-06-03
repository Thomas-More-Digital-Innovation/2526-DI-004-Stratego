package game

// Export internal fields for testing
func (gs *Session) GetPlayer1Pieces() []*Piece {
	return gs.player1Pieces
}

func (gs *Session) GetPlayer2Pieces() []*Piece {
	return gs.player2Pieces
}
