package aivsai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunAiVsAi(t *testing.T) {
	t.Run("successful tournament", func(t *testing.T) {
		summary := runAIvsAI("fato", "fato", 2, false)

		assert.Equal(t, 2, summary.Matches)
		assert.NotEmpty(t, summary.Player1data.Name)
		assert.NotEmpty(t, summary.Player2data.Name)
		assert.Positive(t, summary.AverageRounds)
	})

	t.Run("draw edge case", func(t *testing.T) {
		summary := runAIvsAI("fato", "fato", 0, false)

		assert.Equal(t, 0, summary.Matches)
	})

	t.Run("invalid player 1 AI type panics", func(t *testing.T) {
		assert.Panics(t, func() {
			runAIvsAI("unknown", "fato", 1, false)
		})
	})

	t.Run("invalid player 2 AI type panics", func(t *testing.T) {
		assert.Panics(t, func() {
			runAIvsAI("fato", "unknown", 1, false)
		})
	})
}

func TestRunAIvsAIExported(t *testing.T) {
	// TODO: actually check if formats are correct.

	t.Run("default format", func(_ *testing.T) {
		RunAIvsAI("fato", "fato", 1, "default", false)

	})

	t.Run("markdown format", func(_ *testing.T) {
		RunAIvsAI("fato", "fato", 1, "md", false)
	})
}
