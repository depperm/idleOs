package utils

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/inancgumus/screen"
)

func typingSpeed(gameState *GameState) float64 {
	//             wpm / seconds * milliseconds / characterPerWord
	return float64(gameState.Player.Wpm) / 60.0 * 1000 / 5.0 // ~133 with 40wpm
}

func LinesPerSecond(gameState *GameState) float64 {
	dat, err := os.ReadFile("idleOs.go")
	if err != nil {
		os.Exit(12)
	}
	contents := strings.TrimSpace(string(dat))
	fileLen := len(contents) // ~3349
	// timed contents printing, took ~450sec for 132 lines = .29lps
	return float64(fileLen) / (typingSpeed(gameState) * 100) // ~25, with * 100 its .25
}

func StartCode(gameState *GameState) {
	screen.Clear()
	fmt.Println("Ctrl-C to stop application at any point")
	shouldStop := false

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				shouldStop = true
				return
			}
		}
	}()

	dat, err := os.ReadFile("idleOs.go")
	if err != nil {
		os.Exit(11)
	}
	contents := strings.TrimSpace(string(dat))
	lines := strings.Split(contents, "\n")
	sleepTime := int(typingSpeed(gameState))
	for _, line := range lines {
		for _, letter := range line {
			fmt.Printf("%s", string(letter))
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
			if shouldStop {
				fmt.Println()
				return
			}
		}
		fmt.Println()
		gameState.Player.Lines += 1 //TODO multiplier
	}
}
