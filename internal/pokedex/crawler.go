package pokedex

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Pokemon struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	BaseExp int    `json:"base_exp"`
}

func FetchPokedexData(url string) ([]Pokemon, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pokemons []Pokemon
	err = json.Unmarshal(body, &pokemons)
	if err != nil {
		return nil, err
	}
	return pokemons, nil
}

func SavePokedexData(pokemons []Pokemon, filename string) error {
	data, err := json.MarshalIndent(pokemons, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadPokedexData(filename string) ([]Pokemon, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var pokemons []Pokemon
	err = json.Unmarshal(data, &pokemons)
	if err != nil {
		return nil, err
	}
	return pokemons, nil
}
