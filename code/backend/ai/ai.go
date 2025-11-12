package ai

import (
	"digital-innovation/stratego/engine"
)

// AI is the interface that all AI implementations must satisfy.
// It extends the PlayerController interface
type AI interface {
	engine.PlayerController
}

type BaseAI struct {
	player *engine.Player
	memory *AIMemory
}

func NewBaseAI(player *engine.Player, hasMemory bool) *BaseAI {
	var memory *AIMemory = nil
	if hasMemory {
		memory = NewAIMemory()
	}
	return &BaseAI{
		player: player,
		memory: memory,
	}
}

// GetPlayer returns the player associated with the AI.
func (ai *BaseAI) GetPlayer() *engine.Player {
	return ai.player
}

// GetControllerType returns the type of the AI controller, which is AIController.
func (ai *BaseAI) GetControllerType() engine.ControllerType {
	return engine.AIController
}

// GetMemory returns the AI's memory system (O(1) position lookup)
func (ai *BaseAI) GetMemory() *AIMemory {
	return ai.memory
}

// AnalyzeMove is called after opponent moves - override in subclasses for learning
// Default implementation updates memory automatically
func (ai *BaseAI) AnalyzeMove(move engine.Move, opponent *engine.Player, round int) {
	if ai.memory == nil {
		return
	}

	from := move.GetFrom()
	to := move.GetTo()

	if ai.memory.Recall(from) != nil {
		ai.memory.MovePiece(from, to)
	}
}

// ObserveCombat is called when combat occurs - override for learning from reveals
// Default implementation updates memory with revealed pieces
func (ai *BaseAI) ObserveCombat(attackerPos, defenderPos engine.Position, attackerPiece, defenderPiece *engine.Piece, round int) {
	if ai.memory == nil {
		return
	}

	ai.memory.UpdateFromCombat(attackerPos, defenderPos, attackerPiece, defenderPiece, round)
}
