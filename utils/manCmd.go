package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/inancgumus/screen"
)

type Man struct {
	Name        string     `json:"name"`
	Synopsis    string     `json:"synopsis"`
	Description string     `json:"description"`
	Options     [][]string `json:"options"`
	Examples    [][]string `json:"examples"`
}

func ManCmd(tokens []string) {
	screen.Clear()
	screen.MoveTopLeft()
	manualPage := "man"
	if len(tokens) == 2 {
		manualPage = tokens[1]
	}
	file, err := os.Open(strings.Join([]string{"man/", manualPage, ".json"}, ""))
	if err != nil {
		// no man page for given command
		fmt.Printf("bash: %s: command not found\n", tokens[0])
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		// Handle error
		os.Exit(5)
	}
	var manPage Man
	if err := json.Unmarshal(data, &manPage); err != nil {
		os.Exit(6)
	}

	w, _ := screen.Size()
	fmt.Printf("%-*sUser Commands%*s\n\n", (w-13)/2, manualPage, (w-13)/2, manualPage)
	fmt.Printf("NAME\n\t%s\n\n", manPage.Name)
	fmt.Printf("SYNOPSIS\n\t%s\n\n", manPage.Synopsis)
	fmt.Printf("DESCRIPTION\n\t%s\n\n", manPage.Description)
	for _, option := range manPage.Options {
		fmt.Printf("\t%s\n\t\t%s\n\n", option[0], option[1])
	}
	fmt.Print("EXAMPLES\n\n")
	for _, example := range manPage.Examples {
		fmt.Printf("%s\n\t$ %s\n", example[0], example[1])
	}
}
