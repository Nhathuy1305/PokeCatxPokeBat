// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"sync"

// 	"github.com/gocolly/colly/v2"
// )

// type Pokemon struct {
// 	ID                int                `json:"id"`
// 	Name              string             `json:"name"`
// 	Types             []string           `json:"types"`
// 	Stats             map[string]int     `json:"stats"`
// 	DamageMultipliers map[string]float64 `json:"damage_multipliers"`
// }

// func main() {
// 	// Create a new collector
// 	c := colly.NewCollector(
// 		colly.AllowedDomains("pokedex.org"),
// 	)

// 	// User agent
// 	c.UserAgent = "Colly - Golang scraping framework"

// 	// Slice to hold all Pokemon data
// 	var pokemons []Pokemon
// 	var mx sync.Mutex

// 	// Iterate over each Pokémon entry
// 	for i := 1; i <= 200; i++ {
// 		url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i)

// 		// Clone the collector to avoid reusing handlers
// 		cloned := c.Clone()

// 		cloned.OnHTML(".detail-panel", func(e *colly.HTMLElement) {
// 			pokemon := Pokemon{
// 				ID:                i,
// 				Name:              e.ChildText(".detail-panel-header"),
// 				Types:             []string{},
// 				Stats:             make(map[string]int),
// 				DamageMultipliers: make(map[string]float64),
// 			}

// 			// Extract types
// 			e.ForEach(".monster-type", func(_ int, el *colly.HTMLElement) {
// 				pokemon.Types = append(pokemon.Types, el.Text)
// 			})

// 			// Extract stats
// 			e.ForEach(".detail-stats-row", func(_ int, el *colly.HTMLElement) {
// 				statName := el.ChildText("span:first-child")
// 				statValue, _ := strconv.Atoi(el.ChildText(".stat-bar-fg"))
// 				pokemon.Stats[statName] = statValue
// 			})

// 			// Extract damage when attacked
// 			e.ForEach(".when-attacked-row", func(_ int, el *colly.HTMLElement) {
// 				el.ForEach("span.monster-type", func(index int, elem *colly.HTMLElement) {
// 					damageType := elem.Text
// 					multiplierText := el.DOM.Find("span.monster-multiplier").Eq(index).Text()
// 					multiplierText = strings.Trim(multiplierText, "x")
// 					multiplier, _ := strconv.ParseFloat(multiplierText, 64)
// 					pokemon.DamageMultipliers[damageType] = multiplier
// 				})
// 			})

// 			mx.Lock()
// 			pokemons = append(pokemons, pokemon)
// 			mx.Unlock()
// 		})

// 		cloned.Visit(url)
// 	}

// 	// Serialize data to JSON
// 	file, err := json.MarshalIndent(pokemons, "", "  ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Write to file
// 	if err = os.WriteFile("../../data/pokedex.json", file, 0644); err != nil {
// 		log.Fatal(err)
// 	}
// }
