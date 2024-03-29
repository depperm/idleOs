package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/depperm/idleOs/utils"
	"github.com/inancgumus/screen"
)

type GameState = utils.GameState

func handleInput(input string, gameState *GameState) {
	tokens := strings.Split(strings.TrimSpace(input), " ")
	cmd := tokens[0]
	options := make(map[string]int)
	var positional []string
	if len(strings.TrimSpace(input)) > 0 {
		gameState.Player.History = append(gameState.Player.History, input)
	}
	if len(tokens) > 1 {
		// get options
		for j := 1; j < len(tokens); j++ {
			if strings.HasPrefix(tokens[j], "--") {
				options[tokens[j][2:]] = 1
				// todo some options have # with -t 5
				// grab from positional later on?
				// j += 1
			} else if strings.HasPrefix(tokens[j], "-") {
				if len(tokens[j]) == 2 {
					options[tokens[j][1:]] = 1
				} else {
					for _, flag := range tokens[j][1:] {
						options[string(flag)] = 1
					}
				}
			} else {
				positional = append(positional, tokens[j])
			}
		}
	}
	// fmt.Println(tokens)
	switch cmd {
	case "":
		fmt.Print("")
	case "man":
		utils.ManCmd(tokens)
	case "help":
		fmt.Println("TODO should print something")
	case "whoami":
		fmt.Println(gameState.Player.Username)
	case "history":
		for _, cmd := range gameState.Player.History {
			fmt.Println(cmd)
		}
	case "clear":
		screen.Clear()
		screen.MoveTopLeft()
	case "cd":
		utils.CdCmd(gameState, tokens)
		// fmt.Print("")
		// if len(tokens) == 1 {
		// 	// change to root dir
		// 	gameState.CurrentDir = gameState.Player.Dirs.Name
		// } else {
		// 	dst := strings.Split(strings.TrimRight(tokens[1], "/"), "/")
		// 	fmt.Print(dst)
		// }
	case "pwd":
		fmt.Println(gameState.CurrentDir)
	case "ls":
		utils.LsCmd(gameState, tokens, options)
	}
}

func GameLoop(gameState *GameState) {
	var userInput string
	screen.Clear()
	// TODO
	// IDLE OS info
	// basic info, help, man, examples, etc

	// fmt.Printf("loaded: %+v\n", gameState)
	screen.MoveTopLeft()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("[%s@IDLE %s]$ ", gameState.Player.Username, gameState.CurrentDir)
		scanner.Scan()
		userInput = scanner.Text()
		if userInput == "exit" || userInput == "quit" || userInput == "logout" {
			break
		}
		handleInput(userInput, gameState)
	}
}

func main() {
	gameState, err := utils.LoadGame()
	if err != nil {
		os.Exit(1)
	}
	go utils.AutoSave(&gameState.Player)
	GameLoop(gameState)
}
