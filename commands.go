package main

import (
	"fmt"
	"os"

	"github.com/NeriusZar/pokedexcli/internal/api"
	"github.com/NeriusZar/pokedexcli/internal/models"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	pagination models.Pagination
	api        api.PokeApi
}

func NewConfig() config {
	return config{
		pagination: models.Pagination{},
		api:        api.NewPokeApi(),
	}
}

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
	}
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(c *config) error {
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

func commandMapb(c *config) error {
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
