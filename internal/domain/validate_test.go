package domain

import (
	"testing"
)

func TestValidateTrue(t *testing.T) {
	var gs gameService
	incoming := Board{
		{0, X, X},
		{X, 0, O},
		{O, O, X},
	}

	stored := Board{
		{0, X, X},
		{X, 0, O},
		{O, O, 0},
	}

	actual := gs.validate(incoming, stored, X)
	var expected error

	if actual != expected {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestValidateFalseTwo(t *testing.T) {
	var gs gameService

	incoming := Board{
		{0, X, X},
		{X, X, O},
		{O, O, X},
	}

	stored := Board{
		{0, X, X},
		{X, 0, O},
		{O, O, 0},
	}

	actual := gs.validate(incoming, stored, X)
	expected := ErrInvalidBoard

	if actual.Error() != expected.Error() {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestValidateInvalidNumber(t *testing.T) {
	var gs gameService

	incoming := Board{
		{0, X, X},
		{X, 0, 0},
		{O, O, 0},
	}

	stored := Board{
		{0, X, X},
		{X, 0, O},
		{O, O, 0},
	}

	actual := gs.validate(incoming, stored, X)
	expected := ErrInvalidBoard

	if actual.Error() != expected.Error() {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestValidateStart(t *testing.T) {
	var gs gameService

	incoming := Board{
		{0, 0, 0},
		{X, 0, 0},
		{0, 0, 0},
	}

	stored := Board{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	actual := gs.validate(incoming, stored, X)
	var expected error

	if actual != expected {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestValidateCell(t *testing.T) {
	var gs gameService

	incoming := Board{
		{0, 0, 0},
		{7, 0, 0},
		{0, 0, 0},
	}

	stored := Board{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	actual := gs.validate(incoming, stored, X)
	expected := ErrInvalidCell

	if actual != expected {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}
