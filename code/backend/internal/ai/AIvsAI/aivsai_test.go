package aivsai

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunAiVsAi(t *testing.T) {
	t.Run("successful tournament with logging", func(t *testing.T) {
		summary := runAIvsAI(models.Fato, models.Fato, 2, true)

		assert.Equal(t, 2, summary.Matches)
		assert.NotEmpty(t, summary.Player1data.Name)
		assert.NotEmpty(t, summary.Player2data.Name)
		assert.Positive(t, summary.AverageRounds)
	})

	t.Run("draw edge case", func(t *testing.T) {
		summary := runAIvsAI(models.Fato, models.Fato, 0, false)

		assert.Equal(t, 0, summary.Matches)
	})

	t.Run("invalid player 1 AI type panics", func(t *testing.T) {
		assert.Panics(t, func() {
			runAIvsAI("unknown", models.Fato, 1, false)
		})
	})

	t.Run("invalid player 2 AI type panics", func(t *testing.T) {
		assert.Panics(t, func() {
			runAIvsAI(models.Fato, "unknown", 1, false)
		})
	})
}

func TestRunAIvsAIExported(t *testing.T) {
	// TODO: actually check if formats are correct.

	t.Run("default format", func(_ *testing.T) {
		RunAIvsAI(models.Fato, models.Fato, 1, "default", true)
	})

	t.Run("markdown format", func(_ *testing.T) {
		RunAIvsAI(models.Fato, models.Fato, 1, "md", true)
	})
}

func TestTrainAI(t *testing.T) {
	tempFile := "temp_train_ai_parameters.json"
	ai.SetFallbackFile(tempFile)
	defer func() { _ = os.Remove(tempFile) }()

	err := TrainAI(models.Heuristic, models.Fato, 1, 1, true)
	assert.NoError(t, err)
}
