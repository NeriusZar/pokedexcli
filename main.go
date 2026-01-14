package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := NewConfig()
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

			err := command.callback(&config, cleanedInput[1:]...)
			if err != nil {
				fmt.Println("Failed to execute command", err)
			}
		}
	}
}
