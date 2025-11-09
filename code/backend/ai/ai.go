package ai

import "digital-innovation/stratego/engine"

// AI is the interface that all AI implementations must satisfy.
// It extends the PlayerController interface, meaning all AIs must implement
// GetPlayer(), GetControllerType(), and MakeMove(board).
//
// Any AI type (FafoAI, FatoAI, etc.) that implements these methods
// automatically satisfies this interface.
type AI interface {
	engine.PlayerController
}

type BaseAI struct {
	player *engine.Player
}

func NewBaseAI(player *engine.Player) *BaseAI {
	return &BaseAI{player: player}
}

func (ai *BaseAI) GetPlayer() *engine.Player {
	return ai.player
}

func (ai *BaseAI) GetControllerType() engine.ControllerType {
	return engine.AIController
}
