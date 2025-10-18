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

func (m *Move) GetFrom() Position {
	return m.from
}

func (m *Move) GetTo() Position {
	return m.to
}

func (m *Move) GetPlayer() *Player {
	return m.player
}

func (m *Move) String() string {
	return m.from.String() + " -> " + m.to.String()
}
