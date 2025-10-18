package engine

type Player struct {
	id         int
	name       string
	pieceScore int
	avatar     string
}

func NewPlayer(id int, name string, avatar string) Player {
	return Player{
		id:         id,
		name:       name,
		pieceScore: 0,
		avatar:     avatar,
	}
}

func (pl *Player) GetID() int {
	return pl.id
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
