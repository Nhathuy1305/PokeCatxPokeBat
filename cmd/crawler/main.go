package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pokedex/pkg/crawler"

	"github.com/gocolly/colly/v2"
)

func main() {
	crawlerInstance := crawler.NewCrawler()

	var pokemons []crawler.Pokemon

	// Start crawling process
	for i := 1; i <= 200; i++ {
		url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i)
		crawlerInstance.Collector.OnHTML(".detail-panel", func(e *colly.HTMLElement) {
			// Ensure a new Pokemon instance is created each time
			pokemon := crawler.ParsePokemon(e)

			// Use the correct index in a safe manner
			crawlerInstance.Mutex.Lock()
			pokemons = append(pokemons, pokemon)
			crawlerInstance.Mutex.Unlock()
		})
		crawlerInstance.Collector.Visit(url)
	}

	crawlerInstance.Collector.Wait() // Wait for all collectors to finish

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
}
