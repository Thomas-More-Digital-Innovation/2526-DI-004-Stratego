package ai_test

import (
	"digital-innovation/gostrategy/ai"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory_CoreFlow(t *testing.T) {
	mem := ai.NewMemory()
	pos := engine.NewPosition(2, 3)
	player := engine.NewPlayer(1, "enemy", "blue")
	piece := engine.NewPiece(models.Marshal, &player)

	// Test Remember and Recall
	mem.Remember(pos, piece, 0.8, 5)
	entry := mem.Recall(pos)
	assert.NotNil(t, entry)
	assert.Equal(t, piece, entry.Piece)
	assert.Equal(t, 0.8, entry.Confidence)
	assert.Equal(t, 5, entry.LastSeen)

	// Test Forget
	mem.Forget(pos)
	entry = mem.Recall(pos)
	assert.Nil(t, entry)
}

func TestMemory_MovePiece(t *testing.T) {
	mem := ai.NewMemory()
	from := engine.NewPosition(1, 1)
	to := engine.NewPosition(1, 2)
	player := engine.NewPlayer(1, "enemy", "blue")
	piece := engine.NewPiece(models.General, &player)

	mem.Remember(from, piece, 1.0, 1)
	mem.MovePiece(from, to)

	assert.Nil(t, mem.Recall(from))
	entry := mem.Recall(to)
	assert.NotNil(t, entry)
	assert.Equal(t, piece, entry.Piece)
}

func TestMemory_UpdateFromCombat(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")

	t.Run("attacker survives", func(t *testing.T) {
		mem := ai.NewMemory()
		attackerPos := engine.NewPosition(0, 0)
		defenderPos := engine.NewPosition(0, 1)
		attacker := engine.NewPiece(models.Marshal, &player1)
		defender := engine.NewPiece(models.General, &player2)
		defender.Eliminate()

		mem.UpdateFromCombat(attackerPos, defenderPos, attacker, defender, 10)
		assert.Nil(t, mem.Recall(attackerPos))
		entry := mem.Recall(defenderPos)
		assert.NotNil(t, entry)
		assert.Equal(t, attacker, entry.Piece)
		assert.Equal(t, 1.0, entry.Confidence)
	})

	t.Run("attacker dies, defender survives", func(t *testing.T) {
		mem := ai.NewMemory()
		attackerPos := engine.NewPosition(0, 0)
		defenderPos := engine.NewPosition(0, 1)
		attacker := engine.NewPiece(models.General, &player1)
		defender := engine.NewPiece(models.Marshal, &player2)
		attacker.Eliminate()

		mem.UpdateFromCombat(attackerPos, defenderPos, attacker, defender, 10)
		assert.Nil(t, mem.Recall(attackerPos))
		entry := mem.Recall(defenderPos)
		assert.NotNil(t, entry)
		assert.Equal(t, defender, entry.Piece)
	})

	t.Run("both die", func(t *testing.T) {
		mem := ai.NewMemory()
		attackerPos := engine.NewPosition(0, 0)
		defenderPos := engine.NewPosition(0, 1)
		attacker := engine.NewPiece(models.Miner, &player1)
		defender := engine.NewPiece(models.Miner, &player2)
		attacker.Eliminate()
		defender.Eliminate()

		mem.UpdateFromCombat(attackerPos, defenderPos, attacker, defender, 10)
		assert.Nil(t, mem.Recall(attackerPos))
		assert.Nil(t, mem.Recall(defenderPos))
	})
}

func TestMemory_ClearAndGetKnownEnemyPositions(t *testing.T) {
	mem := ai.NewMemory()
	player := engine.NewPlayer(1, "enemy", "blue")
	piece := engine.NewPiece(models.Miner, &player)

	mem.Remember(engine.NewPosition(1, 1), piece, 1.0, 1)
	mem.Remember(engine.NewPosition(2, 2), piece, 0.5, 1)

	positions := mem.GetKnownEnemyPositions()
	assert.Len(t, positions, 2)
	assert.Contains(t, positions, engine.NewPosition(1, 1))
	assert.Contains(t, positions, engine.NewPosition(2, 2))

	mem.Clear()
	assert.Empty(t, mem.GetKnownEnemyPositions())
}

func TestMemory_DecayConfidence(t *testing.T) {
	mem := ai.NewMemory()
	player := engine.NewPlayer(1, "enemy", "blue")
	piece := engine.NewPiece(models.Miner, &player)

	// Keep confidence 1.0 (revealed)
	mem.Remember(engine.NewPosition(1, 1), piece, 1.0, 1)
	// Guesses
	mem.Remember(engine.NewPosition(2, 2), piece, 0.8, 1)
	mem.Remember(engine.NewPosition(3, 3), piece, 0.3, 1)

	// Decay
	mem.DecayConfidence(0.5, 0.2) // new confidence: 0.8 * 0.5 = 0.4, 0.3 * 0.5 = 0.15 (decayed to nil)

	entry1 := mem.Recall(engine.NewPosition(1, 1))
	assert.NotNil(t, entry1)
	assert.Equal(t, 1.0, entry1.Confidence)

	entry2 := mem.Recall(engine.NewPosition(2, 2))
	assert.NotNil(t, entry2)
	assert.Equal(t, 0.4, entry2.Confidence)

	entry3 := mem.Recall(engine.NewPosition(3, 3))
	assert.Nil(t, entry3)
}

func TestBaseAI_ObserveAndAnalyze(t *testing.T) {
	player := engine.NewPlayer(0, "ai", "red")
	enemy := engine.NewPlayer(1, "enemy", "blue")

	t.Run("no memory", func(t *testing.T) {
		bAI := ai.NewBaseAI(&player, false)
		assert.Nil(t, bAI.GetMemory())

		// Should not panic or crash
		move := engine.NewMove(engine.NewPosition(0, 0), engine.NewPosition(0, 1), &enemy)
		bAI.AnalyzeMove(move, &enemy, 1)
		bAI.ObserveCombat(engine.NewPosition(0, 0), engine.NewPosition(0, 1), nil, nil, 1)
	})

	t.Run("with memory", func(t *testing.T) {
		bAI := ai.NewBaseAI(&player, true)
		assert.NotNil(t, bAI.GetMemory())

		piece := engine.NewPiece(models.General, &enemy)
		bAI.GetMemory().Remember(engine.NewPosition(1, 1), piece, 1.0, 1)

		// Test AnalyzeMove updates memory
		move := engine.NewMove(engine.NewPosition(1, 1), engine.NewPosition(1, 2), &enemy)
		bAI.AnalyzeMove(move, &enemy, 2)

		assert.Nil(t, bAI.GetMemory().Recall(engine.NewPosition(1, 1)))
		entry := bAI.GetMemory().Recall(engine.NewPosition(1, 2))
		assert.NotNil(t, entry)
		assert.Equal(t, piece, entry.Piece)

		// Test ObserveCombat updates memory
		attacker := engine.NewPiece(models.General, &enemy)
		defender := engine.NewPiece(models.Miner, &player)
		defender.Eliminate()

		bAI.ObserveCombat(engine.NewPosition(1, 2), engine.NewPosition(1, 3), attacker, defender, 3)
		assert.Nil(t, bAI.GetMemory().Recall(engine.NewPosition(1, 2)))
		entry = bAI.GetMemory().Recall(engine.NewPosition(1, 3))
		assert.NotNil(t, entry)
		assert.Equal(t, attacker, entry.Piece)
	})
}
