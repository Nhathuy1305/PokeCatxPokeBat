package crawler

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ParsePokemon extracts data from HTML element into a Pokemon struct.
func ParsePokemon(e *colly.HTMLElement) Pokemon {
	var pokemon Pokemon
	pokemon.Name = e.ChildText(".detail-panel-header")
	pokemon.Types = make([]string, 0)
	pokemon.Stats = make(map[string]int)                 // Initialize the map for Stats
	pokemon.DamageMultipliers = make(map[string]float64) // Initialize the map for DamageMultipliers

	// Extract types
	e.ForEach(".monster-type", func(_ int, el *colly.HTMLElement) {
		pokemon.Types = append(pokemon.Types, el.Text)
	})

	// Extract stats
	e.ForEach(".detail-stats-row", func(_ int, el *colly.HTMLElement) {
		statName := el.ChildText("span:first-child")
		statValue, _ := strconv.Atoi(el.ChildText(".stat-bar-fg"))
		pokemon.Stats[statName] = statValue
	})

	// Extract damage when attacked
	e.ForEach(".when-attacked-row", func(_ int, el *colly.HTMLElement) {
		el.ForEach("span.monster-type", func(index int, elem *colly.HTMLElement) {
			damageType := elem.Text
			multiplierText := el.DOM.Find("span.monster-multiplier").Eq(index).Text()
			multiplierText = strings.Trim(multiplierText, "x")
			multiplier, _ := strconv.ParseFloat(multiplierText, 64)
			pokemon.DamageMultipliers[damageType] = multiplier
		})
	})

	return pokemon
}
