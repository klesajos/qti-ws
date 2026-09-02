// Package input reads single keypresses from the terminal in raw mode (no
// Enter required). It supports the arrow keys and WASD; 'q' quits.
package input

import (
	"bufio"
	"os"

	"golang.org/x/term"
)

// Command is what a keypress maps to.
type Command int

// The commands a keypress can produce.
const (
	None Command = iota
	Up
	Down
	Left
	Right
	Quit
)

// Input owns the terminal's raw mode for the lifetime of the game. Create it
// with New and always call Close so the previous settings are restored.
type Input struct {
	fd       int
	previous *term.State
	reader   *bufio.Reader
}

// New puts the terminal into raw mode and returns a reader for keypresses.
func New() (*Input, error) {
	fd := int(os.Stdin.Fd())
	previous, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	return &Input{fd: fd, previous: previous, reader: bufio.NewReader(os.Stdin)}, nil
}

// Close restores the terminal settings captured by New.
func (in *Input) Close() error {
	return term.Restore(in.fd, in.previous)
}

// Next blocks until a key is pressed and maps it to a Command.
func (in *Input) Next() Command {
	c, err := in.reader.ReadByte()
	if err != nil {
		return Quit // stdin closed
	}
	switch c {
	case 'q', 'Q', 0x03: // 0x03 is Ctrl+C, which raw mode no longer turns into a signal.
		return Quit
	case 'w', 'W':
		return Up
	case 's', 'S':
		return Down
	case 'a', 'A':
		return Left
	case 'd', 'D':
		return Right
	case 0x1b: // Arrow keys arrive as ESC '[' followed by A/B/C/D.
		if next, err := in.reader.ReadByte(); err != nil || next != '[' {
			return None
		}
		final, err := in.reader.ReadByte()
		if err != nil {
			return None
		}
		switch final {
		case 'A':
			return Up
		case 'B':
			return Down
		case 'C':
			return Right
		case 'D':
			return Left
		default:
			return None
		}
	default:
		return None
	}
}
