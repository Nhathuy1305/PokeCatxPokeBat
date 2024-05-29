// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"sync"

// 	"github.com/chromedp/chromedp"
// )

// type Pokemon struct {
// 	ID                int                `json:"id"`
// 	Name              string             `json:"name"`
// 	Types             []string           `json:"types"`
// 	Stats             map[string]int     `json:"stats"`
// 	DamageMultipliers map[string]float64 `json:"damage_multipliers"`
// }

// func main() {
// 	opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 		chromedp.Flag("headless", true), // run headless Chrome
// 	)
// 	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
// 	defer cancel()

// 	ctx, cancel = chromedp.NewContext(ctx)
// 	defer cancel()

// 	var pokemons []Pokemon
// 	var mx sync.Mutex

// 	var wg sync.WaitGroup
// 	for i := 1; i <= 5; i++ {
// 		wg.Add(1)
// 		go func(id int) {
// 			defer wg.Done()
// 			pokemon, err := fetchPokemonData(ctx, id)
// 			if err != nil {
// 				log.Printf("Failed to fetch data for ID %d: %v", id, err)
// 				return
// 			}
// 			mx.Lock()
// 			pokemons = append(pokemons, *pokemon)
// 			mx.Unlock()
// 		}(i)
// 	}

// 	wg.Wait()

// 	// Serialize data to JSON
// 	file, err := json.MarshalIndent(pokemons, "", "  ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Write to file
// 	if err = os.WriteFile("pokedex.json", file, 0644); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func fetchPokemonData(ctx context.Context, id int) (*Pokemon, error) {
// 	var name string
// 	var types []string
// 	var stats = make(map[string]int)
// 	var damageMultipliers = make(map[string]float64)

// 	tasks := chromedp.Tasks{
// 		chromedp.Navigate(fmt.Sprintf("https://pokedex.org/#/pokemon/%d", id)),
// 		chromedp.WaitVisible(".detail-panel-header", chromedp.ByQuery),
// 		chromedp.Text(".detail-panel-header", &name),
// 		chromedp.Evaluate(`
// 			Array.from(document.querySelectorAll('.monster-type')).map(el => el.innerText);
// 		`, &types),
// 		chromedp.Evaluate(`
// 			Object.fromEntries(
// 				Array.from(document.querySelectorAll('.detail-stats-row')).map(el => [
// 					el.querySelector('span').textContent.trim(),
// 					parseInt(el.querySelector('.stat-bar-fg').textContent.trim())
// 				])
// 			);
// 		`, &stats),
// 		chromedp.Evaluate(`
// 			Object.fromEntries(
// 				Array.from(document.querySelectorAll('.when-attacked-row')).map(el => {
// 					let type = el.querySelector('.monster-type').innerText.trim();
// 					let multiplier = parseFloat(el.querySelector('.monster-multiplier').innerText.trim().replace('x', ''));
// 					return [type, multiplier];
// 				})
// 			);
// 		`, &damageMultipliers),
// 	}

// 	err := chromedp.Run(ctx, tasks)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Pokemon{
// 		ID:                id,
// 		Name:              name,
// 		Types:             types,
// 		Stats:             stats,
// 		DamageMultipliers: damageMultipliers,
// 	}, nil
// }
