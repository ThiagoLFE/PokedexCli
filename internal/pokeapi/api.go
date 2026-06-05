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
