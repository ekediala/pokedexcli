package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ekediala/pokedexcli/internals"
	"github.com/ekediala/pokedexcli/internals/cache"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	cfg := &internals.Config{
		Next: "https://pokeapi.co/api/v2/location-area",
		Prev: "",
	}

	c := cache.New(time.Minute * 5)

	go c.ReapLoop()

	type cliCommand struct {
		name        string
		description string
		callback    func(cfg *internals.Config, cache *cache.Cache, arg string) error
	}

	cliCommands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    internals.CallbackHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    internals.CallbackExit,
		},
		"map": {
			name:        "map",
			description: "Get the next 20 items",
			callback:    internals.CallbackMap,
		},
		"mapb": {
			name:        "map back",
			description: "get the previous 20 items",
			callback:    internals.CallbackMapB,
		},
		"explore": {
			name:        "explore",
			description: "get the pokemons in current location",
			callback:    internals.CallBackExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch a given pokemon",
			callback:    internals.CallbackCatchPokemon,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect a given pokemon",
			callback:    internals.CallbackInspectPokemon,
		},
		"pokedex": {
			name:        "pokedex",
			description: "list pokemons caught",
			callback:    internals.CallbackPokedex,
		},
	}

	for {
		fmt.Print("\npokedex> ")
		text, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading input", err)
		}

		text = text[:len(text)-1] // remove newline character

		var arg string

		if strings.Contains(text, "explore") {
			arguments := strings.Split(text, " ")

			if len(arguments) < 2 {
				fmt.Println("Please provide pokemon location", text)
				continue
			}

			arg = arguments[1]
			text = arguments[0]
		}

		if strings.Contains(text, "catch") {
			arguments := strings.Split(text, " ")

			if len(arguments) < 2 {
				fmt.Println("Please provide pokemon name", text)
				continue
			}

			arg = arguments[1]
			text = arguments[0]
		}

		if strings.Contains(text, "inspect") {
			arguments := strings.Split(text, " ")

			if len(arguments) < 2 {
				fmt.Println("Please provide pokemon name", text)
				continue
			}

			arg = arguments[1]
			text = arguments[0]
		}

		command, ok := cliCommands[text]

		if !ok {
			fmt.Println("Command not found", text)
			continue
		}

		err = command.callback(cfg, c, arg)

		if err != nil {
			fmt.Println("Error executing command: ", err)
		}

	}
}
