package main

import (
	"fmt"
	"pokemon/internal/pokedex"
)

func main() {
	url := "https://example.com/pokedex"
	pokemons, err := pokedex.FetchPokedexData(url)
	if err != nil {
		fmt.Println("Error fetching pokedex data:", err)
		return
	}
	err = pokedex.SavePokedexData(pokemons, "data/pokedex.json")
	if err != nil {
		fmt.Println("Error saving pokedex data:", err)
		return
	}
	fmt.Println("Pokedex data saved successfully")
}
