package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	
)

// gets user stdin to decide whether to attack or switch
func ChooseAction(input *UserInput) *UserInput {
	fmt.Println("Type a number and hit ENTR to choose an action:")
	fmt.Println("(1) Fight")
	fmt.Println("(2) Pokemon")
	fmt.Println("(3) Bag")
	fmt.Println("(4) Run")
	fmt.Println()
	var mv string
	canContinue := true
	for canContinue {
		fmt.Scanln(&mv)
		fmt.Println()
		switch mv {
		case "1":
			canContinue = false
		case "2":
			canContinue = false
		case "3":
			fmt.Println("Your bag is empty, try again\n")
		case "4":
			fmt.Println("Only cowards run away, try again\n")
		default:
			fmt.Println("[[ INVALID INPUT ]] Try again\n")
		}
	}
	action, _ := strconv.Atoi(mv)
	res := input
	switch action {
	case 1:
		input.action = "attack"
		res = ChooseMove(input)
	case 2:
		input.action = "switch"
		res = SwitchPokemon(input)
	default:
		fmt.Println("[[ SOMETHING UNEXPECTED OCCURRED (controls.go:line 43) ]]")
	}
	return res
}

func ChooseMove(input *UserInput) *UserInput {
	pokemon := input.activePokemon
	fmt.Println("Type a number and hit ENTER to choose a move:")
	
	// Check if the pokemon has any moves
	if len(pokemon.Moves) == 0 {
		fmt.Println("[[ NO MOVES AVAILABLE ]]")
		return input
	}
	
	for i := 0; i < len(pokemon.Moves); i++ {
		fmt.Printf("(%v) %v\n", i+1, pokemon.Moves[i].MoveName)
	}
	fmt.Println()

	var mv string
	for {
		fmt.Scanln(&mv)
		fmt.Println()
		moveIdx, err := strconv.Atoi(mv)
		if err == nil && moveIdx > 0 && moveIdx <= len(pokemon.Moves) {
			input.move = pokemon.Moves[moveIdx-1].MoveName
			break
		}
		fmt.Println("[[ INVALID INPUT ]] Try again\n")
	}
	
	return input
}


// gets user stdin to decide which Pokemon to send out
// (after previous one has fainted)
func ReplaceFaintedPokemon(input *UserInput) *UserInput {
	team := input.team
	var nonFaintedPokemon []*Pokemon
	for i := 0; i < len(team); i++ {
		if !team[i].fainted {
			nonFaintedPokemon = append(nonFaintedPokemon, team[i])
		}
	}

	// Are there more pokemon to fight?
	if len(nonFaintedPokemon) == 0 {
		input.gameOver = true
		return input
	}

	fmt.Println("Type a number and hit ENTR to choose a pokemon:")
	for i := 0; i < len(nonFaintedPokemon); i++ {
		fmt.Printf("(%v) %v\n", i+1, nonFaintedPokemon[i].Name)
	}
	fmt.Println()
	var mv string
	var res int
	for {
		fmt.Scanln(&mv)
		fmt.Println()
		pokemonIdx, err := strconv.Atoi(mv)
		if err == nil && (0 < pokemonIdx) && (pokemonIdx < len(nonFaintedPokemon)+1) {
			res = pokemonIdx - 1
			break
		}
		fmt.Println("[[ INVALID INPUT ]] Try again")
	}
	fmt.Println("[[", input.username, "]]", "Go", nonFaintedPokemon[res].Name, "\n")
	input.activePokemon = nonFaintedPokemon[res]
	return input
}

