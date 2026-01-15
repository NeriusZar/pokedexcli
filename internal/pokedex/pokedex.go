package pokedex

import (
	"sync"

	"github.com/NeriusZar/pokedexcli/internal/models"
)

type Pokedex struct {
	pokemons map[string]models.Pokemon
	mu       *sync.Mutex
}

func NewPokedex() Pokedex {
	return Pokedex{
		pokemons: map[string]models.Pokemon{},
		mu:       &sync.Mutex{},
	}
}

func (p Pokedex) Add(pokemon models.Pokemon) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pokemons[pokemon.Name] = pokemon
}

func (p Pokedex) Get(name string) (models.Pokemon, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pokemon, ok := p.pokemons[name]
	return pokemon, ok
}

func (p Pokedex) GetAll() []models.Pokemon {
	p.mu.Lock()
	defer p.mu.Unlock()

	pokemons := make([]models.Pokemon, 0, len(p.pokemons))
	for _, v := range p.pokemons {
		pokemons = append(pokemons, v)
	}

	return pokemons
}
