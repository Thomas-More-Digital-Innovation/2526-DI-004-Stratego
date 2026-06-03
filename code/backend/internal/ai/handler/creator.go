// Package AIhandler provides a factory for creating AI instances.
package AIhandler

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/ai/fafo"
	"digital-innovation/gostrategy/internal/ai/fato"
	"digital-innovation/gostrategy/internal/ai/heuristic"
	"digital-innovation/gostrategy/internal/ai/mcts"
	"digital-innovation/gostrategy/internal/ai/minimax"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"fmt"
)

const defaultName = "default"

// CreateAI returns an AI instance based on the given type with default options.
func CreateAI(aiType string, player *game.Player) (ai.AI, error) {
	return CreateAIWithOptions(aiType, player, nil)
}

// CreateAIWithOptions creates an AI instance with customizable configuration options.
func CreateAIWithOptions(aiType string, player *game.Player, options map[string]any) (ai.AI, error) {
	switch aiType {
	case models.Fafo:
		return fafo.NewAI(player, false), nil
	case models.Fato:
		aggression := 0.5
		if options != nil {
			if val, ok := options["aggression"].(float64); ok {
				aggression = val
			}
		}
		return fato.NewAIWithAggression(player, true, aggression), nil
	case models.Heuristic:
		name := defaultName
		if options != nil {
			if val, ok := options["name"].(string); ok {
				name = val
			}
		}
		params, _ := ai.Load(models.Heuristic, name)
		if options != nil {
			if val, ok := options["aggression"].(float64); ok {
				params.Aggression = val
			}
		}
		return heuristic.NewAIWithParams(player, true, params), nil
	case models.Minimax:
		name := defaultName
		if options != nil {
			if val, ok := options["name"].(string); ok {
				name = val
			}
		}
		params, _ := ai.Load(models.Minimax, name)
		if options != nil {
			if val, ok := options["aggression"].(float64); ok {
				params.Aggression = val
			}
			if val, ok := options["depth"].(float64); ok {
				params.Config["depth"] = val
			}
		}
		return minimax.NewAIWithParams(player, true, params), nil
	case models.Mcts:
		name := defaultName
		if options != nil {
			if val, ok := options["name"].(string); ok {
				name = val
			}
		}
		params, _ := ai.Load(models.Mcts, name)
		if options != nil {
			if val, ok := options["aggression"].(float64); ok {
				params.Aggression = val
			}
			if val, ok := options["iterations"].(float64); ok {
				params.Config["iterations"] = val
			}
		}
		return mcts.NewAIWithParams(player, true, params), nil
	default:
		return nil, fmt.Errorf("i don't know that AI! %s", aiType)
	}
}
