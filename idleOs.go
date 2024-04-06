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

func incLines(gameState *GameState, amt float64) {
	gameState.Player.Lines += amt // todo add bonuses/%/multipliers/etc
}
func decLines(gameState *GameState, amt float64) {
	gameState.Player.Lines -= amt // todo add bonuses/%/multipliers/etc
}

func handleInput(input string, gameState *GameState) {
	tokens := strings.Split(strings.TrimSpace(input), " ")
	cmd := tokens[0]
	options := make(map[string]int)
	var positional []string
	incLines(gameState, 0.25)
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
	if strings.HasPrefix(cmd, "./") {
		// todo check for & and args/positional
		app := strings.TrimSuffix(cmd[2:], ".exe")
		if utils.HasExe(gameState, app) {
			switch app {
			case "code":
				utils.StartCode(gameState)
			default:
				fmt.Printf("missing %s", app)
			}
		} else {
			fmt.Printf("%s: does not exist\n", app)
		}
	} else {
		switch cmd {
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
		case "stats":
			fmt.Printf("Lines: %f\nTODO\n", gameState.Player.Lines)
		case "":
			fmt.Print("")
			decLines(gameState, 0.25)
		default:
			fmt.Printf("%s: command not found\n", cmd)
			decLines(gameState, 0.25)
		}
	}
}

func GameLoop(gameState *GameState) {
	var userInput string

	// fmt.Printf("loaded: %+v\n", gameState)
	// screen.MoveTopLeft()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("[%s@IDLE %s]$ ", gameState.Player.Username, gameState.CurrentDir)
		scanner.Scan()
		userInput = scanner.Text()
		if utils.Contains([]string{"shutdown", "quit", "logout", "exit"}, userInput) {
			break
		}
		handleInput(userInput, gameState)
	}
}

func main() {
	// TODO
	// boot/loading info
	// [ OK ] Start[ed/ing] some service ...
	// IDLE OS info
	// basic info, help, man, examples, etc
	screen.Clear()
	screen.MoveTopLeft()

	fmt.Println("Welcome to IDLE OS emulator game")
	fmt.Println()
	fmt.Println("Starting init process.")
	fmt.Println("INIT: version 1.0-idleos booting")
	fmt.Println("Loading Player Info")
	gameState, err := utils.LoadGame()
	if err != nil {
		fmt.Println("ERROR Loading player info")
		os.Exit(1)
	}
	if gameState.Player.Lines == 0.0 {
		fmt.Println("Loaded New Game Settings")
	} else {
		fmt.Println("Loaded Saved Game Settings")
	}
	fmt.Println("Running Fake Kernel")
	fmt.Printf(" TODO Processing Power\n\n")
	go utils.AutoSave(&gameState.Player)
	GameLoop(gameState)
}
