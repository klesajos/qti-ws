// Package board holds the pure game state and rules for 2048. No I/O lives
// here, which is what makes the package straightforward to unit-test.
package board

import (
	"math/rand/v2"
	"time"
)

// Size is the side length of the (square) board.
const Size = 4

// WinValue is the tile value that counts as a win.
const WinValue = 2048

// Line is a single row or column of the board.
type Line [Size]int

// Grid is the full board: Size rows of Size cells each.
type Grid [Size][Size]int

// Direction is the direction of a move.
type Direction int

// The four move directions.
const (
	Up Direction = iota
	Down
	Left
	Right
)

// Board is the game state: the grid, the score and the RNG used to spawn
// tiles. Create one with New or FromGrid.
type Board struct {
	grid  Grid
	score int
	rng   *rand.Rand
}

// New starts an empty board and seeds it with two random tiles.
func New() *Board {
	b := &Board{rng: newRNG()}
	b.SpawnRandom()
	b.SpawnRandom()
	return b
}

// FromGrid builds a board from an explicit grid (handy for tests). No tiles
// are spawned, so the state is fully deterministic.
func FromGrid(grid Grid, score int) *Board {
	return &Board{grid: grid, score: score, rng: newRNG()}
}

func newRNG() *rand.Rand {
	seed := uint64(time.Now().UnixNano())
	return rand.New(rand.NewPCG(seed, seed^0x9E3779B97F4A7C15))
}

// Score returns the current score.
func (b *Board) Score() int { return b.score }

// Grid returns a copy of the current grid.
func (b *Board) Grid() Grid { return b.grid }

// At returns the tile value at the given row and column (0 = empty).
func (b *Board) At(row, col int) int { return b.grid[row][col] }

// SlideLineLeft slides and merges a single line towards index 0 (the "left")
// and returns the score gained from merges this move. It is exposed as a free
// function so it can be unit-tested in isolation.
func SlideLineLeft(line *Line) int {
	// 1) Compact: push all non-zero tiles to the front, keeping their order.
	var out Line
	n := 0
	for _, value := range *line {
		if value != 0 {
			out[n] = value
			n++
		}
	}

	// 2) Merge equal neighbours, scanning left to right.
	gained := 0
	i := 0
	for i+1 < n {
		if out[i] == out[i+1] {
			out[i] *= 2
			gained += out[i]

			// Remove the consumed tile by shifting the tail left by one.
			for k := i + 1; k < n-1; k++ {
				out[k] = out[k+1]
			}
			out[n-1] = 0
			n--
			// NOTE: index i is intentionally not advanced here.
		} else {
			i++
		}
	}

	*line = out
	return gained
}

// Move applies a move in the given direction. It returns true if the board
// actually changed.
func (b *Board) Move(dir Direction) bool {
	before := b.grid
	gained := 0

	// NOTE: the four branches below are near-identical. They extract a line
	// (row or column, possibly reversed), slide it left, then write it back.
	switch dir {
	case Left:
		for r := 0; r < Size; r++ {
			line := Line(b.grid[r])
			gained += SlideLineLeft(&line)
			b.grid[r] = line
		}
	case Right:
		for r := 0; r < Size; r++ {
			var line Line
			for c := 0; c < Size; c++ {
				line[c] = b.grid[r][Size-1-c]
			}
			gained += SlideLineLeft(&line)
			for c := 0; c < Size; c++ {
				b.grid[r][Size-1-c] = line[c]
			}
		}
	case Up:
		for c := 0; c < Size; c++ {
			var line Line
			for r := 0; r < Size; r++ {
				line[r] = b.grid[r][c]
			}
			gained += SlideLineLeft(&line)
			for r := 0; r < Size; r++ {
				b.grid[r][c] = line[r]
			}
		}
	default: // Down
		for c := 0; c < Size; c++ {
			var line Line
			for r := 0; r < Size; r++ {
				line[r] = b.grid[Size-1-r][c]
			}
			gained += SlideLineLeft(&line)
			for r := 0; r < Size; r++ {
				b.grid[Size-1-r][c] = line[r]
			}
		}
	}

	b.score += gained
	return b.grid != before
}

// SpawnRandom places a single tile (2 with 90% probability, otherwise 4) in a
// random empty cell. It returns false if the board was already full.
func (b *Board) SpawnRandom() bool {
	type cell struct{ row, col int }
	var empties []cell
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			if b.grid[r][c] == 0 {
				empties = append(empties, cell{r, c})
			}
		}
	}
	if len(empties) == 0 {
		return false
	}

	pick := empties[b.rng.IntN(len(empties))]
	value := 2
	if b.rng.IntN(10) == 0 {
		value = 4
	}
	b.grid[pick.row][pick.col] = value
	return true
}

// IsGameOver reports whether no move is possible.
func (b *Board) IsGameOver() bool {
	// A board is "over" only when no move is possible. For now we just look
	// for an empty cell.
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			if b.grid[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

// HasWon reports whether any tile has reached WinValue.
func (b *Board) HasWon() bool {
	for _, row := range b.grid {
		for _, value := range row {
			if value >= WinValue {
				return true
			}
		}
	}
	return false
}
