package internals

import (
	"fmt"
	"math/rand"
)

var pokemons_caught = make(map[string]Pokemon)

func catchPokemon(pokemon Pokemon) (caught bool) {
	guess := rand.Intn(pokemon.BaseExperience)

	if guess == 0 {
		guess = rand.Intn(pokemon.BaseExperience)
	}

	fmt.Println("catch attempt data", guess, pokemon.BaseExperience, pokemon.BaseExperience%guess)

	return pokemon.BaseExperience%guess == 0
}

func saveCatches(pokemon Pokemon) {
	_, ok := pokemons_caught[pokemon.Name]
	if !ok {
		pokemons_caught[pokemon.Name] = pokemon
	}
}
