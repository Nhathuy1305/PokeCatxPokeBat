package catching

import (
	"fmt"
	"math/rand"
	"pokemon/internal/pokemon"
	"time"
)

type PokeWorld struct {
	Size       int
	Players    map[string]*Player
	PokemonMap map[string]*pokemon.Pokemon
}

func NewPokeWorld(size int) *PokeWorld {
	return &PokeWorld{
		Size:       size,
		Players:    make(map[string]*Player),
		PokemonMap: make(map[string]*pokemon.Pokemon),
	}
}

func (pw *PokeWorld) AddPlayer(name string) *Player {
	player := NewPlayer(name)
	pw.Players[name] = player
	player.X = rand.Intn(pw.Size)
	player.Y = rand.Intn(pw.Size)
	return player
}

func (pw *PokeWorld) SpawnPokemons(pokemons []pokemon.Pokemon, count int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		p := pokemons[rand.Intn(len(pokemons))]
		x := rand.Intn(pw.Size)
		y := rand.Intn(pw.Size)
		pw.PokemonMap[fmt.Sprintf("%d:%d", x, y)] = &p
	}
}

func (pw *PokeWorld) MovePlayer(player *Player, direction string) {
	switch direction {
	case "up":
		if player.Y > 0 {
			player.Y--
		}
	case "down":
		if player.Y < pw.Size-1 {
			player.Y++
		}
	case "left":
		if player.X > 0 {
			player.X--
		}
	case "right":
		if player.X < pw.Size-1 {
			player.X++
		}
	}
	pos := fmt.Sprintf("%d:%d", player.X, player.Y)
	if p, exists := pw.PokemonMap[pos]; exists {
		player.CapturePokemon(p)
		delete(pw.PokemonMap, pos)
	}
}

type Player struct {
	Name     string
	X        int
	Y        int
	Pokemons []*pokemon.Pokemon
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:     name,
		Pokemons: []*pokemon.Pokemon{},
	}
}

func (p *Player) CapturePokemon(pokemon *pokemon.Pokemon) {
	p.Pokemons = append(p.Pokemons, pokemon)
	fmt.Println(p.Name, "captured", pokemon.Name)
}
