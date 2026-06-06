package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PokeClient struct {
	BaseUrl string
	Client  *http.Client
}

func GetPokeClient() *PokeClient {
	pc := &PokeClient{
		BaseUrl: "https://pokeapi.co/api/v2",
	}
	pc.Client = &http.Client{
		Timeout: 5 * time.Second,
	}

	return pc
}

type LocationsPaginated struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"Url"`
	} `json:"results"`
}

type PokemonInArea struct {
	Pokemon struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"pokemon"`
	VersionDetails []struct {
		Version struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"version_details"`
		MaxChance        int `json:"max_chance"`
		EncounterDetails []struct {
			MinLevel        int   `json:"min_level"`
			MaxLevel        int   `json:"max_level"`
			ConditionValues []any `json:"condition_values"`
			Chance          int   `json:"chance"`
			Method          struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"method"`
		} `json:"encounter_details"`
	}
}

type LocationArea struct {
	Name              string          `json:"name"`
	ID                int             `json:"id"`
	GameIndex         int             `json:"game_index"`
	PokemonEncounters []PokemonInArea `json:"pokemon_encounters"`

	// This items bellow have into the endpoint that we request, but they is useless for us
	// EncounterMethodRates []any `json:"encounter_method_rates"`
	// Location any `json:"location"`
	// Names []any `json:"names"`
}

type Pokemon struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	BaseExperience         int    `json:"base_experience"`
	IsDefault              bool   `json:"is_default"`
	Order                  int    `json:"order"`
	Weight                 int    `json:"weight"`
	LocationAreaEncounters string `json:"location_area_encounters"`

	// This items bellow have into the endpoint that we request, but they is useless for us
	// Abilities              []any  `json:"abilities"`
	// Forms                  []any  `json:"forms"`
	// GameIndices            []any  `json:"game_indices"`
	// HeldItems              []any  `json:"held_items"`
	// Moves                  []any  `json:"moves"`
	// PastTypes              []any  `json:"past_types"`
	// PastAbilities          []any  `json:"past_abilities"`
	// PastStats              []any  `json:"past_stats"`
	// Sprites                any    `json:"sprites"`
	// Cries                  any    `json:"cries"`
	// Species                any    `json:"species"`
	// Stats                  []any  `json:"stats"`
	// Types                  any    `json:"types"`
}

func (pk *PokeClient) GetLocations(url string) (LocationsPaginated, error) {
	if len(url) == 0 {
		return LocationsPaginated{}, fmt.Errorf("Miss URL")
	}
	res, err := pk.Client.Get(url)

	if err != nil {
		return LocationsPaginated{}, fmt.Errorf("Failed to get locations: %w", err)
	}
	defer res.Body.Close()

	var locations LocationsPaginated
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&locations); err != nil {
		return LocationsPaginated{}, fmt.Errorf("Failed to decode response: %w", err)
	}

	return locations, nil
}

func (pk *PokeClient) GetLocationArea(url string) (LocationArea, error) {
	if len(url) == 0 {
		return LocationArea{}, fmt.Errorf("Miss URL")
	}
	res, err := pk.Client.Get(url)

	if err != nil {
		return LocationArea{}, fmt.Errorf("Fail to request %w", err)
	}

	defer res.Body.Close()

	var locationArea LocationArea

	if res.StatusCode == http.StatusNotFound {
		return LocationArea{}, fmt.Errorf("Area Not Found\n")
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locationArea); err != nil {
		return LocationArea{}, err
	}

	return locationArea, nil
}

func (pk *PokeClient) GetPokemon(url string) (Pokemon, error) {
	if len(url) == 0 {
		return Pokemon{}, fmt.Errorf("Miss URL")
	}

	res, err := pk.Client.Get(url)

	if err != nil {
		return Pokemon{}, err
	}

	defer res.Body.Close()
	var pokemon Pokemon

	if res.StatusCode == http.StatusNotFound {
		return Pokemon{}, fmt.Errorf("Pokemon Not Found\n")
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&pokemon); err != nil {
		return Pokemon{}, fmt.Errorf("Error: %w", err)
	}

	return pokemon, nil
}
