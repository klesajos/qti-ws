// Package renderer draws the board to standard output using ANSI escape codes.
package renderer

import (
	"fmt"
	"io"
	"os"

	"github.com/klesajos/qti-ws/internal/board"
)

// cellWidth is the number of characters per tile, e.g. "  2048".
const cellWidth = 6

// eol ends every line. The terminal is in raw mode while the game runs, so a
// bare "\n" would only move down without returning to the left edge.
const eol = "\r\n"

// Renderer writes the board and messages to an output stream.
type Renderer struct {
	out io.Writer
}

// New returns a renderer that writes to standard output.
func New() *Renderer {
	return &Renderer{out: os.Stdout}
}

// Draw clears the screen and renders the current board state and score.
func (r *Renderer) Draw(b *board.Board) {
	// Clear screen and move the cursor to the top-left corner.
	fmt.Fprint(r.out, "\x1b[2J\x1b[H")
	fmt.Fprintf(r.out, "  2048  —  score: %d%s%s", b.Score(), eol, eol)

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			value := b.At(row, col)
			if value == 0 {
				fmt.Fprintf(r.out, "%*s", cellWidth, " .")
			} else {
				fmt.Fprintf(r.out, "%*d", cellWidth, value)
			}
		}
		fmt.Fprint(r.out, eol, eol)
	}

	fmt.Fprint(r.out, "  Arrows / WASD to move,  q to quit", eol)
}

// Message prints a one-off message below the board (e.g. "You win!").
func (r *Renderer) Message(text string) {
	fmt.Fprint(r.out, eol, "  ", text, eol)
}
