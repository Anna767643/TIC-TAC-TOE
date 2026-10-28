package domain

import (
	"testing"
)

func TestIsFinishedCol(t *testing.T) {
	var gs gameService
	board := Board{
		{O, X, 0},
		{O, X, 0},
		{O, 0, X},
	}

	actual, winner := gs.isFinished(board)
	winnerExpected := O
	expected := true

	if expected != actual || winner != winnerExpected {
		t.Errorf("Result was incorrect, got: %t, %d, want: %t, %d.", actual, winner, expected, winnerExpected)
	}
}

func TestIsFinishedRow(t *testing.T) {
	var gs gameService
	board := Board{
		{0, 0, 0},
		{O, O, O},
		{O, 0, 0},
	}

	actual, winner := gs.isFinished(board)
	winnerExpected := O
	expected := true

	if expected != actual || winner != winnerExpected {
		t.Errorf("Result was incorrect, got: %t, %d, want: %t, %d.", actual, winner, expected, winnerExpected)
	}
}

func TestIsFinishedRandom(t *testing.T) {
	var gs gameService
	board := Board{
		{0, O, 0},
		{O, 0, X},
		{X, 0, 0},
	}

	actual, winner := gs.isFinished(board)
	winnerExpected := 0
	expected := false

	if expected != actual || winner != winnerExpected {
		t.Errorf("Result was incorrect, got: %t, %d, want: %t, %d.", actual, winner, expected, winnerExpected)
	}
}

func TestIsFinishedCross(t *testing.T) {
	var gs gameService
	board := Board{
		{O, O, X},
		{0, X, O},
		{X, 0, 0},
	}

	actual, winner := gs.isFinished(board)
	winnerExpected := X
	expected := true

	if expected != actual || winner != winnerExpected {
		t.Errorf("Result was incorrect, got: %t, %d, want: %t, %d.", actual, winner, expected, winnerExpected)
	}
}

func TestIsFinishedIsFilledBoard(t *testing.T) {
	var gs gameService
	board := Board{
		{O, X, X},
		{X, X, O},
		{O, O, X},
	}

	actual, winner := gs.isFinished(board)
	winnerExpected := 0
	expected := true

	if expected != actual || winner != winnerExpected {
		t.Errorf("Result was incorrect, got: %t, %d, want: %t, %d.", actual, winner, expected, winnerExpected)
	}
}