// gets user stdin to decide which Pokemon to send out
// (switching during turn)
func SwitchPokemon(input *UserInput) *UserInput {
	activePokemon := input.activePokemon
	team := input.team
	fmt.Println("Type a number and hit ENTR to choose a pokemon:")
	var nonFaintedPokemon []*Pokemon
	fmt.Println("(0) [[ GO BACK ]]")
	for i := 0; i < len(team); i++ {
		if !team[i].fainted && team[i] != activePokemon {
			nonFaintedPokemon = append(nonFaintedPokemon, team[i])
			fmt.Printf("(%v) %v\n", len(nonFaintedPokemon), team[i].Name)
		}
	}
	fmt.Println()
	var mv string
	var res int
	for {
		fmt.Scanln(&mv)
		fmt.Println()
		pokemonIdx, err := strconv.Atoi(mv)
		if err == nil && (0 <= pokemonIdx) && (pokemonIdx < len(nonFaintedPokemon)+1) {
			res = pokemonIdx - 1
			if res == -1 {
				input = ChooseAction(input)
				return input
			}
			break
		}
		fmt.Println("[[ INVALID INPUT ]] Try again")
	}
	fmt.Println("[[", input.username, "]]", "Come back", activePokemon.Name)
	fmt.Println("[[", input.username, "]]", "Go", nonFaintedPokemon[res].Name, "\n")
	input.activePokemon = nonFaintedPokemon[res]
	return input
}

// func ChooseName(input *UserInput) *UserInput {
// 	var mv string
// 	fmt.Println("What is your name?\n")
// 	fmt.Scanln(&mv)
// 	fmt.Println()
// 	input.username = mv
// 	return input
// }

// helper function to pad string for formatting
func padString(str string, n int) string {
	if len(str) >= n {
		return str
	}
	return str + strings.Repeat(" ", n-len(str))
}

// helper function to print all available Pokémon names
func PrintAllPokemon() {
	var s string
	var i int
	for _, data := range pokemonList {
		name := data.Name
		if name == "MissingNo" {
			continue
		} else if i < 5 {
			s += padString(name, 12)
			i++
		} else {
			fmt.Println(s)
			s = padString(name, 12)
			i = 1
		}
	}
	if len(s) > 0 {
		fmt.Println(s)
	}
	fmt.Println()
	os.Stdout.Sync()
}


// choice to either build a team or use a random one
func ChooseTeam(input *UserInput) *UserInput {
	fmt.Println("Type a number and press ENTER to choose an option")
	fmt.Println("(1) Build a team")
	fmt.Println("(2) Use a random team\n")
	var mv string
	for {
		fmt.Scanln(&mv)
		fmt.Println()
		if mv == "1" {
			input = BuildTeam(input, false)
			break
		} else if mv == "2" {	
			input = RandomTeam(input, false)
			break
		} else {
			fmt.Println("[[ INVALID INPUT ]] Try again")
			continue
		}
	}
	return input
}

// menu to choose your team
func BuildTeam(input *UserInput, makeStrong bool) *UserInput {
	fmt.Println("[[ CHOOSE YOUR TEAM ]]\n")
	var mv string
	var team []*Pokemon
	var printTeam []string
	for i := 0; i < 6; i++ {
		for {
			PrintAllPokemon()
			fmt.Println("Your team so far:", printTeam)
			fmt.Println("Choose your (", i+1, ") pokemon\n")
			fmt.Scanln(&mv)
			fmt.Println()
			_, isValid := pokemonList[mv]
			if !isValid {
				fmt.Println("[[ INVALID INPUT ]] Try again")
				continue
			} else {
				team = append(team, NewPokemon(mv, makeStrong))
				printTeam = append(printTeam, mv)
				break
			}
		}
	}
	fmt.Println("Your team is:", printTeam, "\n")
	input.team = team
	return input
}

func RandomTeam(input *UserInput, makeStrong bool) *UserInput {
	// fmt.Println(pokemonList)
	// rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	var team []*Pokemon
	var printTeam []string

	// Put all Pokémon names in a list for random choosing
	var allPokemon []string
	for p := range pokemonList {
		if p == "MissingNo" { // Skip "MissingNo" if present
			continue
		}
		allPokemon = append(allPokemon, p)
	}

	// Randomly choose 6 Pokémon
	for i := 0; i < 6; i++ {
		
		// fmt.Println(len(allPokemon))
	
		pokemon := allPokemon[rand.Intn(len(allPokemon))]



		team = append(team, NewPokemon(pokemon, makeStrong))
		printTeam = append(printTeam, pokemon)
	}

	

	fmt.Println("Your team is:", printTeam, "\n")
	input.team = team
	return input
}


