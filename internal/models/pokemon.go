package models

type Pokemon struct {
	Name           string
	BaseExperience int
	ID             int
	Stats          []PokemonStat
	Types          []string
	Weight         int
	Height         int
}

type PokemonStat struct {
	Name     string
	BaseStat int
}

type PokemonShortInfo struct {
	Name string
	Url  string
}
