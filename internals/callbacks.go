package internals

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ekediala/pokedexcli/internals/cache"
)

func (p names) PrintName() {
	fmt.Println(p.Name)
}

func (p result) PrintName() {
	fmt.Println(p.Name)
}

func (p Stat) PrintName() {
	fmt.Println(p.Stat.Name, p.BaseStat)
}

func processResults[T ResultSet](data []T) error {
	for _, element := range data {
		element.PrintName()
	}
	return nil
}
func CallbackMap(cfg *Config, cache *cache.Cache, location string) error {
	cfg.RLock()
	defer cfg.RUnlock()

	url := cfg.Next

	from_cache, ok := cache.Get(url)

	var data pokemon_location

	if ok {
		if err := json.Unmarshal(from_cache, &data); err != nil {
			return err
		}
		return processResults(data.Results)

	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}

	cfg.Next = data.Next
	cfg.Prev = url

	json_data, err := json.Marshal(data)

	if err != nil {
		return err
	}

	cache.Add(url, json_data)

	return processResults(data.Results)
}

func CallbackMapB(cfg *Config, cache *cache.Cache, location string) error {

	if cfg.Prev == "" {
		return errors.New("no previous page")
	}

	cfg.RLock()
	defer cfg.RUnlock()

	url := cfg.Prev

	from_cache, ok := cache.Get(url)

	var data pokemon_location

	if ok {
		if err := json.Unmarshal(from_cache, &data); err != nil {
			return err
		}
		cfg.Next = url
		return processResults(data.Results)
	}

	res, err := http.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err

	}

	previous, _ := data.Previous.(string)
	cfg.Next = url
	cfg.Prev = previous

	json_data, err := json.Marshal(data.Results)

	if err != nil {
		return err
	}

	cache.Add(url, json_data)

	return processResults(data.Results)
}

func CallbackInspectPokemon(cfg *Config, cache *cache.Cache, pokemon_name string) error {
	pokemon, ok := pokemons_caught[pokemon_name]

	if !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	processResults(pokemon.Stats)
	return nil
}

func CallbackPokedex(cfg *Config, cache *cache.Cache, pokemon_name string) error {
	if len(pokemons_caught) == 0 {
		fmt.Println("Your pokedex is empty")
		return nil
	}

	fmt.Println("Your pokedex:")

	for _, value := range pokemons_caught {
		fmt.Println(" - ", value.Name)
	}

	return nil
}

func CallbackCatchPokemon(cfg *Config, cache *cache.Cache, arg string) error {

	fmt.Printf("Throwing a pokeball at %s\n", arg)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", arg)

	poke, ok := cache.Get(arg)

	if ok {
		var pokemon Pokemon
		if err := json.Unmarshal(poke, &pokemon); err != nil {
			return err
		}

		caughtPokemon := catchPokemon(pokemon)

		if caughtPokemon {
			fmt.Printf("%s was caught!\n", arg)
			saveCatches(pokemon)
			return nil
		}

		fmt.Printf("%s escaped!\n", arg)
		return nil
	}

	res, err := http.Get(url)

	if err != nil {
		return err
	}

	var pokemon Pokemon

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&pokemon); err != nil {
		return err
	}

	json_data, err := json.Marshal(pokemon)

	if err != nil {
		return err
	}

	cache.Add(arg, json_data)

	caughtPokemon := catchPokemon(pokemon)

	if caughtPokemon {
		fmt.Printf("%s was caught!\n", arg)
		saveCatches(pokemon)
		return nil
	}

	fmt.Printf("%s escaped!\n", arg)
	return nil
}

func CallBackExplore(cfg *Config, cache *cache.Cache, location string) error {
	cfg.RLock()
	defer cfg.RUnlock()

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)

	fmt.Printf("Exploring %s...\n", location)

	from_cache, ok := cache.Get(url)

	var data PokemonByLocationNameResponse

	if ok {
		if err := json.Unmarshal(from_cache, &data); err != nil {
			return err
		}
		return processResults(data.Names)
	}

	res, err := http.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}

	json_data, err := json.Marshal(data)

	if err != nil {
		return err
	}

	cache.Add(url, json_data)

	return processResults(data.Names)
}

func CallbackHelp(cfg *Config, cache *cache.Cache, location string) error {
	fmt.Println("\nWelcome to the Pokedex!\n\nUsage: \nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func CallbackExit(cfg *Config, cache *cache.Cache, location string) error {
	cache.Clear()
	os.Exit(1)
	return nil
}
