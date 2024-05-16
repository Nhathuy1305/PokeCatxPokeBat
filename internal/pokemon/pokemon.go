package pokemon

import (
	"math/rand"
	"pokemon/internal/pokedex"
	"pokemon/pkg/utils"
	"time"
)

type Pokemon struct {
	Name       string
	Type       string
	Level      int
	Exp        int
	Attributes Attributes
	EV         float64
}

type Attributes struct {
	HP             int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefense int
	Speed          int
}

func NewPokemon(base pokedex.Pokemon) Pokemon {
	rand.Seed(time.Now().UnixNano())
	ev := utils.RandomFloat(0.5, 1.0)
	attributes := GenerateAttributes(base, ev)
	return Pokemon{
		Name:       base.Name,
		Type:       base.Type,
		Level:      1,
		Exp:        0,
		Attributes: attributes,
		EV:         ev,
	}
}

func GenerateAttributes(base pokedex.Pokemon, ev float64) Attributes {
	return Attributes{
		HP:             int(50 * ev),
		Attack:         int(55 * ev),
		Defense:        int(45 * ev),
		SpecialAttack:  int(50 * ev),
		SpecialDefense: int(50 * ev),
		Speed:          int(60 * ev),
	}
}

func (p *Pokemon) LevelUp() {
	if p.Exp >= p.Level*100 {
		p.Level++
		p.Exp = 0
		p.Attributes = GenerateAttributes(pokedex.Pokemon{Name: p.Name, Type: p.Type}, p.EV)
	}
}
