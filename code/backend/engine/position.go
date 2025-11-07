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

// Equals checks if two positions are equal
// It returns true if the X and Y coordinates are the same, false otherwise.
func (pos Position) Equals(other Position) bool {
	return pos.X == other.X && pos.Y == other.Y
}

// Copy returns a deep copy of the Position.
// It is useful for creating a new Position without modifying the original.
// It returns a new Position with the same X and Y coordinates as the original.
func (pos Position) Copy() Position {
	return Position{
		X: pos.X,
		Y: pos.Y,
	}
}

// ToLeft returns a new Position that is one unit to the left of the original.
// It is useful for generating possible moves for a piece.
// It returns a new Position with the same Y coordinate as the original, and an X coordinate that is one less than the original.
func (pos Position) ToLeft() Position {
	return Position{
		X: pos.X - 1,
		Y: pos.Y,
	}
}

// ToRight returns a new Position that is one unit to the right of the original.
// It is useful for generating possible moves for a piece.
// It returns a new Position with the same Y coordinate as the original, and an X coordinate that is one more than the original.
func (pos Position) ToRight() Position {
	return Position{
		X: pos.X + 1,
		Y: pos.Y,
	}
}

// ToUp returns a new Position that is one unit above the original.
// It is useful for generating possible moves for a piece.
// It returns a new Position with the same X coordinate as the original, and a Y coordinate that is one less than the original.
func (pos Position) ToUp() Position {
	return Position{
		X: pos.X,
		Y: pos.Y - 1,
	}
}

// ToDown returns a new Position that is one unit below the original.
// It is useful for generating possible moves for a piece.
// It returns a new Position with the same X coordinate as the original, and a Y coordinate that is one more than the original.
func (pos Position) ToDown() Position {
	return Position{
		X: pos.X,
		Y: pos.Y + 1,
	}
}

// String returns a string representation of the position in the format "(X,Y)"
// where X is a letter (A-H) and Y is a number (0-7).
// It is useful for debugging and logging.
// For example, the position (0,0) would be represented as "(A,0)".
func (pos Position) String() string {
	return "(" + string(rune(pos.X+'A')) + "," + string(rune(pos.Y+'0')) + ")"
}
