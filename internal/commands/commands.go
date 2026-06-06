package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/thiagolfe/pokedexcli/internal/pokeapi"
	"github.com/thiagolfe/pokedexcli/internal/pokecache"
)

type CliCommand struct {
	Description string
	Callback    func(string) error
}

func (c *Config) commandExit(_ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (c *Config) commandHelp(_ string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for key, cmd := range c.Command {
		fmt.Printf("%s: %s\n", key, cmd.Description)
	}
	return nil
}

func (c *Config) commandMap(_ string) error {
	// getting url
	url := c.Client.BaseUrl + "/location-area"
	if c.Next != nil {
		url = *c.Next
	}

	// declaring values
	var locations pokeapi.LocationsPaginated
	var err error

	// getting values
	data, isCached := c.Cache.Get(url)
	if isCached {
		if err := json.Unmarshal(data, &locations); err != nil {
			return err
		}
	} else {
		locations, err = c.Client.GetLocations(url)
		if err != nil {
			return err
		}

		// adding new data to cache
		dataToStore, err := json.Marshal(locations)
		if err != nil {
			return err
		}
		c.Cache.Add(url, dataToStore)
	}

	c.setPages(locations.Previous, locations.Next)

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func (c *Config) commandBackMap(_ string) error {
	// getting url
	url := c.Client.BaseUrl + "/location-area"

	if c.Previous != nil {
		url = *c.Previous
	}

	// declaring values
	var locations pokeapi.LocationsPaginated
	var err error

	// getting values
	data, isCached := c.Cache.Get(url)
	if isCached {
		if err := json.Unmarshal(data, &locations); err != nil {
			return err
		}
	} else {
		locations, err = c.Client.GetLocations(url)
		if err != nil {
			return err
		}

		// adding new data to cache
		dataToStore, err := json.Marshal(locations)
		if err != nil {
			return err
		}
		c.Cache.Add(url, dataToStore)
	}

	c.setPages(locations.Previous, locations.Next)

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func (c *Config) commandExplore(city string) error {
	fullpath := c.Client.BaseUrl + "/location-area/" + city

	var err error
	var area pokeapi.LocationArea

	data, isCache := c.Cache.Get(fullpath)

	// Reading data from cache/request
	if isCache {
		if err := json.Unmarshal(data, &area); err != nil {
			return fmt.Errorf("Error to read the cached data: %w", err)
		}
	} else {
		area, err = c.Client.GetLocationArea(fullpath)
		if err != nil {
			return err
		}

		//Storing the data to cache
		data, err = json.Marshal(area)
		if err != nil {
			return fmt.Errorf("Error to store data into cache: %w", err)
		}
		c.Cache.Add(fullpath, data)
	}

	for _, pokemon := range area.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
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
		Description: "Returns to the last page of locations names",
		Callback:    c.commandBackMap,
	}
	cmd["explore"] = CliCommand{
		Description: "use 'explore <city-name>' for you search pokemons in the area",
		Callback:    c.commandExplore,
	}

	c.Command = cmd
}

type Config struct {
	Command  map[string]CliCommand
	Client   *pokeapi.PokeClient
	Cache    *pokecache.PokeCache
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
	cfg.Cache = pokecache.NewPokeCache(time.Minute * 5)

	return cfg
}

func (c *Config) setPages(previous, next *string) {
	c.Previous = previous
	c.Next = next
}
