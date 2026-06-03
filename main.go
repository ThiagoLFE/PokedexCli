package main

import (
	"bufio"
	"fmt"
	"os"
)

type Controller struct {
	cmd map[string]cliCommand
}

type cliCommand struct {
	Name        string
	Description string
	Callback    func() error
}

func main() {
	// putting into main function to prevent initialization cycle problems

	scan := bufio.NewScanner(os.Stdin)
	for {
		controller := newController()
		fmt.Print("Pokedex > ")
		// Waits until a enter from terminal
		scan.Scan()

		// Verify if happens an error on read
		if err := scan.Err(); err != nil {
			fmt.Printf("Error to read the input: %w\n", err)
		}

		// Taking the user input
		input := scan.Text()

		cmd, ok := controller.cmd[input]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.Callback(); err != nil {
			fmt.Printf("%s", err.Error())
		}

		fmt.Printf("Your command was: %v\n", input)
	}
}

func newController() *Controller {
	c := &Controller{}

	c.cmd = map[string]cliCommand{
		"exit": cliCommand{
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    c.commandExit,
		},
		"help": cliCommand{
			Name:        "help",
			Description: "Display a help messages",
			Callback:    c.commandHelp,
		},
	}
	return c
}

func (c *Controller) commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (c *Controller) commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range c.cmd {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}
