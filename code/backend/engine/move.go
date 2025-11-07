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

// String returns a string representation of the move in the format "from -> to".
func (m *Move) String() string {
	return m.from.String() + " -> " + m.to.String()
}
