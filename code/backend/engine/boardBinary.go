package engine

import (
	"digital-innovation/stratego/models"
	"encoding/base64"
	"fmt"
)

// Binary format: 8 bits per square
// Bit 7: Occupied | Bits 6-2: Piece type | Bit 1: Color | Bit 0: Moved
const (
	BitOccupied    = 1 << 7
	BitColor       = 1 << 1
	BitMoved       = 1 << 0
	MaskPieceType  = 0x7C
	ShiftPieceType = 2
)

const (
	PieceIDFlag       byte = 1
	PieceIDBomb       byte = 2
	PieceIDSpy        byte = 3
	PieceIDScout      byte = 4
	PieceIDMiner      byte = 5
	PieceIDSergeant   byte = 6
	PieceIDLieutenant byte = 7
	PieceIDCaptain    byte = 8
	PieceIDMajor      byte = 9
	PieceIDColonel    byte = 10
	PieceIDGeneral    byte = 11
	PieceIDMarshal    byte = 12
)

var rankToPieceID = map[byte]byte{
	'0': PieceIDFlag, 'B': PieceIDBomb, '1': PieceIDSpy, '2': PieceIDScout,
	'3': PieceIDMiner, '4': PieceIDSergeant, '5': PieceIDLieutenant, '6': PieceIDCaptain,
	'7': PieceIDMajor, '8': PieceIDColonel, '9': PieceIDGeneral, 'M': PieceIDMarshal,
}

var pieceIDToRank = func() [32]byte {
	var a [32]byte
	for k, v := range rankToPieceID {
		a[v] = k
	}
	return a
}()

var pieceTypeToID = map[*models.PieceType]byte{
	&models.Flag: PieceIDFlag, &models.Bomb: PieceIDBomb, &models.Spy: PieceIDSpy,
	&models.Scout: PieceIDScout, &models.Miner: PieceIDMiner, &models.Sergeant: PieceIDSergeant,
	&models.Lieutenant: PieceIDLieutenant, &models.Captain: PieceIDCaptain, &models.Major: PieceIDMajor,
	&models.Colonel: PieceIDColonel, &models.General: PieceIDGeneral, &models.Marshal: PieceIDMarshal,
}

var idToPieceType = map[byte]*models.PieceType{
	PieceIDFlag: &models.Flag, PieceIDBomb: &models.Bomb, PieceIDSpy: &models.Spy,
	PieceIDScout: &models.Scout, PieceIDMiner: &models.Miner, PieceIDSergeant: &models.Sergeant,
	PieceIDLieutenant: &models.Lieutenant, PieceIDCaptain: &models.Captain, PieceIDMajor: &models.Major,
	PieceIDColonel: &models.Colonel, PieceIDGeneral: &models.General, PieceIDMarshal: &models.Marshal,
}

// Encode full 10x10 board to 100 bytes
func EncodeBoard(board *Board, movedPositions map[Position]bool) []byte {
	data := make([]byte, 100)
	field := board.GetField()

	for y := range 10 {
		for x := range 10 {
			piece := field[y][x]
			if piece == nil {
				continue
			}

			cell := byte(BitOccupied)

			// Get piece ID from rank character
			rank := piece.GetType().GetRank()
			pieceID := rankToPieceID[rank]
			cell |= (pieceID << ShiftPieceType) & MaskPieceType

			if piece.GetOwner().GetID() == 2 {
				cell |= BitColor
			}

			if movedPositions != nil && movedPositions[NewPosition(x, y)] {
				cell |= BitMoved
			}

			data[y*10+x] = cell
		}
	}

	return data
}

// Decode 100 bytes to full 10x10 board
func DecodeBoard(data []byte, player1, player2 *Player) (*Board, map[Position]bool) {
	board := NewBoard()
	movedPositions := make(map[Position]bool)

	for y := range 10 {
		for x := range 10 {
			cell := data[y*10+x]
			if cell&BitOccupied == 0 {
				continue
			}

			pieceID := (cell & MaskPieceType) >> ShiftPieceType
			pieceType := idToPieceType[pieceID]

			owner := player1
			if cell&BitColor != 0 {
				owner = player2
			}

			board.SetPieceAt(NewPosition(x, y), NewPiece(*pieceType, owner))

			if cell&BitMoved != 0 {
				movedPositions[NewPosition(x, y)] = true
			}
		}
	}

	return board, movedPositions
}

// Encode to base64 string
func EncodeBoardToBase64(board *Board, movedPositions map[Position]bool) string {
	return base64.StdEncoding.EncodeToString(EncodeBoard(board, movedPositions))
}

// Decode from base64 string
func DecodeBoardFromBase64(encoded string, player1, player2 *Player) (*Board, map[Position]bool, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, nil, err
	}
	if len(data) != 100 {
		return nil, nil, fmt.Errorf("invalid data length: expected 100, got %d", len(data))
	}
	board, moved := DecodeBoard(data, player1, player2)
	return board, moved, nil
}

// Encode 4x10 setup to 40 bytes
func EncodeSetup(rows []string, playerID int) ([]byte, error) {
	if len(rows) != 4 {
		return nil, fmt.Errorf("expected 4 rows, got %d", len(rows))
	}

	data := make([]byte, 40)
	for i, row := range rows {
		if len(row) != 10 {
			return nil, fmt.Errorf("row %d must be 10 chars, got %d", i, len(row))
		}

		for j, char := range row {
			if char == '.' || char == ' ' {
				continue
			}

			cell := byte(BitOccupied)
			cell |= (rankToPieceID[byte(char)] << ShiftPieceType) & MaskPieceType

			if playerID == 2 {
				cell |= BitColor
			}

			data[i*10+j] = cell
		}
	}

	return data, nil
}

// Decode 40 bytes to 4x10 setup
func DecodeSetup(data []byte) ([]string, int, error) {
	if len(data) != 40 {
		return nil, 0, fmt.Errorf("invalid data length: expected 40, got %d", len(data))
	}

	rows := make([]string, 4)
	playerID := 1

	for i := 0; i < 4; i++ {
		row := make([]byte, 10)
		for j := 0; j < 10; j++ {
			cell := data[i*10+j]
			if cell&BitOccupied == 0 {
				row[j] = '.'
			} else {
				pieceID := (cell & MaskPieceType) >> ShiftPieceType
				row[j] = pieceIDToRank[pieceID]

				if cell&BitColor != 0 {
					playerID = 2
				}
			}
		}
		rows[i] = string(row)
	}

	return rows, playerID, nil
}

// Validate setup has correct piece counts
func ValidateSetup(rows []string) error {
	if len(rows) != 4 {
		return fmt.Errorf("setup must have 4 rows, got %d", len(rows))
	}

	counts := make(map[byte]int)
	for _, row := range rows {
		if len(row) != 10 {
			return fmt.Errorf("each row must be 10 chars")
		}
		for _, c := range row {
			if c != '.' && c != ' ' {
				counts[byte(c)]++
			}
		}
	}

	expected := map[byte]int{
		'0': 1, 'B': 6, '1': 1, '2': 8, '3': 5, '4': 4,
		'5': 4, '6': 4, '7': 3, '8': 2, '9': 1, 'M': 1,
	}

	for piece, count := range expected {
		if counts[piece] != count {
			return fmt.Errorf("piece %c: expected %d, got %d", piece, count, counts[piece])
		}
	}

	return nil
}
