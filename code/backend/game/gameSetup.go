package game

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"errors"
	"fmt"
	"math/rand/v2"
)

// SetupGame initializes the board with pieces for both players and prepares the game
func SetupGame(game *Game, player1Pieces, player2Pieces []*engine.Piece) error {
	// Validate piece counts
	if len(player1Pieces) != 40 || len(player2Pieces) != 40 {
		return errors.New("each player must have exactly 40 pieces")
	}

	// Place player 1 pieces (bottom 4 rows: 6-9)
	if err := placePiecesInRows(game.Board, player1Pieces, 6, 9); err != nil {
		return err
	}

	// Place player 2 pieces (top 4 rows: 0-3)
	if err := placePiecesInRows(game.Board, player2Pieces, 0, 3); err != nil {
		return err
	}

	game.InitializePieces()

	game.Players[0].InitializePieceScore(GetPieceListStrategicValue(player1Pieces))
	game.Players[1].InitializePieceScore(GetPieceListStrategicValue(player2Pieces))

	return nil
}

// placePiecesInRows places pieces sequentially in the specified row range
func placePiecesInRows(board *engine.Board, pieces []*engine.Piece, startRow, endRow int) error {
	pieceIndex := 0
	for y := startRow; y <= endRow; y++ {
		for x := 0; x < 10; x++ {
			if pieceIndex >= len(pieces) {
				return nil
			}
			pos := engine.NewPosition(x, y)
			board.SetPieceAt(pos, pieces[pieceIndex])
			pieceIndex++
		}
	}
	return nil
}

// RandomSetup creates a random valid piece placement for a player
func RandomSetup(player *engine.Player) []*engine.Piece {
	pieces := GetPieceList(player)
	// Shuffle pieces for random placement
	rand.Shuffle(len(pieces), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})
	return pieces
}

// CustomSetup allows manual piece placement with validation
func CustomSetup(player *engine.Player, piecePositions map[engine.Position]*engine.Piece) ([]*engine.Piece, error) {
	if len(piecePositions) != 40 {
		return nil, errors.New("must place exactly 40 pieces")
	}

	pieces := make([]*engine.Piece, 0, 40)
	for _, piece := range piecePositions {
		pieces = append(pieces, piece)
	}

	return pieces, nil
}

// ParseSetup converts 40-byte binary setup data (bitpacked) or 40-character rank strings
// into an ordered piece list for the player.
func ParseSetup(player *engine.Player, data []byte) ([]*engine.Piece, error) {
	if len(data) != 40 {
		return nil, fmt.Errorf("setup data must be 40 bytes, got %d", len(data))
	}

	// Detect format: bitpacked vs rank characters
	isBitpacked := false
	for _, b := range data {
		if b&engine.BitOccupied != 0 {
			isBitpacked = true
			break
		}
	}

	pieces := make([]*engine.Piece, 0, 40)
	for i := 0; i < 40; i++ {
		cell := data[i]
		var pieceType *models.PieceType

		if isBitpacked {
			if cell&engine.BitOccupied == 0 {
				return nil, fmt.Errorf("cell %d is empty, all 40 cells must be occupied", i)
			}
			pieceType = engine.GetPieceTypeFromCell(cell)
		} else {
			// Rank character format
			if cell == '.' || cell == ' ' {
				return nil, fmt.Errorf("cell %d is empty in rank string, all 40 cells must be occupied", i)
			}
			pieceID, ok := engine.GetPieceIDFromRank(cell)
			if !ok {
				return nil, fmt.Errorf("invalid piece rank '%c' at position %d", cell, i)
			}
			pieceType = engine.GetPieceTypeFromID(pieceID)
		}

		if pieceType == nil {
			return nil, fmt.Errorf("unknown piece type at index %d", i)
		}

		pieces = append(pieces, engine.NewPiece(*pieceType, player))
	}

	if len(pieces) != 40 {
		return nil, fmt.Errorf("expected 40 pieces, got %d", len(pieces))
	}

	return pieces, nil
}

// QuickStart creates a game with random setups for both players
func QuickStart(controller1, controller2 engine.PlayerController) *Game {
	game := NewGame(controller1, controller2)

	player1 := controller1.GetPlayer()
	player2 := controller2.GetPlayer()

	player1Pieces := RandomSetup(player1)
	player2Pieces := RandomSetup(player2)

	if err := SetupGame(game, player1Pieces, player2Pieces); err != nil {
		panic("Failed to setup game: " + err.Error())
	}

	return game
}
