package aivsai

import (
	"testing"
)

func TestRunAIvsAI(t *testing.T) {
	// Use FATO (fast) for testing
	summary := runAIvsAI("fato", "fato", 2, false)

	if summary.Matches != 2 {
		t.Errorf("Expected 2 matches, got %d", summary.Matches)
	}

	if summary.Player1data.Name == "" || summary.Player2data.Name == "" {
		t.Error("Player names should be populated")
	}

	if summary.AverageRounds <= 0 {
		t.Error("Average rounds should be greater than 0")
	}
}

func TestRunAIvsAIDraw(t *testing.T) {
	// Test with 0 matches just to check edge case
	summary := runAIvsAI("fato", "fato", 0, false)
	
	if summary.Matches != 0 {
		t.Errorf("Expected 0 matches, got %d", summary.Matches)
	}
}
