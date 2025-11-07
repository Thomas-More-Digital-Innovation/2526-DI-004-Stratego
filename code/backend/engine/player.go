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

// GetID returns the unique identifier of the player.
// This is used by the game engine to keep track of players.
func (pl *Player) GetID() int {
	return pl.id
}

// HasWon returns a boolean indicating whether the player has won the game.
// This is used by the game engine to determine the outcome of the game.
func (p *Player) HasWon() bool {
	return p.won
}

// SetWinner marks the player as the winner of the game.
// This is used by the game engine to update the state of the players after a game is complete.
func (p *Player) SetWinner() {
	p.won = true
}

// GetName returns the name of the player. This is used by the game engine to
// display the names of players in the game.
func (pl *Player) GetName() string {
	return pl.name
}

// SetName updates the name of the player.
// This is used by the game engine to update the player's name.
func (pl *Player) SetName(name string) {
	pl.name = name
}

// GetAvatar returns the avatar of the player. This is used by the game engine to
// display the avatar of the player in the game.
func (pl *Player) GetAvatar() string {
	return pl.avatar
}

// SetAvatar updates the avatar of the player.
// This is used by the game engine to update the player's avatar.
func (pl *Player) SetAvatar(avatar string) {
	pl.avatar = avatar
}

// UpdatePieceScore updates the player's piece score by subtracting the strategic value of the eliminated piece.
// The game engine uses this method to update the state of the players after a battle.
// The method also removes the eliminated piece from the player's set of alive pieces.
func (pl *Player) UpdatePieceScore(eliminatedPiece *Piece) {
	pl.pieceScore -= eliminatedPiece.GetStrategicValue()
	pl.RemovePiece(eliminatedPiece)
}

// GetPieceScore returns the current piece score of the player.
// This is used by the game engine to keep track of the player's score.
// The score is the total value of all pieces left on the board for that player.
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
