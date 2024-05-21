package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pokedex/pkg/crawler"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	crawlerInstance := crawler.NewCrawler()

	var pokemons []crawler.Pokemon

	// Start crawling process
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i)
		id := i // Capture the current value of i
		
		c := crawlerInstance.Collector.Clone() // Clone the collector to avoid reusing handlers
		c.OnHTML(".detail-panel", func(e *colly.HTMLElement) {
			// Ensure a new Pokemon instance is created each time
			pokemon := crawler.ParsePokemon(e)
			pokemon.ID = id // Assign the correct ID

			// Use the correct index in a safe manner
			crawlerInstance.Mutex.Lock()
			pokemons = append(pokemons, pokemon)
			crawlerInstance.Mutex.Unlock()
		})
		
		log.Printf("Visiting URL: %s", url)
		c.Visit(url)
		c.Wait() // Wait for this collector to finish before moving to the next
		time.Sleep(2 * time.Second) // Add a delay to be gentle on the server
	}

	// Serialize data to JSON
	file, err := json.MarshalIndent(pokemons, "", "  ")
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	// Write to file
	dataFilePath := "../../data/pokedex.json" // Ensure the directory exists
	if err = os.WriteFile(dataFilePath, file, 0644); err != nil {
		log.Fatal("Error writing to file:", err)
	}
	log.Println("Data successfully written to file.")
}
