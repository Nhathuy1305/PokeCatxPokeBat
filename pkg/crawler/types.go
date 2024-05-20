package crawler

type Pokemon struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	Types             []string           `json:"types"`
	Stats             map[string]int     `json:"stats"`
	DamageMultipliers map[string]float64 `json:"damage_multipliers"`
}
