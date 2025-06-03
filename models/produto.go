package models

type Produto struct {
	ID      int     `json:"id"`
	Nome    string  `json:"nome"`
	Preco   float64 `json:"preco"`
	Estoque int     `json:"estoque"`
}
