package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "math/rand"
    "net"
)

type UserInput struct {
    username      string
    action        string
    activePokemon *Pokemon
    team          []*Pokemon
    move          string
    isAI          bool
    gameOver      bool
}

var pokemonList map[string]PokemonData

type Pokemon struct {
    PokemonData
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
    "Growl":        {"Growl", "Normal", "atk", 0, 100, 4, "None"},
    "Tackle":       {"Tackle", "Normal", "atk", 50, 100, 3.5, "None"},
    "Leech Seed":   {"Leech Seed", "Grass", "spatk", 0, 90, 1, "None"},
    "Scratch":      {"Scratch", "Normal", "atk", 40, 100, 3.5, "None"},
    "Ember":        {"Ember", "Fire", "spatk", 40, 100, 2.5, "None"},
    "Shadow Claw":  {"Shadow Claw", "Ghost", "spatk", 70, 100, 1.5, "None"},
    "Air Slash":    {"Air Slash", "Flying", "spatk", 75, 95, 1.5, "None"},
    "Heat Wave":    {"Heat Wave", "Fire", "spatk", 95, 90, 1, "None"},
    "Tail Whip":    {"Tail Whip", "Normal", "atk", 0, 100, 3, "None"},
    "Bubble":       {"Bubble", "Water", "spatk", 40, 100, 3, "None"},
    "Water Gun":    {"Water Gun", "Water", "spatk", 40, 100, 2.5, "None"},
    "Flash Cannon": {"Flash Cannon", "Steel", "spatk", 80, 100, 1, "None"},
    "String Shot":  {"String Shot", "Bug", "spatk", 0, 95, 4, "None"},
    "Bug Bite":     {"Bug Bite", "Bug", "spatk", 60, 100, 2, "None"},
}


// Move represents the structure of a move.


var Moves map[string]Move

func battle(name string) {
    conn, err := net.Dial("tcp", "localhost:8081")
    if err != nil {
        fmt.Println("Error connecting to server:", err.Error())
        return
    }
    defer conn.Close()

    reader := bufio.NewReader(conn)

    // Read welcome message
    welcomeMessage, _ := reader.ReadString('\n')
    fmt.Print(welcomeMessage)

    // Read Pokémon list
    pokemonListJSON, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading Pokémon list:", err)
        return
    }

    // Parse JSON data
    err = json.Unmarshal([]byte(pokemonListJSON), &pokemonList)
    if err != nil {
        fmt.Println("Error unmarshaling Pokémon list:", err)
        return
    }

    // Example team setup using received Pokémon list
    venusaur := NewPokemon("Venusaur", true)
    charmeleon := NewPokemon("Charmeleon", true)
    wartortle := NewPokemon("Wartortle", true)
    blastoise := NewPokemon("Blastoise", true)
    caterpie := NewPokemon("Caterpie", true)
    bulbasaur := NewPokemon("Bulbasaur", true)

    cynthiasTeam := []*Pokemon{venusaur, charmeleon, wartortle, blastoise, caterpie, bulbasaur}

    myInput := &UserInput{"Ash", "", nil, nil, "", false, false}
    cynthiasInput := &UserInput{"Cynthia", "", charmeleon, cynthiasTeam, "", true, false}

    fmt.Println()
    ChooseName(myInput, name)
    ChooseTeam(myInput)
    myInput.activePokemon = myInput.team[0]

    Battle(myInput, cynthiasInput)
}


func ChooseName(input *UserInput, mv string) *UserInput {
	input.username = mv
	return input
}

// calculateHp calculates a Pokemon's HP stat
func calculateHp(baseHp, level, iv, ev int) int {
    return (((2 * baseHp) + iv + (ev / 4)) * level) / 100 + level + 10
}

// calculateOtherStat calculates a Pokemon's other stats (Attack, Defense, etc.)
func calculateOtherStat(baseStat, level, iv, ev int) int {
    return (((2 * baseStat) + iv + (ev / 4)) * level) / 100 + 5
}

// initializeStats initializes the actual stats of a Pokemon based on its base stats, IVs, etc.
func initializeStats(pokemon *Pokemon, makeStrong bool) {
    if makeStrong { // Cynthia's Pokemon will be stronger
        pokemon.level = 60
        IV := 31
        EV := 252
        pokemon.hp = calculateHp(pokemon.BaseHP, pokemon.level, IV, EV)
        pokemon.atk = calculateOtherStat(pokemon.BaseAtk, pokemon.level, IV, EV)
        pokemon.def = calculateOtherStat(pokemon.BaseDef, pokemon.level, IV, EV)
        pokemon.spatk = calculateOtherStat(pokemon.BaseSpAtk, pokemon.level, IV, EV)
        pokemon.spdef = calculateOtherStat(pokemon.BaseSpDef, pokemon.level, IV, EV)
        pokemon.speed = calculateOtherStat(pokemon.BaseSpeed, pokemon.level, IV, EV)
    } else { // Player's Pokemon will be slightly weaker
        pokemon.level = rand.Intn(10) + 50
        IV := rand.Intn(20) + 10
        EV := rand.Intn(100) + 120
        pokemon.hp = calculateHp(pokemon.BaseHP, pokemon.level, IV, EV)
        pokemon.atk = calculateOtherStat(pokemon.BaseAtk, pokemon.level, IV, EV)
        pokemon.def = calculateOtherStat(pokemon.BaseDef, pokemon.level, IV, EV)
        pokemon.spatk = calculateOtherStat(pokemon.BaseSpAtk, pokemon.level, IV, EV)
        pokemon.spdef = calculateOtherStat(pokemon.BaseSpDef, pokemon.level, IV, EV)
        pokemon.speed = calculateOtherStat(pokemon.BaseSpeed, pokemon.level, IV, EV)
    }
}

// NewPokemon creates a new Pokemon with initialized stats and returns a pointer to it
func NewPokemon(name string, makeStrong bool) *Pokemon {
    template := pokemonList[name]

    pokemon := Pokemon{
        PokemonData: template,
    }

    // Initialize stats
    initializeStats(&pokemon, makeStrong)

    // Initialize battle-specific fields
    pokemon.nonVolatileStatus = ""
    pokemon.volatileStatus = ""
    pokemon.fainted = false

    return &pokemon
}


