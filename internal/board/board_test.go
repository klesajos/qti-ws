package board

import "testing"

// slid runs SlideLineLeft on a line literal and returns the result.
// Convenience for the cases where the score gained is irrelevant.
func slid(line Line) Line {
	SlideLineLeft(&line)
	return line
}

// ---------------------------------------------------------------------------
// SlideLineLeft — the heart of the game.
// ---------------------------------------------------------------------------

func TestSlideLineLeft_EmptyLineStaysEmpty(t *testing.T) {
	if got, want := slid(Line{0, 0, 0, 0}), (Line{0, 0, 0, 0}); got != want {
		t.Errorf("slid(empty) = %v, want %v", got, want)
	}
}

func TestSlideLineLeft_TilesCompactWithoutMerging(t *testing.T) {
	cases := []struct{ in, want Line }{
		{Line{0, 2, 0, 4}, Line{2, 4, 0, 0}},
		{Line{2, 4, 8, 16}, Line{2, 4, 8, 16}},
	}
	for _, tc := range cases {
		if got := slid(tc.in); got != tc.want {
			t.Errorf("slid(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestSlideLineLeft_SinglePairMerges(t *testing.T) {
	line := Line{2, 2, 0, 0}
	gained := SlideLineLeft(&line)
	if want := (Line{4, 0, 0, 0}); line != want {
		t.Errorf("line = %v, want %v", line, want)
	}
	if gained != 4 {
		t.Errorf("gained = %d, want 4", gained)
	}
}

func TestSlideLineLeft_TwoSeparatePairsMergeIntoTwoTiles(t *testing.T) {
	line := Line{2, 2, 2, 2}
	gained := SlideLineLeft(&line)
	if want := (Line{4, 4, 0, 0}); line != want {
		t.Errorf("line = %v, want %v", line, want)
	}
	if gained != 8 {
		t.Errorf("gained = %d, want 8", gained)
	}
}

func TestSlideLineLeft_OnlyMatchingNeighboursMerge(t *testing.T) {
	if got, want := slid(Line{4, 2, 2, 0}), (Line{4, 4, 0, 0}); got != want {
		t.Errorf("slid({4,2,2,0}) = %v, want %v", got, want)
	}
}

// TODO (participants): add a test for the line {4, 4, 8, 0}.
//   What *should* the result be? Run it and see what actually happens.

// ---------------------------------------------------------------------------
// Board.Move — directions.
// ---------------------------------------------------------------------------

func TestMove_LeftCollapsesEachRow(t *testing.T) {
	b := FromGrid(Grid{{2, 2, 0, 0}, {0, 4, 4, 0}, {0, 0, 0, 0}, {8, 0, 8, 0}}, 0)
	if !b.Move(Left) {
		t.Fatal("Move(Left) = false, want true")
	}
	want := Grid{{4, 0, 0, 0}, {8, 0, 0, 0}, {0, 0, 0, 0}, {16, 0, 0, 0}}
	if got := b.Grid(); got != want {
		t.Errorf("grid = %v, want %v", got, want)
	}
}

func TestMove_RightCollapsesTowardsTheRightEdge(t *testing.T) {
	b := FromGrid(Grid{{2, 2, 0, 0}}, 0)
	b.Move(Right)
	if got := b.At(0, 3); got != 4 {
		t.Errorf("At(0,3) = %d, want 4", got)
	}
}

func TestMove_UpCollapsesColumns(t *testing.T) {
	b := FromGrid(Grid{{2, 0, 0, 0}, {2, 0, 0, 0}}, 0)
	b.Move(Up)
	if got := b.At(0, 0); got != 4 {
		t.Errorf("At(0,0) = %d, want 4", got)
	}
}

func TestMove_ThatChangesNothingReturnsFalse(t *testing.T) {
	b := FromGrid(Grid{{2, 4, 2, 4}, {4, 2, 4, 2}, {2, 4, 2, 4}, {4, 2, 4, 2}}, 0)
	if b.Move(Left) {
		t.Error("Move(Left) = true, want false")
	}
}

func TestMove_ScoreAccumulatesAcrossMoves(t *testing.T) {
	b := FromGrid(Grid{{2, 2, 0, 0}}, 0)
	b.Move(Left) // +4
	if got := b.Score(); got != 4 {
		t.Errorf("Score() = %d, want 4", got)
	}
}

// ---------------------------------------------------------------------------
// Win / lose detection.
// ---------------------------------------------------------------------------

func TestHasWon_TriggersAt2048(t *testing.T) {
	b := FromGrid(Grid{{2048, 0, 0, 0}}, 0)
	if !b.HasWon() {
		t.Error("HasWon() = false, want true")
	}
}

func TestIsGameOver_BoardWithAnEmptyCellIsNotOver(t *testing.T) {
	b := FromGrid(Grid{{2, 4, 2, 4}, {4, 2, 4, 2}, {2, 4, 2, 4}, {4, 2, 4, 0}}, 0)
	if b.IsGameOver() {
		t.Error("IsGameOver() = true, want false")
	}
}

func TestIsGameOver_FullyLockedBoardIsOver(t *testing.T) {
	b := FromGrid(Grid{{2, 4, 2, 4}, {4, 2, 4, 2}, {2, 4, 2, 4}, {4, 2, 4, 2}}, 0)
	if !b.IsGameOver() {
		t.Error("IsGameOver() = false, want true")
	}
}

// TODO (participants): a *full* board can still be playable if it has two
//   equal neighbours. Write a test for a full board that contains a mergeable
//   pair and assert that IsGameOver() returns false.
