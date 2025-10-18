package engine

type Position struct {
	X int
	Y int
}

func NewPosition(x int, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

func (pos Position) Equals(other Position) bool {
	return pos.X == other.X && pos.Y == other.Y
}

func (pos Position) Copy() Position {
	return Position{
		X: pos.X,
		Y: pos.Y,
	}
}

func (pos Position) ToLeft() Position {
	return Position{
		X: pos.X - 1,
		Y: pos.Y,
	}
}

func (pos Position) ToRight() Position {
	return Position{
		X: pos.X + 1,
		Y: pos.Y,
	}
}

func (pos Position) ToUp() Position {
	return Position{
		X: pos.X,
		Y: pos.Y - 1,
	}
}

func (pos Position) ToDown() Position {
	return Position{
		X: pos.X,
		Y: pos.Y + 1,
	}
}

func (pos Position) String() string {
	return "(" + string(rune(pos.X+'A')) + "," + string(rune(pos.Y+'0')) + ")"
}
