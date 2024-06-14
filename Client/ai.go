package main

import (
	"fmt"
)

// Handles AI choosing an action (only attack for now)
func ChooseActionAI(aiInput *UserInput, userInput *UserInput) *UserInput {
	aiInput.action = "attack"
	aiInput = ChooseMoveAI(aiInput, userInput)
	return aiInput
}

// Handles AI choosing a move
func ChooseMoveAI(aiInput *UserInput, userInput *UserInput) *UserInput {
    // AI chooses the move with the highest damage
    damage := 0
    var chosenMove string
    for _, move := range aiInput.activePokemon.Moves {
        dmg, _ := DamageCalc(aiInput.activePokemon, userInput.activePokemon, &move)
        if dmg >= damage {
            damage = dmg
            chosenMove = move.MoveName
        }
    }
    aiInput.move = chosenMove
    return aiInput
}


// Handles AI sending out a new pokemon when needed
func ReplaceFaintedPokemonAI(aiInput *UserInput, userInput *UserInput) *UserInput {
    // AI sends out pokemon with strongest move
    damage := 0
    var chosenPokemon *Pokemon
    for _, pokemon := range aiInput.team {
        // Don't consider fainted pokemon
        if pokemon.fainted {
            continue
        }
        // AI chooses the move with the highest damage
        for _, move := range pokemon.Moves {
            dmg, _ := DamageCalc(aiInput.activePokemon, userInput.activePokemon, &move)
            if dmg > damage {
                damage = dmg
                chosenPokemon = pokemon
            }
        }
    }
    
    // Check if a pokemon to send out was found
    if chosenPokemon == nil {
        aiInput.gameOver = true
    } else {
        fmt.Println("[[", aiInput.username, "]] Go", chosenPokemon.Name, "\n")
        aiInput.activePokemon = chosenPokemon
    }
    return aiInput
}
