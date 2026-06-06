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

type Pokemon struct {
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
	Name              string    `json:"name"`
	ID                int       `json:"id"`
	GameIndex         int       `json:"game_index"`
	PokemonEncounters []Pokemon `json:"pokemon_encounters"`

	// This items bellow have into the endpoint that we request, but they is useless for us
	// EncounterMethodRates []any `json:"encounter_method_rates"`
	// Location any `json:"location"`
	// Names []any `json:"names"`
}

func (pk *PokeClient) GetLocations(url string) (LocationsPaginated, error) {
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
	res, err := pk.Client.Get(url)

	if err != nil {
		return LocationArea{}, fmt.Errorf("Fail to request %w", err)
	}

	defer res.Body.Close()

	var locationArea LocationArea
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locationArea); err != nil {
		return LocationArea{}, fmt.Errorf("Not Found\n")
	}

	return locationArea, nil
}
