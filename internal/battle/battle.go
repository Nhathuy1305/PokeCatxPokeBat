package battle

import (
	"math/rand"
	"pokemon/internal/pokemon"
	"time"
)

type Player struct {
	Name     string
	Pokemons []pokemon.Pokemon
}

func NewPlayer(name string, pokemons []pokemon.Pokemon) Player {
	return Player{Name: name, Pokemons: pokemons}
}

func (p *Player) Battle(opponent *Player) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	for i := 0; i < len(p.Pokemons) && i < len(opponent.Pokemons); i++ {
		if p.Pokemons[i].Battle(&opponent.Pokemons[i], r) {
			continue
		} else {
			return opponent.Name
		}
	}
	return p.Name
}

func (p *pokemon.Pokemon) Battle(opponent *pokemon.Pokemon) bool {
	for p.Attributes.HP > 0 && opponent.Attributes.HP > 0 {
		p.attack(opponent)
		if opponent.Attributes.HP > 0 {
			opponent.attack(p)
		}
	}
	return p.Attributes.HP > 0
}

func (p *pokemon.Pokemon) attack(opponent *pokemon.Pokemon) {
	if rand.Intn(2) == 0 {
		damage := p.Attributes.Attack - opponent.Attributes.Defense
		if damage < 0 {
			damage = 0
		}
		opponent.Attributes.HP -= damage
	} else {
		damage := p.Attributes.SpecialAttack*2 - opponent.Attributes.SpecialDefense
		if damage < 0 {
			damage = 0
		}
		opponent.Attributes.HP -= damage
	}
}
