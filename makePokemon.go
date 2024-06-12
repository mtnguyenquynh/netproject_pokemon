/* package main

import "math/rand"

// does math calculation for a pokemon's real HP stat
func calculateHp(pokemon *Pokemon, iv int, ev int) (int) {
    return (((2 * pokemon.baseHp) + iv + (ev / 4)) * pokemon.level) / 100 + pokemon.level + 10
}


// does math calculation for a pokemon's real HP stat
func calculateOtherStat(pokemon *Pokemon, iv int, ev int, baseStat int) (int) {
    return (((2 * baseStat) + iv + (ev / 4)) * pokemon.level) / 100 + 5
}


// stat calculator
// takes a pokemon and a bool indicating whether to make it strong or not
// and fills in that pokemon's actual stats (based on base stats, IVs, etc) 
func initializeStats(pokemon *Pokemon, makeStrong bool) {
    if makeStrong {   // cynthia AI will get stronger pokemon
        pokemon.level = 60
        IV := 31
        EV := 252
        pokemon.hp = calculateHp(pokemon, IV, EV)
        pokemon.atk = calculateOtherStat(pokemon, IV, EV, pokemon.baseAtk)
        pokemon.def = calculateOtherStat(pokemon, IV, EV, pokemon.baseDef)
        pokemon.spatk = calculateOtherStat(pokemon, IV, EV, pokemon.baseSpatk)
        pokemon.spdef = calculateOtherStat(pokemon, IV, EV, pokemon.baseSpdef)
        pokemon.speed = calculateOtherStat(pokemon, IV, EV, pokemon.baseSpeed)
    } else {   // player gets slightly weaker pokemon
        pokemon.level = rand.Intn(10)+50
        pokemon.hp = calculateHp(pokemon, rand.Intn(20)+10, rand.Intn(100)+120)
        pokemon.atk = calculateOtherStat(pokemon, rand.Intn(20)+10, rand.Intn(100)+120, pokemon.baseAtk)
        pokemon.def = calculateOtherStat(pokemon, rand.Intn(20)+10, rand.Intn(100)+120, pokemon.baseDef)
        pokemon.spatk = calculateOtherStat(pokemon, rand.Intn(20)+10, rand.Intn(100)+120, pokemon.baseSpatk)
        pokemon.spdef = calculateOtherStat(pokemon, rand.Intn(20)+10, rand.Intn(100)+120, pokemon.baseSpdef)
        pokemon.speed = calculateOtherStat(pokemon, rand.Intn(20)+10, rand.Intn(100)+120, pokemon.baseSpeed)
    }
}


// new pokemon constructor
// didn't want to use packages so it is very raw
func NewPokemon(name string, makeStrong bool) *Pokemon {
	template := PokemonList[name]
	pokemon := Pokemon{}

	// shared among individuals
    pokemon.name = template.name
    pokemon.pokedexNumber = template.pokedexNumber
    pokemon.Type = template.Type
    pokemon.baseHp = template.baseHp
    pokemon.baseAtk = template.baseAtk
    pokemon.baseDef = template.baseDef
    pokemon.baseSpatk = template.baseSpatk
    pokemon.baseSpdef = template.baseSpdef
    pokemon.baseSpeed = template.baseSpeed
    pokemon.moves = template.moves

    pokemon.Level = template.Level
    pokemon.Description = template.Description
    pokemon.Height = template.Height
    pokemon.Weight = template.Weight
    pokemon.AccumExp = template.AccumExp
    pokemon.Exp = template.Exp


    // level and all stats handled by helper
    initializeStats(&pokemon, makeStrong)

	// remaining battle specific fields
    pokemon.nonVolatileStatus = ""
    pokemon.volatileStatus = ""
    pokemon.fainted = false

	return &pokemon
} */

package main

import (
    "math/rand"
)

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
func NewPokemon(name string, makeStrong bool, pokemonList map[string]PokemonData) *Pokemon {
    template, exists := pokemonList[name]
    if !exists {
        return nil // Return nil if the Pokemon does not exist in the Pokedex
    }

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