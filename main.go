package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/thiagolfe/pokedexcli/internal/commands"
	"github.com/thiagolfe/pokedexcli/internal/repl"
)

func main() {
	// putting into main function to prevent initialization cycle problems
	Config := commands.NewConfig()

	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		// Waits until a enter from terminal
		scan.Scan()

		// Verify if happens an error on read
		if err := scan.Err(); err != nil {
			fmt.Printf("Error to read the input: %s\n", err.Error())
		}

		// Taking the user input
		input := repl.CleanInput(scan.Text())
		commandName := ""
		param := ""

		if len(input) > 0 {
			commandName = input[0]
			if len(input) > 1 {
				param = input[1]
			}
		}

		cmd, ok := Config.Command[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.Callback(param); err != nil {
			fmt.Printf("%s", err.Error())
		}

		fmt.Printf("Your command was: %v\n", input)
	}
}
