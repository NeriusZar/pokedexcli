package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/NeriusZar/pokedexcli/internal/models"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pagination := models.Pagination{}
	commands := getCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		if len(cleanedInput) > 0 {
			command, ok := commands[cleanedInput[0]]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}

			err := command.callback(&pagination)
			if err != nil {
				fmt.Println("Failed to execute command")
			}
		}
	}
}
