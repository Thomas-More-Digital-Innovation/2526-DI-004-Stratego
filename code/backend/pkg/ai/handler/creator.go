// Package AIhandler provides a factory for creating AI instances
package AIhandler

import (
	"digital-innovation/gostrategy/pkg/ai"
	"digital-innovation/gostrategy/pkg/ai/fafo"
	"digital-innovation/gostrategy/pkg/ai/fato"
	"digital-innovation/gostrategy/pkg/game"
	"digital-innovation/gostrategy/pkg/game/models"
	"fmt"
)

// CreateAI is a factory function that returns an AI instance based on the given type
func CreateAI(aiType string, player *game.Player) (ai.AI, error) {
	switch aiType {
	case models.Fafo:
		return fafo.NewAI(player, false), nil
	case models.Fato:
		return fato.NewAI(player, true), nil
	default:
		return nil, fmt.Errorf("i don't know that AI! %s", aiType)
	}
}
