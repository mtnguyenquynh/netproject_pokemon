/* package main

// list of non-volatile status conditions
var StatusList = map[string]bool{
    "PSN": true,  //Poison
    "FRZ": true,  //Freeze
    "BRN": true,  //Burn
    "PRZ": true,  //Paralysis
}


// List of pokemon types
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


// type multipliers based on type
// KEY is attacking type
// VALUES are multiplier based on defending type
var Matchup = map[string][17]float64{
    "Normal": [17]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0.5, 0, 1, 1, 0.5}, 
    "Fire": [17]float64{1, 0.5, 0.5, 1, 2, 2, 1, 1, 1, 1, 1, 2, 0.5, 1, 0.5, 1, 1},
    "Water": [17]float64{1, 2, 0.5, 1, 0.5, 1, 1, 1, 2, 1, 1, 1, 2, 1, 0.5, 1, 1},
    "Electric": [17]float64{1, 1, 2, 0.5, 0.5, 1, 1, 1, 0, 2, 1, 1, 1, 1, 0.5, 1, 1},
    "Grass": [17]float64{1, 0.5, 2, 1, 0.5, 1, 1, 0.5, 2, 0.5, 1, 0.5, 2, 1, 0.5, 1, 0.5},
    "Ice": [17]float64{1, 0.5, 0.5, 1, 2, 0.5, 1, 1, 2, 2, 1, 1, 1, 1, 2, 1, 0.5},
    "Fighting": [17]float64{2, 1, 1, 1, 1, 2, 1, 0.5, 1, 0.5, 0.5, 0.5, 2, 0, 1, 2, 2},
    "Poison": [17]float64{1, 1, 1, 1, 2, 1, 1, 0.5, 0.5, 1, 1, 1, 0.5, 0.5, 1, 1, 0},
    "Ground": [17]float64{1, 2, 1, 2, 0.5, 1, 1, 2, 1, 0, 1, 0.5, 2, 1, 1, 1, 2},
    "Flying": [17]float64{1, 1, 1, 0.5, 2, 1, 2, 1, 1, 1, 1, 2, 0.5, 1, 1, 1, 0.5},
    "Psychic": [17]float64{1, 1, 1, 1, 1, 1, 2, 2, 1, 1, 0.5, 1, 1, 1, 1, 0, 0.5},
    "Bug": [17]float64{1, 0.5, 1, 1, 2, 1, 0.5, 0.5, 1, 0.5, 2, 1, 1, 0.5, 1, 2, 0.5},
    "Rock": [17]float64{1, 2, 1, 1, 1, 2, 0.5, 1, 0.5, 2, 1, 2, 1, 1, 1, 1, 0.5},
    "Ghost": [17]float64{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 2, 1, 0.5, 0.5},
    "Dragon": [17]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 0.5},
    "Dark": [17]float64{1, 1, 1, 1, 1, 1, 0.5, 1, 1, 1, 2, 1, 1, 2, 1, 0.5, 0.5},
    "Steel": [17]float64{1, 0.5, 0.5, 0.5, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 0.5},
}

/* // move, type, special/physical, power, accuracy, secondary effect
type Move struct {
    name string
    moveType string
    atkType string
    power int
    accuracy int
    secondEffectRate float64
    secondEffect string
} 

// Move represents the structure of a move.
type Move struct {
	name        string  `json:"name"`
	moveType        string  `json:"type"`
	atkType     string  `json:"atk_type"`
	power       int     `json:"power"`
	accuracy    float64 `json:"accuracy"`
	secondEffectRate          float64     `json:"pp"`
	description string  `json:"description"`
}
// maps move name to information
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

type PokemonData struct {
    // shared among individuals
    name string
    pokedexNumber int
    type1 string
    type2 string
    baseHp int
    baseAtk int
    baseDef int
    baseSpatk int
    baseSpdef int
    baseSpeed int
    moves [4]string
} 

type PokemonData struct {
	pokedexNumber       string   `json:"index"`
	name        string   `json:"name"`
	Exp         int      `json:"exp"`
	baseHp          int      `json:"hp"`
	baseAtk      int      `json:"attack"`
	baseDef     int      `json:"defense"`
	baseSpatk    int      `json:"sp_attack"`
	baseSpdef   int      `json:"sp_defense"`
	baseSpeed       int      `json:"speed"`
	TotalEVs    int      `json:"total_evs"`
	Type        [2]string `json:"type"`
	Description string   `json:"description"`
	Height      string   `json:"height"`
	Weight      string   `json:"weight"`
	Level       int      `json:"level"`
	AccumExp    int      `json:"accum_exp"`
    moves [4]string

}

type Pokemon struct {
    // shared among individuals
    PokemonData

    // specific per individual
    level int
    hp int
    atk int
    def int
    spatk int
    spdef int
    speed int
    nonVolatileStatus string
    volatileStatus string
    fainted bool
}

var PokemonList = map[string]PokemonData{
    "Spiritomb": PokemonData{"108", "Spiritomb", 108, 50, 92, 108, 92, 108, 35, 0, [2]string{"Ghost", "Dark"}, " ", "", "", 0, 0, [4]string{"Dark Pulse", "Shadow Ball", "Bug Buzz", "Psychic"}}, 
    "Togekiss": PokemonData{"175", "Togekiss", 175, 85, 50, 95, 120, 115, 80, 0, [2]string{"Normal", "Flying"}, "", "", "", 0, 0, [4]string{"Air Slash", "Aura Sphere", "Surf", "Thunderbolt"}}, 
    "Roserade": PokemonData{"027", "Roserade", 27, 60, 70, 55, 125, 105, 90,  0,[2]string{"Grass", "Poison"}, "", "", "", 0, 0, [4]string{"Energy Ball", "Sludge Bomb", "Extrasensory", "Shadow Ball"}}, 
    "Milotic": PokemonData{"139", "Milotic", 139, 95, 60, 79, 100, 125, 81,  0,[2]string{"Water"}, "", "", "", 0, 0, [4]string{"Surf", "Ice Beam", "Dragon Pulse", "Flamethrower"}}, 
    "Garchomp": PokemonData{"111", "Garchomp", 111, 108, 130, 95, 80, 85, 102, 0, [2]string{"Dragon", "Ground"}, "", "", "", 0, 0, [4]string{"Dragon Claw", "Earthquake", "Flamethrower", "Rock Slide"}},
    "Torterra": PokemonData{"003", "Torterra", 3, 95, 109, 105, 75, 85, 56,  0,[2]string{"Grass", "Ground"}, "", "", "", 0, 0, [4]string{"Energy Ball", "Earthquake", "Crunch", "Rock Slide"}},
    "Infernape": PokemonData{"006", "Infernape", 6, 76, 104, 71, 104, 71, 108,  0,[2]string{"Fire", "Fighting"}, "", "", "", 0, 0, [4]string{"Flamethrower", "Brick Break", "Aura Sphere", "Poison Jab"}},
    "Empoleon": PokemonData{"009", "Empoleon", 9, 84, 86, 88, 111, 101, 60,  0,[2]string{"Water", "Steel"}, "", "", "", 0, 0, [4]string{"Surf", "Flash Cannon", "Aerial Ace", "Strength"}},
    "Staraptor": PokemonData{"012", "Staraptor", 12, 85, 120, 70, 50, 50, 100,  0, [2]string{"Normal", "Flying"}, "", "", "", 0, 0, [4]string{"Strength", "Aerial Ace", "Brick Break", "Leaf Blade"}},
    "Bibarel": PokemonData{"014", "Bibarel", 14, 79, 85, 60, 55, 60, 71,  0, [2]string{"Normal", "Water"}, "", "", "", 0, 0, [4]string{"Strength", "Waterfall", "Aerial Ace", "Brick Break"}},
    "Luxray": PokemonData{"019", "Luxray", 19, 80, 120, 79, 95, 79, 70,  0,[2]string{"Electric"}, "", "", "", 0, 0, [4]string{"Spark", "Crunch", "Brick Break", "Strength"}},
    "Alakazam": PokemonData{"022", "Alakazam", 22, 55, 50, 45, 135, 85, 120,  0,[2]string{"Psychic"}, "", "", "", 0, 0, [4]string{"Psychic", "Shadow Ball", "Dark Pulse", "Aura Sphere"}},
    "Gyarados": PokemonData{"024", "Gyarados", 24, 95, 125, 79, 60, 100, 81,  0,[2]string{"Water", "Flying"}, "", "", "", 0, 0, [4]string{"Waterfall", "Aerial Ace", "Earthquake", "Hyper Voice"}},
    "Steelix": PokemonData{"035", "Steelix", 35, 75, 85, 200, 55, 65, 30,  0,[2]string{"Steel", "Ground"}, "", "", "", 0, 0, [4]string{"Iron Head", "Earthquake", "Strength", "Stone Edge"}},
    "Machamp": PokemonData{"042", "Machamp", 42, 90, 130, 80, 65, 85, 55,  0,[2]string{"Fighting"}, "", "", "", 0, 0, [4]string{"Brick Break", "Rock Slide", "Strength", "Earthquake"}},
    "Gastrodon": PokemonData{"061", "Gastrodon", 61, 111, 83, 68, 92, 82, 39,  0,[2]string{"Water", "Ground"}, "", "", "", 0, 0, [4]string{"Earth Power", "Surf", "Ice Beam", "Sludge Bomb"}},
   // Add more entries as needed
}



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
} */

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

var PokemonList = make(map[string]PokemonData)

func InitData() (map[string]PokemonData, error) {
	// Read data from pokedex.json and populate PokemonList
	file, err := os.Open("./crawler/pokedex.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	pokedex := []PokemonData{}
	if err := decoder.Decode(&pokedex); err != nil {
		return nil, err
	}

	pokemonList := make(map[string]PokemonData)
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

	return pokemonList, nil
}


/* func main() {
    // Initialize data from pokedex.json
    err := InitData()
    if err != nil {
        fmt.Println("Error initializing data:", err)
        return
    }

    // Now you can call other functions that use the loaded data
    PrintAllPokemon()

    // Other code...
}
 */