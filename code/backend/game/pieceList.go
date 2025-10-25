package game

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
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

func GetPieceList(player *engine.Player) []*engine.Piece {
	pieceList := make([]*engine.Piece, 0, 40) // Pre-allocate for 40 pieces

	for _, pieceType := range pieceTypes {
		for range pieceType.GetCount() {
			pieceList = append(pieceList, engine.NewPiece(pieceType, player))
		}
	}
	return pieceList
}

func GetPieceListStrategicValue(pieceList []*engine.Piece) int {
	strategicValue := 0

	for _, piece := range pieceList {
		strategicValue += piece.GetType().GetStrategicValue()
	}
	return strategicValue
}
