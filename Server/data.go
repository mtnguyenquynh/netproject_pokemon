

package main

import (
    "encoding/json"
    "fmt"
    "os"

    rl "github.com/gen2brain/raylib-go/raylib"


)

var WinMessages = []string{"That was excellent. Truly, an outstanding battle. You gave the support your Pokémon needed to maximize their power. And you guided them with certainty to secure victory. You have both passion and calculating coolness. Together, you and your Pokémon can overcome any challenge that may come your way. Those are the impressions I got from our battle. I'm glad I got to take part in the crowning of Sinnoh's new Champion! Come with me. We'll take the lift."}   // win messages
var LoseMessages = []string{"Smell ya later!", "Better luck next time", "Keep training", "Time to soft-reset", "You whited out...", "Come back when you're stronger", "Do or do not...there is no try", "Looks like you're blasting off again"}   // lose messages




// List of Pokemon types
var Types = map[string]int{
    "Normal": 0,
    "Fire": 1,
    "Water": 2,
    "Electric": 3,
    "Grass": 4,
    "Ice": 5,
    "Fighting": 6,
    "Poison": 7,
    "Ground": 8,
    "Flying": 9,
    "Psychic": 10,
    "Bug": 11,
    "Rock": 12,
    "Ghost": 13,
    "Dragon": 14,
    "Dark": 15,
    "Steel": 16,
}

// Type multipliers based on type
// KEY is attacking type
// VALUES are multiplier based on defending type
var Matchup = map[string][17]float64{
    "Normal":   {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0.5, 0, 1, 1, 0.5},
    "Fire":     {1, 0.5, 0.5, 1, 2, 2, 1, 1, 1, 1, 1, 2, 0.5, 1, 0.5, 1, 1},
    "Water":    {1, 2, 0.5, 1, 0.5, 1, 1, 1, 2, 1, 1, 1, 2, 1, 0.5, 1, 1},
    "Electric": {1, 1, 2, 0.5, 0.5, 1, 1, 1, 0, 2, 1, 1, 1, 1, 0.5, 1, 1},
    "Grass":    {1, 0.5, 2, 1, 0.5, 1, 1, 0.5, 2, 0.5, 1, 0.5, 2, 1, 0.5, 1, 0.5},
    "Ice":      {1, 0.5, 0.5, 1, 2, 0.5, 1, 1, 2, 2, 1, 1, 1, 1, 2, 1, 0.5},
    "Fighting": {2, 1, 1, 1, 1, 2, 1, 0.5, 1, 0.5, 0.5, 0.5, 2, 0, 1, 2, 2},
    "Poison":   {1, 1, 1, 1, 2, 1, 1, 0.5, 0.5, 1, 1, 1, 0.5, 0.5, 1, 1, 0},
    "Ground":   {1, 2, 1, 2, 0.5, 1, 1, 2, 1, 0, 1, 0.5, 2, 1, 1, 1, 2},
    "Flying":   {1, 1, 1, 0.5, 2, 1, 2, 1, 1, 1, 1, 2, 0.5, 1, 1, 1, 0.5},
    "Psychic":  {1, 1, 1, 1, 1, 1, 2, 2, 1, 1, 0.5, 1, 1, 1, 1, 0, 0.5},
    "Bug":      {1, 0.5, 1, 1, 2, 1, 0.5, 0.5, 1, 0.5, 2, 1, 1, 0.5, 1, 2, 0.5},
    "Rock":     {1, 2, 1, 1, 1, 2, 0.5, 1, 0.5, 2, 1, 2, 1, 1, 1, 1, 0.5},
    "Ghost":    {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 2, 1, 0.5, 0.5},
    "Dragon":   {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 0.5},
    "Dark":     {1, 1, 1, 1, 1, 1, 0.5, 1, 1, 1, 2, 1, 1, 2, 1, 0.5, 0.5},
    "Steel":    {1, 0.5, 0.5, 0.5, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 0.5},
}



// Define structs to hold the data from pokedex.json
type PokemonData struct {
    PokedexNumber  string   `json:"index"`
    Name           string   `json:"name"`
    BaseHP         int      `json:"hp"`
    BaseAtk        int      `json:"attack"`
    BaseDef        int      `json:"defense"`
    BaseSpAtk      int      `json:"sp_attack"`
    BaseSpDef      int      `json:"sp_defense"`
    BaseSpeed      int      `json:"speed"`
    TotalEVs       int      `json:"total_evs"`
    Type           [2]string `json:"type"`
    Description    string   `json:"description"`
    Height         string   `json:"height"`
    Weight         string   `json:"weight"`
    ImageURL    string       `json:"image_url"`
    Exp            int      `json:"exp"`
    Moves          []Move   `json:"moves"`
    Texture     rl.Texture2D `json:"-"`
	Position    rl.Vector2   `json:"-"`
}

// 	ImageURL    string       `json:"image_url"`




type Pokemon struct {
    // Shared among individuals
    PokemonData

    // Specific per individual
    level             int
    hp                int
    atk               int
    def               int
    spatk             int
    spdef             int
    speed             int
    nonVolatileStatus string
    volatileStatus    string
    fainted           bool
}

// Move represents the structure of a move.
type Move struct {
    MoveName            string  `json:"name"`
    MoveType        string  `json:"type"`
    AtkType         string  `json:"atk_type"`
    Power           int     `json:"power"`
    Accuracy        int     `json:"accuracy"`
    SecondEffectRate float64 `json:"pp"`
    SecondEffect    string  `json:"description"`
}

var Moves map[string]Move

var pokemonList map[string]PokemonData



func InitData()  {
    pokemonList = make(map[string]PokemonData)
	// Read data from pokedex.json and populate PokemonList
	file, err := os.Open("./crawler/pokedex.json")
	if err != nil {
		// return nil, err
        fmt.Println(err)
        os.Exit(1)
	}
	defer file.Close()


	decoder := json.NewDecoder(file)
	pokedex := []PokemonData{}
	if err := decoder.Decode(&pokedex); err != nil {
		// return nil, err
        fmt.Println(err)
        os.Exit(1)
	}

	// pokemonList := make(map[string]PokemonData)
	// Iterate over the decoded data and populate PokemonList
	for _, pokemon := range pokedex {
		pokemonList[pokemon.Name] = pokemon
	}

	fmt.Println("PokemonList loaded")

	// Debugging: Print out decoded PokemonList data
	fmt.Println("Decoded PokemonList data:")
	for _, pokemon := range pokemonList {
		fmt.Println(pokemon.Name)
	}

	// return pokemonList, nil
}

