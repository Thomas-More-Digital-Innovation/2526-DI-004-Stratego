package engine

type Player struct {
	id             int
	name           string
	pieceScore     int
	avatar         string
	won            bool
	alivePieces    []*Piece
	piecePositions map[*Piece]Position
}

func NewPlayer(id int, name string, avatar string) Player {
	return Player{
		id:             id,
		name:           name,
		pieceScore:     0,
		avatar:         avatar,
		alivePieces:    make([]*Piece, 0, 40),
		piecePositions: make(map[*Piece]Position, 40),
	}
}

func (pl *Player) GetID() int {
	return pl.id
}

func (p *Player) HasWon() bool {
	return p.won
}

func (p *Player) SetWinner() {
	p.won = true
}

func (pl *Player) GetName() string {
	return pl.name
}

func (pl *Player) SetName(name string) {
	pl.name = name
}

func (pl *Player) GetAvatar() string {
	return pl.avatar
}

func (pl *Player) SetAvatar(avatar string) {
	pl.avatar = avatar
}

func (pl *Player) UpdatePieceScore(eliminatedPiece *Piece) {
	pl.pieceScore -= eliminatedPiece.GetStrategicValue()
	pl.RemovePiece(eliminatedPiece)
}

func (pl *Player) GetPieceScore() int {
	return pl.pieceScore
}

func (pl *Player) ResetPieceScore() {
	pl.pieceScore = 0
}

func (pl *Player) InitializePieceScore(initialScore int) {
	pl.pieceScore = initialScore
}

// AddPiece adds a piece to the player's tracking (O(1))
func (pl *Player) AddPiece(piece *Piece, pos Position) {
	pl.alivePieces = append(pl.alivePieces, piece)
	pl.piecePositions[piece] = pos
}

// RemovePiece removes a piece from tracking when eliminated (O(n) but rare)
func (pl *Player) RemovePiece(piece *Piece) {
	for i, p := range pl.alivePieces {
		if p == piece {
			pl.alivePieces = append(pl.alivePieces[:i], pl.alivePieces[i+1:]...)
			break
		}
	}
	delete(pl.piecePositions, piece)
}

// UpdatePiecePosition updates a piece's position (O(1))
func (pl *Player) UpdatePiecePosition(piece *Piece, newPos Position) {
	pl.piecePositions[piece] = newPos
}

// GetPiecePosition returns a piece's position (O(1))
func (pl *Player) GetPiecePosition(piece *Piece) (Position, bool) {
	pos, exists := pl.piecePositions[piece]
	return pos, exists
}

// GetAlivePieces returns all alive pieces (O(1))
func (pl *Player) GetAlivePieces() []*Piece {
	return pl.alivePieces
}
