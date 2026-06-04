package fato

import (
	ai_const "digital-innovation/gostrategy/internal/ai/const"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAI(t *testing.T) {
	player := game.NewPlayer(0, "player", "red")
	ai := NewAI(&player, true)

	assert.NotNil(t, ai.GetPlayer())
	assert.NotNil(t, ai.GetMemory())
	assert.Equal(t, 0.5, ai.GetAggression())
}

func TestAggressionClamping(t *testing.T) {
	player := game.NewPlayer(0, "player", "red")

	t.Run("NewAIWithAggression clamp", func(t *testing.T) {
		aiLow := NewAIWithAggression(&player, true, -0.5)
		assert.Equal(t, 0.0, aiLow.GetAggression())

		aiHigh := NewAIWithAggression(&player, true, 1.5)
		assert.Equal(t, 1.0, aiHigh.GetAggression())
	})

	t.Run("SetAggression clamp", func(t *testing.T) {
		ai := NewAI(&player, true)
		ai.SetAggression(-1.0)
		assert.Equal(t, 0.0, ai.GetAggression())

		ai.SetAggression(2.0)
		assert.Equal(t, 1.0, ai.GetAggression())
	})
}

func TestMakeMove_FullFidelity(t *testing.T) {
	p1 := game.NewPlayer(0, "ai", "red")
	p2 := game.NewPlayer(1, "human", "blue")

	t.Run("scout move detection in AnalyzeMove", func(t *testing.T) {
		aiObj := NewAI(&p1, true)

		normalMove := game.NewMove(game.NewPosition(1, 1), game.NewPosition(1, 2), &p2)
		aiObj.AnalyzeMove(normalMove, &p2, 1)
		assert.Nil(t, aiObj.GetMemory().Recall(normalMove.GetTo()))

		scoutMove := game.NewMove(game.NewPosition(1, 1), game.NewPosition(1, 4), &p2)
		aiObj.AnalyzeMove(scoutMove, &p2, 1)
		remembered := aiObj.GetMemory().Recall(scoutMove.GetTo())
		assert.NotNil(t, remembered)
		assert.Equal(t, ai_const.Scout, remembered.Piece.GetType().GetName())
	})

	t.Run("findExplorationMove towards player 1 enemy territory", func(t *testing.T) {
		p1Alt := game.NewPlayer(1, "ai", "blue")
		aiObj := NewAI(&p1Alt, false)
		p1Piece := game.NewPiece(models.Marshal, &p1Alt)

		boardAlt := game.NewBoard()
		boardAlt.SetPieceAt(game.NewPosition(0, 8), p1Piece)
		p1Alt.AddPiece(p1Piece, game.NewPosition(0, 8))

		move := aiObj.MakeMove(boardAlt)
		assert.NotEqual(t, game.Move{}, move)
		assert.Equal(t, game.NewPosition(0, 9), move.GetTo())
	})
}

func TestEvaluateAttack_Matchups(t *testing.T) {
	p1 := game.NewPlayer(0, "ai", "red")
	p2 := game.NewPlayer(1, "human", "blue")
	aiObj := NewAI(&p1, true)

	miner := game.NewPiece(models.Miner, &p1)
	spy := game.NewPiece(models.Spy, &p1)
	marshal := game.NewPiece(models.Marshal, &p1)
	scout := game.NewPiece(models.Scout, &p1)
	major := game.NewPiece(models.Major, &p1)
	lieutenant := game.NewPiece(models.Lieutenant, &p1)
	sergeant := game.NewPiece(models.Sergeant, &p1)

	t.Run("known Flag target", func(t *testing.T) {
		flag := game.NewPiece(models.Flag, &p2)
		flag.Reveal()
		score := aiObj.evaluateAttack(marshal, flag, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.Equal(t, 10000.0, score)
	})

	t.Run("known Spy vs Marshal", func(t *testing.T) {
		enemyMarshal := game.NewPiece(models.Marshal, &p2)
		enemyMarshal.Reveal()
		score := aiObj.evaluateAttack(spy, enemyMarshal, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.Equal(t, 500.0, score)
	})

	t.Run("known Bomb target", func(t *testing.T) {
		bomb := game.NewPiece(models.Bomb, &p2)
		bomb.Reveal()
		scoreMiner := aiObj.evaluateAttack(miner, bomb, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.Equal(t, 200.0, scoreMiner)
		scoreSpy := aiObj.evaluateAttack(spy, bomb, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.Equal(t, -1000.0, scoreSpy)
	})

	t.Run("known matchups comparison", func(t *testing.T) {
		oppGeneral := game.NewPiece(models.General, &p2)
		oppGeneral.Reveal()
		score1 := aiObj.evaluateAttack(marshal, oppGeneral, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.Equal(t, 10.0, score1)

		oppMarshal := game.NewPiece(models.Marshal, &p2)
		oppMarshal.Reveal()
		score2 := aiObj.evaluateAttack(marshal, oppMarshal, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.Equal(t, -20.0, score2)

		score3 := aiObj.evaluateAttack(spy, oppGeneral, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.Equal(t, -1000.0, score3)
	})

	t.Run("unknown target rank tiers", func(t *testing.T) {
		unrevealed := game.NewPiece(models.General, &p2)

		s1 := aiObj.evaluateAttack(marshal, unrevealed, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.True(t, s1 >= 15.0 && s1 <= 25.0)

		s2 := aiObj.evaluateAttack(major, unrevealed, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.True(t, s2 >= 15.0 && s2 <= 25.0)

		s3 := aiObj.evaluateAttack(lieutenant, unrevealed, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.True(t, s3 >= 5.0 && s3 <= 15.0)

		s4 := aiObj.evaluateAttack(sergeant, unrevealed, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.True(t, s4 >= 0.0 && s4 <= 10.0)

		s5 := aiObj.evaluateAttack(scout, unrevealed, game.NewPosition(0, 1), aiObj.GetMemory())
		assert.True(t, s5 >= -5.0 && s5 <= 5.0)
	})
}
