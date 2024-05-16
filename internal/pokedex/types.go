package pokedex

type Pokemon struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	BaseExp int    `json:"base_exp"`
}
