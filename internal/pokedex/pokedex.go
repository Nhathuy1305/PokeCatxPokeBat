package pokedex

import (
	"pokemon/pkg/logger"
)

var Pokedex []Pokemon

func InitPokedex(filename string) {
	pokemons, err := LoadPokedexData(filename)
	if err != nil {
		logger.Error("Failed to load Pokedex:", err)
		return
	}
	Pokedex = pokemons
	logger.Info("Pokedex initialized with", len(Pokedex), "pokemons")
}
