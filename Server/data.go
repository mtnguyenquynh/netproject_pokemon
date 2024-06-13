

package main

import (
    "encoding/json"
    "fmt"
    "os"
)

var WinMessages = []string{"That was excellent. Truly, an outstanding battle. You gave the support your Pokémon needed to maximize their power. And you guided them with certainty to secure victory. You have both passion and calculating coolness. Together, you and your Pokémon can overcome any challenge that may come your way. Those are the impressions I got from our battle. I'm glad I got to take part in the crowning of Sinnoh's new Champion! Come with me. We'll take the lift."}   // win messages
var LoseMessages = []string{"Smell ya later!", "Better luck next time", "Keep training", "Time to soft-reset", "You whited out...", "Come back when you're stronger", "Do or do not...there is no try", "Looks like you're blasting off again"}   // lose messages

type UserInput struct {
    username string
    action string
    activePokemon *Pokemon
    team []*Pokemon
    move string
    isAI bool
    gameOver bool
}

// List of non-volatile status conditions
var StatusList = map[string]bool{
    "PSN": true, // Poison
    "FRZ": true, // Freeze  
    "BRN": true, // Burn
    "PRZ": true, // Paralysis
}

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

var MoveList = map[string]Move{
    "Strength": Move{"Strength", "Normal", "atk", 80, 100, 0, "None"},
    "Hyper Voice": Move{"Hyper Voice", "Normal", "spatk", 90, 100, 0, "None"},
    "Flamethrower": Move{"Flamethrower", "Fire", "spatk", 95, 100, 0.1, "BRN"},
    "Flame Wheel": Move{"Flame Wheel", "Fire", "atk", 60, 100, 0.1, "BRN"},
    "Surf": Move{"Surf", "Water", "spatk", 95, 100, 0, "None"},
    "Waterfall": Move{"Waterfall", "Water", "atk", 80, 100, 0.2, "flinch"},
    "Thunderbolt": Move{"Thunderbolt", "Electric", "spatk", 95, 100, 0.1, "PRZ"},
    "Spark": Move{"Spark", "Electric", "atk", 65, 100, 0.3, "PRZ"},
    "Energy Ball": Move{"Energy Ball", "Grass", "spatk", 80, 100, 0.1, "spdef"},
    "Seed Bomb": Move{"Seed Bomb", "Grass", "atk", 80, 100, 0, "None"},
    "Leaf Blade": Move{"Leaf Blade", "Grass", "atk", 90, 100, 2, "crit"},
    "Ice Beam": Move{"Ice Beam", "Ice", "spatk", 95, 100, 0.1, "FRZ"},
    "Avalanche": Move{"Avalanche", "Ice", "atk", 80, 100, 0, "None"},
    "Aura Sphere": Move{"Aura Sphere", "Fighting", "spatk", 90, 100, 0, "None"},
    "Brick Break": Move{"Brick Break", "Fighting", "atk", 75, 100, 0, "None"},
    "Sludge Bomb": Move{"Sludge Bomb", "Poison", "spatk", 90, 100, 0.3, "PSN"},
    "Poison Jab": Move{"Poison Jab", "Poison", "atk", 90, 100, 0.3, "PSN"},
    "Earthquake": Move{"Earthquake", "Ground", "atk", 100, 100, 0, "None"},
    "Earth Power": Move{"Earth Power", "Ground", "spatk", 90, 100, 10, "spdef"},
    "Air Slash": Move{"Air Slash", "Flying", "spatk", 75, 95, 0.3, "flinch"},
    "Aerial Ace": Move{"Aerial Ace", "Flying", "atk", 60, 100, 0, "None"},
    "Psychic": Move{"Psychic", "Psychic", "spatk", 90, 100, 0.1, "spdef"},
    "Psycho Cut": Move{"Psycho Cut", "Psychic", "atk", 90, 100, 0, "None"},
    "Extrasensory": Move{"Extrasensory", "Psychic", "spatk", 80, 100, 0.1, "flinch"},
    "X-scissor": Move{"X-scissor", "Bug", "atk", 80, 100, 0, "None"},
    "Bug Buzz": Move{"Bug Buzz", "Bug", "spatk", 90, 100, 0.1, "spdef"},
    "Rock Slide": Move{"Rock Slide", "Rock", "atk", 75, 90, 0.3, "flinch"},
    "Stone Edge": Move{"Stone Edge", "Rock", "atk", 100, 80, 2, "crit"},
    "Power Gem": Move{"Power Gem", "Rock", "spatk", 80, 100, 0, "None"},
    "Shadow Ball": Move{"Shadow Ball", "Ghost", "spatk", 80, 100, 0.2, "spdef"},
    "Shadow Punch": Move{"Shadow Punch", "Ghost", "atk", 60, 100, 0, "None"},
    "Dragon Pulse": Move{"Dragon Pulse", "Dragon", "spatk", 90, 100, 0, "None"},
    "Dragon Claw": Move{"Dragon Claw", "Dragon", "atk", 80, 100, 0, "None"},
    "Dark Pulse": Move{"Dark Pulse", "Dark", "spatk", 80, 100, 0.2, "flinch"},
    "Night Slash": Move{"Night Slash", "Dark", "atk", 70, 100, 2, "crit"},
    "Crunch": Move{"Crunch", "Dark", "atk", 80, 100, 0.2, "def"},
    "Iron Head": Move{"Iron Head", "Steel", "atk", 80, 100, 0.3, "flinch"},
    "Flash Cannon": Move{"Flash Cannon", "Steel", "spatk", 80, 100, 0.1, "spdef"},
}

// Define structs to hold the data from pokedex.json
type PokemonData struct {
    PokedexNumber  string   `json:"index"`
    Name           string   `json:"name"`
    Exp            int      `json:"exp"`
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
    Level          int      `json:"level"`
    AccumExp       int      `json:"accum_exp"`
    Moves         []Move  `json:"moves"`
}

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

