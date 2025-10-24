package engine

import (
	"digital-innovation/stratego/models"
)

type Piece struct {
	pieceType models.PieceType
	player    *Player
	alive     bool
	revealed  bool
}

// NewPiece creates a new Piece with the given pieceType and player.
// The piece is initially alive and not revealed to the opponent.
// The game engine uses this method to create new pieces on the board.
func NewPiece(pieceType models.PieceType, player *Player) *Piece {
	return &Piece{
		pieceType: pieceType,
		player:    player,
		alive:     true,
		revealed:  false,
	}
}

// GetOwner returns the owner of the piece.
// The owner is the player that controls the piece.
// The game engine uses this method to determine which player owns a piece.
func (p *Piece) GetOwner() *Player {
	return p.player
}

// GetType returns the type of the piece.
func (p *Piece) GetType() *models.PieceType {
	return &p.pieceType
}

// GetValue returns the value of the piece. The value is used by the game engine to
// determine the strength of a piece in battle. The value is also used to calculate
// the score of a player in the game. The score is the total value of all pieces
// captured by the player. The player with the highest score at the end of the game
// wins.
func (p *Piece) GetStrategicValue() int {
	return p.pieceType.GetStrategicValue()
}

// IsRevealed returns a boolean indicating whether the piece has been revealed to the opponent.
// A piece is revealed when it is attacked by an opponent's piece.
// The game engine uses this method to determine whether a piece is visible to the opponent.
func (p *Piece) IsRevealed() bool {
	return p.revealed
}

// GetRank returns the rank of the piece.
// The rank is a byte value that indicates the piece's strength in battle.
// The rank is used by the game engine to determine the outcome of a battle.
// A piece with a higher rank will always win against a piece with a lower rank.
// If the ranks are equal, both pieces are eliminated from the game.
func (p *Piece) GetRank() byte {
	return p.pieceType.GetRank()
}

// CanMove returns a boolean indicating whether the piece can move on the board.
// A piece can move if its PieceType is movable. The game engine uses this method
// to determine whether a piece is valid to move on the board.
func (p *Piece) CanMove() bool {
	return p.pieceType.IsMovable()
}

// IsAlive returns a boolean indicating whether the piece is still alive in the game.
// A piece is eliminated from the game when it is attacked by a piece of higher rank
// or when it is attacked by a piece of equal rank. The game engine uses this method
// to determine whether a piece is still valid on the board.
func (p *Piece) IsAlive() bool {
	return p.alive
}

// Reveal sets the Revealed field of the piece to true, indicating that
// the piece has been revealed to the opponent. This is used by the
// game engine to update the state of the pieces after a battle.
func (p *Piece) Reveal() {
	p.revealed = true
}

// Eliminate marks the piece as eliminated from the game.
// It is used by the game engine to update the state of the pieces after a battle.
func (p *Piece) Eliminate() {
	p.alive = false
	player := p.GetOwner()
	player.UpdatePieceScore(p)
}

// Attack is called when a piece is attacked by another piece.
// It determines the outcome of the battle based on the rank of the two pieces.
// If the target is the enemy's flag, the attacking piece wins and the game is over.
// If the attacking piece is a spy and the target is a marshal, the spy wins.
// If the target is a bomb, the attacking piece is eliminated unless it is a miner.
// If neither of the above conditions are met, the outcome of the battle is determined by the standard rules:
// the piece with the higher rank wins, and if the ranks are equal, both pieces are eliminated.
//
// returns an array of two pieces: the attacking piece and the target piece
func (p *Piece) Attack(target *Piece) [2]*Piece {
	switch {
	case target.GetRank() == models.Flag.GetRank():
		p.resolveFlagCapture()
	case p.GetRank() == models.Spy.GetRank() && target.GetRank() == models.Marshal.GetRank():
		p.resolveSpyAttackingMarshal(target)
	case target.GetRank() == models.Bomb.GetRank():
		p.resolveBombAttack(target)
	default:
		p.resolveStandardAttack(target)
	}
	return [2]*Piece{p, target}
}

func (p *Piece) resolveFlagCapture() {
	// implement win logic
}

func (p *Piece) resolveSpyAttackingMarshal(target *Piece) {
	target.Eliminate()
}

// in this case, target is the bomb
func (p *Piece) resolveBombAttack(target *Piece) {
	if p.GetRank() == models.Miner.GetRank() {
		target.Eliminate()
	} else {
		p.Eliminate()
	}
}

func (p *Piece) resolveStandardAttack(target *Piece) {
	switch {
	case p.GetRank() > target.GetRank():
		target.Eliminate()
	case p.GetRank() < target.GetRank():
		p.Eliminate()
	default:
		p.Eliminate()
		target.Eliminate()
	}
}
