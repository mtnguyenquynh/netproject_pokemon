package main

import (
	"fmt"
)

func main() {
	// Initialize data from pokedex.json
    //server
	InitData()


	venusaur := NewPokemon("Venusaur", true)
	charmeleon := NewPokemon("Charmeleon", true)
	wartortle := NewPokemon("Wartortle", true)
	blastoise := NewPokemon("Blastoise", true)
	caterpie := NewPokemon("Caterpie", true)
	bulbasaur := NewPokemon("Bulbasaur", true)

	cynthiasTeam := []*Pokemon{venusaur, charmeleon, wartortle, blastoise, caterpie, bulbasaur}

    // client 
    // get name from the client 
	myInput := &UserInput{"Ash", "", nil, nil, "", false, false}

    //server
	cynthiasInput := &UserInput{"Cynthia", "", venusaur, cynthiasTeam, "", true, false}

	fmt.Println()
	ChooseName(myInput)
	ChooseTeam(myInput)
	myInput.activePokemon = myInput.team[0]

    //fmt.Println("error")

	Battle(myInput, cynthiasInput)
}
