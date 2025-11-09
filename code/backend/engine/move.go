package engine

type Move struct {
	from   Position
	to     Position
	player *Player
}

func NewMove(from Position, to Position, player *Player) Move {
	return Move{
		from:   from,
		to:     to,
		player: player,
	}
}

// GetFrom returns the starting position of the move.
func (m *Move) GetFrom() Position {
	return m.from
}

// GetTo returns the ending position of the move.
func (m *Move) GetTo() Position {
	return m.to
}

// GetPlayer returns the player that made the move.
func (m *Move) GetPlayer() *Player {
	return m.player
}

// IsEmpty returns a boolean indicating whether the move is empty (i.e., has no player).
// An empty move is used to indicate that no move has been made, such as when the game is in a waiting state.
// The method returns true if the move is empty, and false otherwise.
func (m *Move) IsEmpty() bool {
	return m.player == nil
}

// String returns a string representation of the move in the format "from -> to".
func (m *Move) String() string {
	return m.from.String() + " -> " + m.to.String()
}
