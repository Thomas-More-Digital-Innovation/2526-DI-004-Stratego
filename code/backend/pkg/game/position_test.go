package game

import (
	"testing"
)

func TestEquals(t *testing.T) {
	// setup
	position1 := NewPosition(6, 7)
	position2 := NewPosition(7, 6)
	position3 := NewPosition(6, 7)

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
	original := NewPosition(3, 4)

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
	original := NewPosition(5, 5)

	// test
	left := original.ToLeft()
	expected := NewPosition(4, 5)
	if !left.Equals(expected) {
		t.Errorf("Expected position to the left to be %v, got %v", expected, left)
	}
}

func TestToRight(t *testing.T) {
	// setup
	original := NewPosition(5, 5)

	// test
	right := original.ToRight()
	expected := NewPosition(6, 5)
	if !right.Equals(expected) {
		t.Errorf("Expected position to the right to be %v, got %v", expected, right)
	}
}
func TestToUp(t *testing.T) {
	// setup
	original := NewPosition(5, 5)

	// test
	up := original.ToUp()
	expected := NewPosition(5, 4)
	if !up.Equals(expected) {
		t.Errorf("Expected position above to be %v, got %v", expected, up)
	}
}
func TestToDown(t *testing.T) {
	// setup
	original := NewPosition(5, 5)

	// test
	down := original.ToDown()
	expected := NewPosition(5, 6)
	if !down.Equals(expected) {
		t.Errorf("Expected position below to be %v, got %v", expected, down)
	}
}

func TestPositionString(t *testing.T) {
	pos := NewPosition(0, 0)
	if pos.String() != "(A,0)" {
		t.Errorf("Expected (A,0), got %s", pos.String())
	}
	pos2 := NewPosition(9, 9)
	if pos2.String() != "(J,9)" {
		t.Errorf("Expected (J,9), got %s", pos2.String())
	}
}
