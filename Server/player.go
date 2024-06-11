package main
import (
	"os"
	"encoding/json"
	"fmt"
	"strings"
)

func createPlayer(pokedex []Pokemon, playerName string) Player {
	// Create the player

	player := Player{
		Name:        playerName,
		PokemonList: []Pokemon{},
	}
	// Choose 3 starter Pokemon
	for _, p := range starters {
		pokemon, _ := findPokemon(pokedex, p)
		player.PokemonList = append(player.PokemonList, pokemon)
	}
	// Load the existing players from the JSON file
	file, _ := os.Open("./crawler/players.json")
	decoder := json.NewDecoder(file)
	existingPlayers := []Player{}
	_ = decoder.Decode(&existingPlayers)
	file.Close()

	// Append the new player to the list
	existingPlayers = append(existingPlayers, player)

	// Save the updated list of players to the JSON file
	file, _ = os.Create("./crawler/players.json")
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Set indent to 4 spaces
	_ = encoder.Encode(existingPlayers)

	fmt.Printf("Player %s created\n", playerName)
	return player
}

func findPlayer(name string) (Player, bool) {
	// Load the Pokedex
	file, _ := os.Open("./crawler/players.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	players := []Player{}
	_ = decoder.Decode(&players)
	for _, p := range players {
		if strings.EqualFold(p.Name, name) {
			return p, true
		}
	}
	return Player{}, false
}