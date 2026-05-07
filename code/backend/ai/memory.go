package ai

import (
	"digital-innovation/stratego/engine"
	"sync"
)

// MemoryEntry represents a single piece known to the AI memory
type MemoryEntry struct {
	Piece      *engine.Piece
	Confidence float64 // how sure are we? 1.0 = revealed, <1.0 = guess
	LastSeen   int
}

// Memory represents the AI's internal knowledge of the board
type Memory struct {
	field [10][10]*MemoryEntry
	mutex sync.RWMutex
}

// NewMemory creates a new Memory instance
func NewMemory() *Memory {
	return &Memory{}
}

// Remember updates memory with new information about a piece
func (m *Memory) Remember(pos engine.Position, piece *engine.Piece, confidence float64, round int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.field[pos.Y][pos.X] = &MemoryEntry{
		Piece:      piece,
		Confidence: confidence,
		LastSeen:   round,
	}
}

// Recall retrieves what we know about a position (O(1))
func (m *Memory) Recall(pos engine.Position) *MemoryEntry {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.field[pos.Y][pos.X]
}

// Forget clears memory at a position
func (m *Memory) Forget(pos engine.Position) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.field[pos.Y][pos.X] = nil
}

// MovePiece updates memory when a piece moves (critical for correctness)
func (m *Memory) MovePiece(from, to engine.Position) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.field[from.Y][from.X] != nil {
		m.field[to.Y][to.X] = m.field[from.Y][from.X]
		m.field[from.Y][from.X] = nil
	}
}

// UpdateFromCombat processes combat results to update memory
// Call this when pieces are revealed in combat
func (m *Memory) UpdateFromCombat(attackerPos, defenderPos engine.Position, attackerPiece, defenderPiece *engine.Piece, round int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// attacker survived and moved, update memory
	if attackerPiece != nil && attackerPiece.IsAlive() {
		m.field[defenderPos.Y][defenderPos.X] = &MemoryEntry{
			Piece:      attackerPiece,
			Confidence: 1.0,
			LastSeen:   round,
		}
		m.field[attackerPos.Y][attackerPos.X] = nil
	} else {
		// Attacker died, clear both positions
		m.field[attackerPos.Y][attackerPos.X] = nil

		// defender survived, remember it
		if defenderPiece != nil && defenderPiece.IsAlive() {
			m.field[defenderPos.Y][defenderPos.X] = &MemoryEntry{
				Piece:      defenderPiece,
				Confidence: 1.0,
				LastSeen:   round,
			}
		} else {
			// both died
			m.field[defenderPos.Y][defenderPos.X] = nil
		}
	}
}

// Clear resets all memory (for new game)
func (m *Memory) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.field = [10][10]*MemoryEntry{}
}

// GetKnownEnemyPositions returns all positions where we remember enemy pieces
// Useful for targeting decisions
func (m *Memory) GetKnownEnemyPositions() []engine.Position {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	positions := make([]engine.Position, 0, 10)
	for y := range 10 {
		for x := 0; x < 10; x++ {
			if m.field[y][x] != nil {
				positions = append(positions, engine.NewPosition(x, y))
			}
		}
	}
	return positions
}

// DecayConfidence reduces confidence over time for guesses
// Call periodically to forget old, uncertain information
func (m *Memory) DecayConfidence(decayRate float64, minConfidence float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for y := range 10 {
		for x := range 10 {
			if entry := m.field[y][x]; entry != nil && entry.Confidence < 1.0 {
				entry.Confidence *= (1.0 - decayRate)
				if entry.Confidence < minConfidence {
					m.field[y][x] = nil // Forget low-confidence guesses
				}
			}
		}
	}
}
