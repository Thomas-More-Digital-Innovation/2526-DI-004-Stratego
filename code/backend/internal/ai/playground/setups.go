package playground

import (
	"digital-innovation/gostrategy/internal/game"
	"fmt"
	"strings"
)

// BoardSetup represents a piece configuration layout
type BoardSetup struct {
	Name        string
	Description string
	Data        string
}

// PredefinedSetups holds default setups in memory
var PredefinedSetups = []BoardSetup{
	{
		Name:        "Honey Pot",
		Description: "A defensive layout shielding the Flag with bombs and high ranks.",
		Data:        "220M22B5BB971877364422342258B56663533B4B",
	},
}

// NormalizeSetup cleans character representations and validates format/length
func NormalizeSetup(data string) (string, error) {
	data = strings.TrimSpace(data)
	data = strings.ReplaceAll(data, "F", "0")
	if len(data) != 40 {
		return "", fmt.Errorf("setup data must be exactly 40 characters, got %d", len(data))
	}
	return data, nil
}

// ValidateSetupString checks if the setup contains valid rank characters in correct counts
func ValidateSetupString(player *game.Player, data string) error {
	normalized, err := NormalizeSetup(data)
	if err != nil {
		return err
	}
	_, err = game.ParseSetup(player, []byte(normalized))
	return err
}
