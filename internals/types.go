package internals

import "sync"

type Config struct {
	Next string
	Prev string
	sync.RWMutex
}

type ResultSet interface {
	PrintName()
}

type pokemon_location struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous any      `json:"previous"`
	Results  []result `json:"results"`
}

type result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type names struct {
	Language result `json:"language"`
	Name     string `json:"name"`
}

type PokemonByLocationNameResponse struct {
	ID       int     `json:"id"`
	Location result  `json:"location"`
	Name     string  `json:"name"`
	Names    []names `json:"names"`
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Weight         int    `json:"weight"`
	Stats          []Stat  `json:"stats"`
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}
