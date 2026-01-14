package models

type Pokemon struct {
	Name           string
	BaseExperience int
	ID             int
}

type PokemonShortInfo struct {
	Name string
	Url  string
}
