package engine_test

import (
	"digital-innovation/stratego/engine"
	"testing"
)

func TestEquals(t *testing.T) {
	// setup
	position1 := engine.NewPosition(6, 7)
	position2 := engine.NewPosition(7, 6)
	position3 := engine.NewPosition(6, 7)

	// test
	if position1.Equals(position2) {
		t.Errorf("Expected positions to be unequal")
	}

	if !position1.Equals(position3) {
		t.Errorf("Expected positions to be equal")
	}
}

func TestCopy(t *testing.T) {
	// setup
	original := engine.NewPosition(3, 4)

	// test
	copied := original.Copy()
	if !original.Equals(copied) {
		t.Errorf("Expected copied position to be equal to original")
	}

	// modify copied position
	copied.X = 5
	if original.Equals(copied) {
		t.Errorf("Expected original position to remain unchanged after modifying the copy")
	}
}

func TestToLeft(t *testing.T) {
	// setup
	original := engine.NewPosition(5, 5)

	// test
	left := original.ToLeft()
	expected := engine.NewPosition(4, 5)
	if !left.Equals(expected) {
		t.Errorf("Expected position to the left to be %v, got %v", expected, left)
	}
}

func TestToRight(t *testing.T) {
	// setup
	original := engine.NewPosition(5, 5)

	// test
	right := original.ToRight()
	expected := engine.NewPosition(6, 5)
	if !right.Equals(expected) {
		t.Errorf("Expected position to the right to be %v, got %v", expected, right)
	}
}
func TestToUp(t *testing.T) {
	// setup
	original := engine.NewPosition(5, 5)

	// test
	up := original.ToUp()
	expected := engine.NewPosition(5, 4)
	if !up.Equals(expected) {
		t.Errorf("Expected position above to be %v, got %v", expected, up)
	}
}
func TestToDown(t *testing.T) {
	// setup
	original := engine.NewPosition(5, 5)

	// test
	down := original.ToDown()
	expected := engine.NewPosition(5, 6)
	if !down.Equals(expected) {
		t.Errorf("Expected position below to be %v, got %v", expected, down)
	}
}
