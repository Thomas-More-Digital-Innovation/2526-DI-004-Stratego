// Package AIhandler provides a factory for creating AI instances
package AIhandler

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/ai/fafo"
	"digital-innovation/gostrategy/internal/ai/fato"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"fmt"
)

// CreateAI is a factory function that returns an AI instance based on the given type
func CreateAI(aiType string, player *game.Player) (ai.AI, error) {
	return CreateAIWithOptions(aiType, player, nil)
}

// CreateAIWithOptions creates an AI instance with customizable configuration options
func CreateAIWithOptions(aiType string, player *game.Player, options map[string]any) (ai.AI, error) {
	switch aiType {
	case models.Fafo:
		return fafo.NewAI(player, false), nil
	case models.Fato:
		aggression := 0.5
		if options != nil {
			if val, ok := options["aggression"]; ok {
				if f, ok := val.(float64); ok {
					aggression = f
				} else if i, ok := val.(int); ok {
					aggression = float64(i)
				}
			}
		}
		return fato.NewAIWithAggression(player, true, aggression), nil
	default:
		return nil, fmt.Errorf("i don't know that AI! %s", aiType)
	}
}
