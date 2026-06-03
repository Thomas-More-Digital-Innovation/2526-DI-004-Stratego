// Package ai provides base interfaces and common functionality for AI implementations
package ai

import (
	"digital-innovation/gostrategy/internal/game"
)

// AI is the interface that all AI implementations must satisfy.
// It extends the PlayerController interface
type AI interface {
	game.PlayerController
}

// BaseAI provides common functionality for all AI types
type BaseAI struct {
	player *game.Player
	memory *Memory
}

// NewBaseAI creates a new BaseAI instance
func NewBaseAI(player *game.Player, hasMemory bool) *BaseAI {
	var memory *Memory
	if hasMemory {
		memory = NewMemory()
	}
	return &BaseAI{
		player: player,
		memory: memory,
	}
}

// GetPlayer returns the player associated with the AI.
func (ai *BaseAI) GetPlayer() *game.Player {
	return ai.player
}

// GetControllerType returns the type of the AI controller, which is AIController.
func (ai *BaseAI) GetControllerType() game.ControllerType {
	return game.AIController
}

// GetMemory returns the AI's memory system (O(1) position lookup)
func (ai *BaseAI) GetMemory() *Memory {
	return ai.memory
}

// AnalyzeMove is called after opponent moves - override in subclasses for learning
// Default implementation updates memory automatically
func (ai *BaseAI) AnalyzeMove(move game.Move, _ *game.Player, _ int) {
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
func (ai *BaseAI) ObserveCombat(attackerPos, defenderPos game.Position, attackerPiece, defenderPiece *game.Piece, round int) {
	if ai.memory == nil {
		return
	}

	ai.memory.UpdateFromCombat(attackerPos, defenderPos, attackerPiece, defenderPiece, round)
}
