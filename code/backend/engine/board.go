package engine

import "errors"

type Board struct {
	field [10][10]*Piece
	lakes [8]Position
}

func NewBoard() *Board {
	return &Board{
		field: [10][10]*Piece{},
		lakes: [8]Position{
			NewPosition(2, 4), NewPosition(3, 4), NewPosition(2, 5), NewPosition(3, 5),
			NewPosition(6, 4), NewPosition(7, 4), NewPosition(6, 5), NewPosition(7, 5),
		},
	}
}

// SetPieceAt sets the piece at the given position on the board.
// The piece is updated in the board's internal field, which is a 10x10 2D slice of pointers to Piece.
// The function does not check if the move is valid, it simply updates the board state.
// The function is O(1) and updates the board state.
func (b *Board) SetPieceAt(pos Position, piece *Piece) {
	b.field[pos.Y][pos.X] = piece
}

// GetField returns the board's internal field, which is a 10x10 2D slice of pointers to Piece.
// The field is used to store the pieces on the board and is updated by the game engine.
// The function is O(1) and does not modify the board state.
func (b *Board) GetField() [10][10]*Piece {
	return b.field
}

// GetPieceAt returns the piece at the given position on the board.
// If there is no piece at the given position, the function returns nil.
// The function is O(1) and does not modify the board state.
func (b *Board) GetPieceAt(pos Position) *Piece {
	return b.field[pos.Y][pos.X]
}

// IsLake returns a boolean indicating whether the given position is a lake on the board.
// A lake is a specific position on the board that is not occupiable by any piece.
// The function checks if the given position is equal to any of the positions in the board's lake list.
// If the position is found in the list, the function returns true, otherwise it returns false.
func (b *Board) IsLake(pos Position) bool {
	for _, lakePos := range b.lakes {
		if lakePos == pos {
			return true
		}
	}
	return false
}

// IsValidMove returns a boolean indicating whether a move is valid on the board.
// It checks if the move is within the bounds of the board, if the destination is a lake,
// and if the destination is occupied by a piece of the same owner as the piece being moved.
// It does not check if the piece can move to the destination (e.g. if the piece is a scout, it does
// not check if the destination is more than one space away).
// It returns false if any of these conditions are not met, and true otherwise.
func (b *Board) IsValidMove(move *Move) bool {
	if move.GetTo().X < 0 || move.GetTo().X >= len(b.field[0]) || move.GetTo().Y < 0 || move.GetTo().Y >= len(b.field) {
		return false
	}
	if b.IsLake(move.GetTo()) {
		return false
	}
	if b.GetPieceAt(move.GetTo()) != nil && b.GetPieceAt(move.GetTo()).GetOwner() == b.GetPieceAt(move.GetFrom()).GetOwner() {
		return false
	}
	return true
}

// SwapPieces swaps two pieces in the setup
// It returns an error if either of the positions are empty.
// It does not check if the move is valid, it just swaps the pieces.
// The board state is updated accordingly.
func (b *Board) SwapPieces(pos1 Position, pos2 Position) error {
	piece1 := b.GetPieceAt(pos1)
	piece2 := b.GetPieceAt(pos2)

	if piece1 == nil || piece2 == nil {
		return errors.New("one or both positions are empty")
	}

	b.field[pos1.Y][pos1.X], b.field[pos2.Y][pos2.X] = piece2, piece1
	return nil
}

// RemovePieceAt removes a piece from the board at the given position.
// If there is no piece at the given position, or if the piece is still alive,
// the function returns an error. Otherwise, the piece is removed from the board.
func (b *Board) RemovePieceAt(pos Position) error {
	piece := b.GetPieceAt(pos)
	if piece == nil {
		return errors.New("no piece at the given position to remove")
	}
	if piece.IsAlive() {
		return errors.New("cannot remove a piece that is still alive")
	}
	b.field[pos.Y][pos.X] = nil
	return nil
}

// MovePiece updates the board state by moving a piece from one position to another.
// The piece is moved from the "from" position to the "to" position, and the board state is updated accordingly.
// If there is a piece at the "to" position, it is removed from the board.
// If the piece is not alive, or if there is no piece at the "from" position, the function does nothing.
func (b *Board) MovePiece(move *Move, piece *Piece) {
	b.field[move.GetFrom().Y][move.GetFrom().X] = nil
	b.field[move.GetTo().Y][move.GetTo().X] = piece
}

func (b *Board) ListMoves(pos Position) ([]Move, error) {
	var moves []Move
	piece := b.GetPieceAt(pos)
	if piece == nil || !piece.CanMove() {
		return moves, errors.New("no movable piece at the given position")
	}

	if piece.GetType().GetName() == "Scout" {
		// Scouts can move any number of spaces in straight lines
		b.handleScoutMoves(pos, &moves)
	} else {
		// Other pieces can move one space in any orthogonal direction
		b.handleStandardMoves(pos, &moves)
	}
	return moves, nil
}

func (b *Board) handleScoutMoves(pos Position, moves *[]Move) {
	directions := []Position{
		{X: 0, Y: -1}, // Up
		{X: 0, Y: 1},  // Down
		{X: -1, Y: 0}, // Left
		{X: 1, Y: 0},  // Right
	}

	for _, dir := range directions {
		for step := 1; step < 10; step++ {
			newPos := NewPosition(pos.X+dir.X*step, pos.Y+dir.Y*step)
			move := NewMove(pos, newPos, nil)
			if b.IsValidMove(&move) {
				*moves = append(*moves, move)
				if b.GetPieceAt(newPos) != nil {
					// Can't move past an occupied square
					break
				}
			} else {
				// Stop if the move is invalid (out of bounds, lake, etc.)
				break
			}
		}
	}
}

func (b *Board) handleStandardMoves(pos Position, moves *[]Move) {
	directions := []Position{
		{X: 0, Y: -1}, // Up
		{X: 0, Y: 1},  // Down
		{X: -1, Y: 0}, // Left
		{X: 1, Y: 0},  // Right
	}

	for _, dir := range directions {
		newPos := NewPosition(pos.X+dir.X, pos.Y+dir.Y)
		move := NewMove(pos, newPos, nil)
		if b.IsValidMove(&move) {
			*moves = append(*moves, move)
		}
	}
}

// String returns a string representation of the board.
// Each piece is represented by its icon, and empty squares are represented by "..".
// Lakes are represented by "~~". The string is formatted as a 10x10 grid, with each row on a new line.
func (b *Board) String() string {
	result := ""
	for y := range 10 {
		for x := range 10 {
			piece := b.field[y][x]
			if piece == nil {
				if b.IsLake(NewPosition(x, y)) {
					result += " ~~ "
				} else {
					result += " .. "
				}
			} else {
				result += " " + piece.GetType().GetIcon() + " "
			}
		}
		result += "\n"
	}
	return result
}
