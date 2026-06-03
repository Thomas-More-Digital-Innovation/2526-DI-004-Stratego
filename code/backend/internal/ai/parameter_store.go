package ai

import (
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Parameters defines configuration and weights stored for an AI type.
type Parameters struct {
	AIType     string             `json:"ai_type"`
	Name       string             `json:"name"`
	Aggression float64            `json:"aggression"`
	Weights    map[string]float64 `json:"weights"`
	Config     map[string]any     `json:"config"`
}

var (
	fileMutex    sync.Mutex
	fallbackFile = "ai_parameters.json"
)

// SetFallbackFile overrides the default JSON fallback persistence path.
func SetFallbackFile(path string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	fallbackFile = path
}

// GetDefault retrieves standard hardcoded settings for the specified AI type.
func GetDefault(aiType string, name string) *Parameters {
	defaults := &Parameters{
		AIType:     aiType,
		Name:       name,
		Aggression: 0.5,
		Weights:    make(map[string]float64),
		Config:     make(map[string]any),
	}

	switch aiType {
	case "fato":
		defaults.Aggression = 0.5
	case "heuristic":
		defaults.Aggression = 0.5
		defaults.Weights["Flag"] = 10000.0
		defaults.Weights["Bomb"] = 0.0
		defaults.Weights["Spy"] = 90.0
		defaults.Weights["Scout"] = 20.0
		defaults.Weights["Miner"] = 50.0
		defaults.Weights["Sergeant"] = 30.0
		defaults.Weights["Lieutenant"] = 40.0
		defaults.Weights["Captain"] = 50.0
		defaults.Weights["Major"] = 60.0
		defaults.Weights["Colonel"] = 70.0
		defaults.Weights["General"] = 80.0
		defaults.Weights["Marshal"] = 100.0
		defaults.Weights["explorationWeight"] = 5.0
		defaults.Weights["combatWeight"] = 10.0
	case "minimax":
		defaults.Aggression = 0.5
		defaults.Config["depth"] = 2.0
		defaults.Weights["Flag"] = 10000.0
		defaults.Weights["Bomb"] = 0.0
		defaults.Weights["Spy"] = 90.0
		defaults.Weights["Scout"] = 20.0
		defaults.Weights["Miner"] = 50.0
		defaults.Weights["Sergeant"] = 30.0
		defaults.Weights["Lieutenant"] = 40.0
		defaults.Weights["Captain"] = 50.0
		defaults.Weights["Major"] = 60.0
		defaults.Weights["Colonel"] = 70.0
		defaults.Weights["General"] = 80.0
		defaults.Weights["Marshal"] = 100.0
		defaults.Weights["explorationWeight"] = 5.0
		defaults.Weights["combatWeight"] = 10.0
	case "mcts":
		defaults.Aggression = 0.5
		defaults.Config["iterations"] = 100.0
		defaults.Config["exploration_constant"] = 1.414
	}

	return defaults
}

// Load loads AI parameters from DB or fallback file.
func Load(aiType string, name string) (*Parameters, error) {
	if db.DB != nil {
		var model models.AIParameter
		err := db.DB.Where("ai_type = ? AND name = ?", aiType, name).First(&model).Error
		if err == nil {
			var weights map[string]float64
			var config map[string]any
			_ = json.Unmarshal([]byte(model.Weights), &weights)
			_ = json.Unmarshal([]byte(model.Config), &config)
			return &Parameters{
				AIType:     model.AIType,
				Name:       model.Name,
				Aggression: model.Aggression,
				Weights:    weights,
				Config:     config,
			}, nil
		}
	}

	fileMutex.Lock()
	defer fileMutex.Unlock()

	data, err := os.ReadFile(fallbackFile)
	if err == nil {
		var list []Parameters
		if json.Unmarshal(data, &list) == nil {
			for _, p := range list {
				if p.AIType == aiType && p.Name == name {
					return &p, nil
				}
			}
		}
	}

	return GetDefault(aiType, name), nil
}

// Save saves parameters to database and updates JSON backup.
func Save(params *Parameters) error {
	if params == nil {
		return errors.New("cannot save nil parameters")
	}

	weightsJSON, _ := json.Marshal(params.Weights)
	configJSON, _ := json.Marshal(params.Config)

	if db.DB != nil {
		var model models.AIParameter
		err := db.DB.Where("ai_type = ? AND name = ?", params.AIType, params.Name).First(&model).Error
		if err == nil {
			model.Aggression = params.Aggression
			model.Weights = string(weightsJSON)
			model.Config = string(configJSON)
			model.UpdatedAt = time.Now()
			db.DB.Save(&model)
		} else {
			newModel := models.AIParameter{
				AIType:     params.AIType,
				Name:       params.Name,
				Aggression: params.Aggression,
				Weights:    string(weightsJSON),
				Config:     string(configJSON),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			db.DB.Create(&newModel)
		}
	}

	fileMutex.Lock()
	defer fileMutex.Unlock()

	var list []Parameters
	data, err := os.ReadFile(fallbackFile)
	if err == nil {
		_ = json.Unmarshal(data, &list)
	}

	found := false
	for i, p := range list {
		if p.AIType == params.AIType && p.Name == params.Name {
			list[i] = *params
			found = true
			break
		}
	}
	if !found {
		list = append(list, *params)
	}

	newData, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fallbackFile, newData, 0600)
}
