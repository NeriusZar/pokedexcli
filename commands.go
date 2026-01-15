package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/NeriusZar/pokedexcli/internal/api"
	"github.com/NeriusZar/pokedexcli/internal/models"
	"github.com/NeriusZar/pokedexcli/internal/pokedex"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	pagination models.Pagination
	api        api.PokeApi
	pokedex    pokedex.Pokedex
}

func NewConfig() config {
	return config{
		pagination: models.Pagination{},
		api:        api.NewPokeApi(),
		pokedex:    pokedex.NewPokedex(),
	}
}

const difficultyConf = 40

func getCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next page of the names of 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous page of the names of 20 location areas in the Pokemon world.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Takes location argument and lists all the Pokemons located there.",
			callback:    explore,
		},
		"catch": {
			name:        "catch",
			description: "Takes name of pokemon and attempts to catch the pokemon.",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows details of caught Pokemon",
			callback:    inspect,
		},
	}
}

func commandExit(c *config, a ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, a ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(c *config, a ...string) error {
	areas, pagination, err := c.api.RetrieveAreas(c.pagination.Next)
	if err != nil {
		return err
	}

	c.pagination.Next = pagination.Next
	c.pagination.Previous = pagination.Previous

	for _, a := range areas {
		fmt.Println(a.Name)
	}
	return nil
}

func commandMapb(c *config, a ...string) error {
	if c.pagination.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	areas, pagination, err := c.api.RetrieveAreas(c.pagination.Previous)
	if err != nil {
		return err
	}

	c.pagination.Next = pagination.Next
	c.pagination.Previous = pagination.Previous

	for _, a := range areas {
		fmt.Println(a.Name)
	}
	return nil
}

func explore(c *config, a ...string) error {
	if len(a) < 1 {
		fmt.Println("You didn't provide any areas to explore.")
		return nil
	}
	area := a[0]

	fmt.Printf("Exploring %s...\n", area)

	pokemons, err := c.api.RetrievePokemonsInArea(area)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, p := range pokemons {
		fmt.Printf(" - %s\n", p.Name)
	}

	return nil
}

func catch(c *config, a ...string) error {
	if len(a) < 1 {
		fmt.Println("You didn't provide any pokemon to catch")
		return nil
	}

	name := a[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := c.api.GetPokemonDetails(name)
	if err != nil {
		return err
	}

	if rand.IntN(pokemon.BaseExperience) > difficultyConf {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}

	c.pokedex.Add(pokemon)

	fmt.Printf("%s was caught!\n", name)

	return nil
}

func inspect(c *config, a ...string) error {
	if len(a) < 1 {
		fmt.Println("You didn't provide pokemon name")
		return nil
	}

	name := a[0]

	pokemon, ok := c.pokedex.Get(name)
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Height: %d\n", pokemon.Height)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t)
	}

	return nil
}
