// Command 2048 runs the terminal version of the game.
package main

import (
	"fmt"
	"os"

	"github.com/klesajos/qti-ws/internal/game"
)

func main() {
	g, err := game.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "2048: cannot switch the terminal to raw mode:", err)
		os.Exit(1)
	}
	defer g.Close()

	g.Run()
}
