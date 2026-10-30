package datasource

import (
	"reflect"
	"src/internal/domain"
	"testing"
)

func TestToDomainBoard(t *testing.T) {
	expectedBoard := domain.Board{
		{0, 0, 2},
		{1, 0, 0},
		{0, 1, 0},
	}

	boardStr := "0, 0, 2, 1, 0, 0, 0, 1, 0"

	dboard, err := toDomainBoard(boardStr)
	if err != nil {
		t.Errorf("ToDomainBoard() get error %v, want not error", err)
		return
	}

	if !reflect.DeepEqual(dboard, expectedBoard) {
		t.Errorf("toDomainBoard() get %v, but want %v", dboard, expectedBoard)
	}
}

func TestToDomainBoardError(t *testing.T) {
	boardStr := "0, 0, 2, 1, 0, u, 0, 1, 0"

	_, err := toDomainBoard(boardStr)
	if err == nil {
		t.Errorf("toDomainBoard() get error %v, want error", err)
		return
	}
}

func TestToEntityBoard(t *testing.T) {
	board := domain.Board{
		{0, 0, 2},
		{1, 0, 0},
		{0, 1, 0},
	}

	expectedBoardtStr := "0, 0, 2, 1, 0, 0, 0, 1, 0"

	boardStr := toEntityBoard(board)

	if !reflect.DeepEqual(boardStr, expectedBoardtStr) {
		t.Errorf("toEntityBoard() get %s, but want %s", boardStr, expectedBoardtStr)
	}
}
