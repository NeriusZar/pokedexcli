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
	callback    func(*models.Pagination) (error)
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

func commandExit(p *models.Pagination) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(p *models.Pagination) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(p *models.Pagination) error {
	areas, pagination, err := api.RetrieveAreas(p.Next)
	if err != nil {
		return err
	}

	p.Next = pagination.Next
	p.Previous = pagination.Previous

	for _, a := range areas {
		fmt.Println(a.Name)
	}
	return nil
}

func commandMapb(p *models.Pagination) error {
	if p.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	areas, pagination, err := api.RetrieveAreas(*p.Previous)
	if err != nil {
		return err
	}

	p.Next = pagination.Next
	p.Previous = pagination.Previous

	for _, a := range areas {
		fmt.Println(a.Name)
	}
	return nil
}
