package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerAddress string `yaml:"server_address"`
	PokeWorldSize int    `yaml:"pokeworld_size"`
	PokemonWave   int    `yaml:"pokemon_wave"`
}

var Conf Config

func InitConfig() {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
