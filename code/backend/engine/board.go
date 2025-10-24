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

func (b *Board) SetPieceAt(pos Position, piece *Piece) {
	b.field[pos.Y][pos.X] = piece
}

func (b *Board) GetField() [10][10]*Piece {
	return b.field
}

func (b *Board) GetPieceAt(pos Position) *Piece {
	return b.field[pos.Y][pos.X]
}

func (b *Board) IsLake(pos Position) bool {
	for _, lakePos := range b.lakes {
		if lakePos == pos {
			return true
		}
	}
	return false
}

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

func (b *Board) RandomizeSetup(player *Player, pieces []Piece) error {
	panic("unimplemented")
}

func (b *Board) SwapPieces(pos1 Position, pos2 Position) error {
	piece1 := b.GetPieceAt(pos1)
	piece2 := b.GetPieceAt(pos2)

	if piece1 == nil || piece2 == nil {
		return errors.New("one or both positions are empty")
	}

	b.field[pos1.Y][pos1.X], b.field[pos2.Y][pos2.X] = piece2, piece1
	return nil
}

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
