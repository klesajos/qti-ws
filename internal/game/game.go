// Package game wires the board, renderer and input together and runs the main
// loop.
package game

import (
	"fmt"

	"github.com/klesajos/qti-ws/internal/board"
	"github.com/klesajos/qti-ws/internal/input"
	"github.com/klesajos/qti-ws/internal/renderer"
)

// Game owns one board and the terminal I/O around it.
type Game struct {
	board    *board.Board
	renderer *renderer.Renderer
	input    *input.Input
}

// New creates a game with a freshly seeded board and switches the terminal
// into raw mode. Call Close when done.
func New() (*Game, error) {
	in, err := input.New()
	if err != nil {
		return nil, err
	}
	return &Game{board: board.New(), renderer: renderer.New(), input: in}, nil
}

// Close restores the terminal.
func (g *Game) Close() error {
	return g.input.Close()
}

// Run plays the game until the player quits or no move is possible.
func (g *Game) Run() {
	g.renderer.Draw(g.board)

	for {
		cmd := g.input.Next()
		if cmd == input.Quit {
			break
		}
		if cmd == input.None {
			continue
		}

		g.board.Move(toDirection(cmd))

		// Place a new tile and redraw.
		g.board.SpawnRandom()
		g.renderer.Draw(g.board)

		if g.board.HasWon() {
			g.renderer.Message("You reached 2048! Keep going or press q.")
		}
		if g.board.IsGameOver() {
			g.renderer.Message(fmt.Sprintf("Game over. Final score: %d", g.board.Score()))
			break
		}
	}
}

// toDirection maps a movement Command to a board Direction. The caller must
// ensure cmd is one of the four movement commands.
func toDirection(cmd input.Command) board.Direction {
	switch cmd {
	case input.Up:
		return board.Up
	case input.Down:
		return board.Down
	case input.Left:
		return board.Left
	default:
		return board.Right
	}
}
