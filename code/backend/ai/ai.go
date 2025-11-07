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
