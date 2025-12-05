package engine

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Compact board-setup encoding for player setups (40 cells = 4 rows × 10 cols):
// - map each piece-icon char to a 4-bit id (0..15)
// - pack two cells per byte: (hi<<4)|lo

const (
	PlayerSetupCells = 40
	PlayerSetupRows  = 4
	PlayerSetupCols  = 10
)

var (
	// char -> 4-bit id. Keep this stable across versions.
	charToID = map[byte]byte{
		'0': 0,  // Flag
		'B': 1,  // Bomb
		'1': 2,  // Spy
		'2': 3,  // Scout
		'3': 4,  // Miner
		'4': 5,  // Sergeant
		'5': 6,  // Lieutenant
		'6': 7,  // Captain
		'7': 8,  // Major
		'8': 9,  // Colonel
		'9': 10, // General
		'M': 11, // Marshal
		'.': 15, // empty (optional)
	}
	// reverse lookup for unpacking
	idToChar = func() [16]byte {
		var a [16]byte
		for k, v := range charToID {
			a[v] = k
		}
		// ensure unused ids are set to '?' for visibility
		for i := 0; i < 16; i++ {
			if a[i] == 0 {
				a[i] = '?'
			}
		}
		return a
	}()
)

// EncodeBoardSetup accepts either:
// - rows: a slice of 4 strings each 10 chars long, or
// - a single string of length 40 (concatenated rows)
// It returns a compact base64 string representing the packed 4-bit ids.
func EncodeBoardSetup(rows []string) (string, error) {
	if len(rows) == 0 {
		return "", errors.New("no rows provided")
	}
	// if a single long string provided, split it into rows of 10
	var conc string
	if len(rows) == 1 && len(rows[0]) == PlayerSetupCells {
		conc = rows[0]
	} else {
		if len(rows) != PlayerSetupRows {
			return "", fmt.Errorf("expected %d rows, got %d", PlayerSetupRows, len(rows))
		}
		for _, r := range rows {
			if len(r) != PlayerSetupCols {
				return "", fmt.Errorf("each row must be %d chars, got %d", PlayerSetupCols, len(r))
			}
			conc += r
		}
	}
	if len(conc) != PlayerSetupCells {
		return "", fmt.Errorf("setup length must be %d chars, got %d", PlayerSetupCells, len(conc))
	}

	packed := make([]byte, PlayerSetupCells/2) // 40/2 = 20 bytes
	for i := 0; i < PlayerSetupCells; i += 2 {
		hiCh := conc[i]
		loCh := conc[i+1]
		hi, ok := charToID[hiCh]
		if !ok {
			return "", fmt.Errorf("unknown piece char: %c", hiCh)
		}
		lo, ok := charToID[loCh]
		if !ok {
			return "", fmt.Errorf("unknown piece char: %c", loCh)
		}
		packed[i/2] = (hi << 4) | (lo & 0x0F)
	}
	return base64.StdEncoding.EncodeToString(packed), nil
}

func DecodeBoardSetup(encoded string) ([]string, error) {
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	if len(b) != PlayerSetupCells/2 {
		return nil, fmt.Errorf("decoded length must be %d bytes, got %d", PlayerSetupCells/2, len(b))
	}
	cells := make([]byte, PlayerSetupCells)
	for i := 0; i < len(b); i++ {
		hi := (b[i] >> 4) & 0x0F
		lo := b[i] & 0x0F
		cells[i*2] = idToChar[hi]
		cells[i*2+1] = idToChar[lo]
	}
	// split into rows
	rows := make([]string, PlayerSetupRows)
	for r := 0; r < PlayerSetupRows; r++ {
		rows[r] = string(cells[r*PlayerSetupCols : (r+1)*PlayerSetupCols])
	}
	return rows, nil
}

// ParseBoardSetupSmart accepts a raw setup string from storage. It detects and
// accepts either the current human-friendly layout (JSON array of rows or a
// concatenated 40-char string) and returns
// the human-readable rows. Use this when loading setups from DB/files.
func ParseBoardSetupSmart(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if strings.Contains(raw, "\n") || strings.Contains(raw, ",") || len(raw) == PlayerSetupCells {
		sep := '\n'
		if strings.Contains(raw, ",") {
			sep = ','
		}
		parts := splitKeepNonEmpty(raw, sep)
		if len(parts) == PlayerSetupRows {
			for _, r := range parts {
				if len(r) != PlayerSetupCols {
					return nil, fmt.Errorf("expected row length %d", PlayerSetupCols)
				}
			}
			return parts, nil
		}
		if len(parts) == 1 && len(parts[0]) == PlayerSetupCells {
			out := make([]string, PlayerSetupRows)
			for i := range PlayerSetupRows {
				out[i] = parts[0][i*PlayerSetupCols : (i+1)*PlayerSetupCols]
			}
			return out, nil
		}
	}
	return nil, errors.New("unrecognized board setup format")
}

func splitKeepNonEmpty(raw string, sep rune) []string {
	var parts []string
	for p := range strings.SplitSeq(raw, string(sep)) {
		p = strings.TrimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}
