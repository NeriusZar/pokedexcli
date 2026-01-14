package pokedex

import (
	"sync"

	"github.com/NeriusZar/pokedexcli/internal/models"
)

type Pokedex struct {
	Pokemons map[string]models.Pokemon
	mu       *sync.Mutex
}

func NewPokedex() Pokedex {
	return Pokedex{
		Pokemons: map[string]models.Pokemon{},
		mu:       &sync.Mutex{},
	}
}

func (p Pokedex) Add(pokemon models.Pokemon) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Pokemons[pokemon.Name] = pokemon
}
