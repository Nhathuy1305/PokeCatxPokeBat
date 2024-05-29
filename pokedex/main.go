package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

type Pokemon struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Types             []string          `json:"types"`
	Stats             map[string]string `json:"stats"`
	DamageMultipliers map[string]string `json:"damage_multipliers"`
	EXP               string            `json:"exp"`
}

func main() {
	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Timeout for our operations
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var pokemons []Pokemon

	// Navigate and extract data from pokedex.org
	for i := 1; i <= 3; i++ {
		var pokemon Pokemon
		err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i)),
			chromedp.Sleep(5*time.Second),
			chromedp.Evaluate(`document.querySelector(".detail-header .detail-national-id").innerText.replace("#", "")`, &pokemon.ID),
			chromedp.Evaluate(`document.querySelector(".detail-panel-header").innerText`, &pokemon.Name),
			chromedp.Evaluate(`Array.from(document.querySelectorAll('.detail-types span.monster-type')).map(elem => elem.innerText)`, &pokemon.Types),
			chromedp.Evaluate(`Object.fromEntries(Array.from(document.querySelectorAll('.detail-stats-row')).map(row => {
				const label = row.querySelector('span:first-child').innerText;
				const value = row.querySelector('.stat-bar-fg').innerText;
				return [label, value];
			}))`, &pokemon.Stats),
			chromedp.Evaluate(`Object.fromEntries(Array.from(document.querySelectorAll('.when-attacked-row')).map(row => {
				const types = row.querySelectorAll('span.monster-type');
				const multipliers = row.querySelectorAll('span.monster-multiplier');
				return Array.from(types).map((type, index) => {
					const key = type.innerText.trim();
					const value = multipliers[index]?.innerText.trim() || '';
					return key && value ? [key, value] : null;
				}).filter(pair => pair !== null);
			}).flat())`, &pokemon.DamageMultipliers),
		)
		if err != nil {
			log.Fatalf("Failed to extract data for ID %d: %v", i, err)
		}
		pokemons = append(pokemons, pokemon)
		fmt.Printf("Crawled data for Pokemon ID %d\n", i)
	}

	// Extract EXP data from Bulbapedia
	var expData [][2]string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield_(Generation_IX)"),
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('table.roundy tbody tr')).slice(1).map(row => {
				const cols = row.querySelectorAll('td');
				const id = cols[0]?.innerText.trim().replace(/^0+/, '');
				const exp = cols[3]?.innerText.trim();
				return id && exp ? [id, exp] : null;
			}).filter(entry => entry !== null)`, &expData),
	)
	if err != nil {
		log.Fatalf("Failed to extract EXP data from Bulbapedia: %v", err)
	}

	// Convert extracted data to a map
	expMap := make(map[string]string)
	for _, entry := range expData {
		expMap[entry[0]] = entry[1]
	}

	// Merge EXP data with existing Pokémon data
	for i := range pokemons {
		if exp, found := expMap[pokemons[i].ID]; found {
			pokemons[i].EXP = exp
		}
	}

	// Save to JSON file
	file, err := os.Create("pokedex.json")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(pokemons); err != nil {
		log.Fatal("Cannot encode to JSON", err)
	}
}
