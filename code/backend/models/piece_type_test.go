package models_test

import (
	"digital-innovation/gostrategy/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPieceType_Getters(t *testing.T) {
	pt := models.NewPieceType("Spy", 1, true, "Finds bombs", "spy-icon", 1, 100)

	assert.Equal(t, "Spy", pt.GetName())
	assert.Equal(t, byte(1), pt.GetRank())
	assert.True(t, pt.IsMovable())
	assert.Equal(t, "Finds bombs", pt.GetDescription())
	assert.Equal(t, "spy-icon", pt.GetIcon())
	assert.Equal(t, 1, pt.GetCount())
	assert.Equal(t, 100, pt.GetStrategicValue())
}
