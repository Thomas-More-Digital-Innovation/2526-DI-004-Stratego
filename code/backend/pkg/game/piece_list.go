package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
)

var pieceTypes = []models.PieceType{
	models.Flag,
	models.Bomb,
	models.Spy,
	models.Scout,
	models.Miner,
	models.Sergeant,
	models.Lieutenant,
	models.Captain,
	models.Major,
	models.Colonel,
	models.General,
	models.Marshal,
}

// GetPieceList returns a list of 40 pieces for a player according to the game rules
func GetPieceList(player *Player) []*Piece {
	pieceList := make([]*Piece, 0, 40)

	for _, pieceType := range pieceTypes {
		for range pieceType.GetCount() {
			pieceList = append(pieceList, NewPiece(pieceType, player))
		}
	}
	return pieceList
}

// GetPieceListStrategicValue calculates the total strategic value of a piece list
func GetPieceListStrategicValue(pieceList []*Piece) int {
	strategicValue := 0

	for _, piece := range pieceList {
		strategicValue += piece.GetType().GetStrategicValue()
	}
	return strategicValue
}
