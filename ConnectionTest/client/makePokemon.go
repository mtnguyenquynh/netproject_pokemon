

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