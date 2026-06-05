package commands

import (
	"fmt"
	"os"

	"github.com/thiagolfe/pokedexcli/internal/pokeapi"
)

type CliCommand struct {
	Description string
	Callback    func() error
}

func (c *Config) commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (c *Config) commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for key, cmd := range c.Command {
		fmt.Printf("%s: %s\n", key, cmd.Description)
	}
	return nil
}

func (c *Config) commandMap() error {
	nextPath := ""
	if c.Next != nil {
		nextPath = *c.Next
	}

	locations, err := c.Client.GetLocations(nextPath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	c.setPages(locations.Previous, locations.Next)

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func (c *Config) commandBackMap() error {

	previousPath := ""
	if c.Previous != nil {
		previousPath = *c.Previous
	}
	locations, err := c.Client.GetLocations(previousPath)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	c.setPages(locations.Previous, locations.Next)

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func getCommands(c *Config) {
	cmd := make(map[string]CliCommand, 0)

	cmd["exit"] = CliCommand{
		Description: "Exit the Pokedex",
		Callback:    c.commandExit,
	}
	cmd["help"] = CliCommand{
		Description: "Display a help messages",
		Callback:    c.commandHelp,
	}

	cmd["map"] = CliCommand{
		Description: "Shows the name of locations",
		Callback:    c.commandMap,
	}

	cmd["bmap"] = CliCommand{
		Description: "Shows the name of locations",
		Callback:    c.commandBackMap,
	}

	c.Command = cmd
}

type Config struct {
	Command  map[string]CliCommand
	Client   *pokeapi.PokeClient
	Next     *string
	Previous *string
}

func NewConfig() *Config {
	cfg := &Config{
		Next:     nil,
		Previous: nil,
	}

	getCommands(cfg)
	cfg.Client = pokeapi.GetPokeClient()

	return cfg
}

func (c *Config) setPages(previous, next *string) {
	c.Previous = previous
	c.Next = next
}
