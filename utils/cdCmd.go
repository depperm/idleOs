package utils

import (
	"fmt"
	"strings"
)

func CdCmd(gameState *GameState, tokens []string) {
	// fmt.Print("")
	if len(tokens) == 1 {
		// change to root dir
		gameState.CurrentDir = gameState.Player.Dirs.Name
	} else {
		dst := strings.Split(strings.TrimRight(tokens[1], "/"), "/")
		fmt.Print(dst)
	}
}
