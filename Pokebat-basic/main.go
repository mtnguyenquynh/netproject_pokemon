package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Pokemon struct to hold all necessary attributes
type Pokemon struct {
	Index       string   `json:"index"`
	Name        string   `json:"name"`
	Level       int      `json:"level"`
	Experience  int      `json:"experience"`
	HP          int      `json:"hp"`
	Attack      int      `json:"attack"`
	Defense     int      `json:"defense"`
	SpAttack    int      `json:"sp_attack"`
	SpDefense   int      `json:"sp_defense"`
	Speed       int      `json:"speed"`
	TotalEVs    int      `json:"total_evs"`
	Type        []string `json:"type"`
	Description string   `json:"description"`
	Height      string   `json:"height"`
	Weight      string   `json:"weight"`
}

// Calculate the required experience for the next level
func calculateNextLevelExp(level int) int {
	return int(math.Pow(2, float64(level-1)) * 100)
}

// Level up a Pokemon if it has enough experience
func levelUp(p *Pokemon) {
	for p.Experience >= calculateNextLevelExp(p.Level) {
		p.Level++

		// Recalculate attributes
		p.HP = int(float64(p.HP) * 1.1)
		p.Attack = int(float64(p.Attack) * 1.1)
		p.Defense = int(float64(p.Defense) * 1.1)
		p.SpAttack = int(float64(p.SpAttack) * 1.1)
		p.SpDefense = int(float64(p.SpDefense) * 1.1)
		p.TotalEVs = p.HP + p.Attack + p.Defense + p.SpAttack + p.SpDefense + p.Speed

		fmt.Printf("%s leveled up to %d!\n", p.Name, p.Level)
	}
}

// Transfer experience from one Pokemon to another of the same type
func transferExp(source, target *Pokemon) {
	if hasSameType(source.Type, target.Type) {
		target.Experience += source.Experience
		fmt.Printf("%s transferred %d EXP to %s\n", source.Name, source.Experience, target.Name)
	} else {
		fmt.Printf("Cannot transfer EXP: %s and %s are not of the same type\n", source.Name, target.Name)
	}
}

// Check if two Pokemon have at least one same type
func hasSameType(types1, types2 []string) bool {
	typeSet := make(map[string]bool)
	for _, t := range types1 {
		typeSet[t] = true
	}
	for _, t := range types2 {
		if typeSet[t] {
			return true
		}
	}
	return false
}

// Calculate damage based on attack type
func calculateDamage(attacker, defender *Pokemon, attackType string) int {
	if attackType == "normal" {
		return max(attacker.Attack-defender.Defense, 0)
	}
	elementalMultiplier := 1.2 // Example value, should be dynamic
	return max(int(float64(attacker.SpAttack)*elementalMultiplier)-defender.SpDefense, 0)
}

// Helper function to get the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Battle function
func battle(pokemon1, pokemon2 *Pokemon) (*Pokemon, *Pokemon) {
	var first, second *Pokemon
	if pokemon1.Speed > pokemon2.Speed {
		first, second = pokemon1, pokemon2
	} else {
		first, second = pokemon2, pokemon1
	}

	for first.HP > 0 && second.HP > 0 {
		// Player's turn to choose attack type
		var playerAttackType string
		fmt.Println("Choose your attack type: normal or special")
		fmt.Scan(&playerAttackType)

		// Calculate and apply damage
		dmg := calculateDamage(first, second, playerAttackType)
		second.HP -= dmg
		fmt.Printf("%s attacked %s with %s attack dealing %d damage.\n", first.Name, second.Name, playerAttackType, dmg)

		if second.HP <= 0 {
			fmt.Printf("%s fainted.\n", second.Name)
			break
		}

		// Machine's turn to choose attack type randomly
		machineAttackTypes := []string{"normal", "special"}
		machineAttackType := machineAttackTypes[rand.Intn(len(machineAttackTypes))]

		// Calculate and apply damage
		dmg = calculateDamage(second, first, machineAttackType)
		first.HP -= dmg
		fmt.Printf("%s attacked %s with %s attack dealing %d damage.\n", second.Name, first.Name, machineAttackType, dmg)

		if first.HP <= 0 {
			fmt.Printf("%s fainted.\n", first.Name)
			break
		}
	}

	var winner, loser *Pokemon
	if second.HP <= 0 {
		winner, loser = first, second
	} else {
		winner, loser = second, first
	}
	fmt.Printf("%s wins!\n", winner.Name)
	return winner, loser
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Example Pokemon
	bulbasaur := &Pokemon{
		Index:       "001",
		Name:        "Bulbasaur",
		Level:       5,
		Experience:  2000,
		HP:          45,
		Attack:      49,
		Defense:     49,
		SpAttack:    65,
		SpDefense:   65,
		Speed:       45,
		TotalEVs:    318,
		Type:        []string{"grass", "poison"},
		Description: "For some time after its birth, it grows by gaining nourishment from the seed on its back.",
		Height:      "0.7 m",
		Weight:      "6.9 kg",
	}

	charmander := &Pokemon{
		Index:       "004",
		Name:        "Charmander",
		Level:       5,
		Experience:  500,
		HP:          39,
		Attack:      52,
		Defense:     43,
		SpAttack:    60,
		SpDefense:   50,
		Speed:       65,
		TotalEVs:    309,
		Type:        []string{"fire"},
		Description: "The fire on the tip of its tail is a measure of its life. If healthy, its tail burns intensely.",
		Height:      "0.6 m",
		Weight:      "8.5 kg",
	}

	squirtle := &Pokemon{
		Index:       "007",
		Name:        "Squirtle",
		Level:       5,
		Experience:  500,
		HP:          44,
		Attack:      48,
		Defense:     65,
		SpAttack:    50,
		SpDefense:   64,
		Speed:       43,
		TotalEVs:    314,
		Type:        []string{"water"},
		Description: "It shelters itself in its shell then strikes back with spouts of water at every opportunity.",
		Height:      "0.5 m",
		Weight:      "9 kg",
	}

	// Get player choices
	var playerPokemon *Pokemon
	fmt.Println("Choose your Pokemon: 1 for Bulbasaur, 2 for Charmander, 3 for Squirtle")
	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1:
		playerPokemon = bulbasaur
	case 2:
		playerPokemon = charmander
	case 3:
		playerPokemon = squirtle
	default:
		fmt.Println("Invalid choice")
		return
	}

	// Machine (random) choice for Pokemon
	machineChoices := []*Pokemon{bulbasaur, charmander, squirtle}
	machinePokemon := machineChoices[rand.Intn(len(machineChoices))]

	fmt.Printf("You chose %s.\n", playerPokemon.Name)
	fmt.Printf("Machine chose %s.\n", machinePokemon.Name)

	// Battle between player and machine
	winner, loser := battle(playerPokemon, machinePokemon)
	fmt.Printf("Winner: %+v\nLoser: %+v\n", winner, loser)
}
