package main

import (
	"fmt"
)

func main() {
    // Initialize data from pokedex.json
    pokemonList, err := InitData()
    if err != nil {
        fmt.Println("Error initializing data:", err)
        return
    }

    spiritomb := NewPokemon("Spiritomb", true, pokemonList)
    lucario := NewPokemon("Lucario", true, pokemonList)
    togekiss := NewPokemon("Togekiss", true, pokemonList)
    roserade := NewPokemon("Roserade", true, pokemonList)
    milotic := NewPokemon("Milotic", true, pokemonList)
    garchomp := NewPokemon("Garchomp", true, pokemonList)

    cynthiasTeam := []*Pokemon{spiritomb, lucario, togekiss, roserade, milotic, garchomp}

    myInput := &UserInput{"Ash", "", nil, nil, "", false, false}
    cynthiasInput := &UserInput{"Cynthia", "", spiritomb, cynthiasTeam, "", true, false}

    fmt.Println()
    ChooseName(myInput)
    myInput = ChooseTeam(myInput, pokemonList)
    myInput.activePokemon = myInput.team[0]

    Battle(myInput, cynthiasInput)
}
 
   