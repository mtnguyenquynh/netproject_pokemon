package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Pokemon struct to hold all necessary attributes
type Pokemon struct {
	Index      int
	Name       string
	Level      int
	Experience int
	HP         int
	Attack     int
	Defense    int
	SpAttack   int
	SpDefense  int
	Speed      int
	TotalEVs   int
	Type       []string
	Description string
	Height     string
	Weight     string
	EV         float64
}

// Calculate the required experience for the next level
func calculateNextLevelExp(level int) int {
	return int(math.Pow(2, float64(level-1)) * 100)
}

// Level up a Pokemon if it has enough experience
func levelUp(p *Pokemon) {
	for p.Experience >= calculateNextLevelExp(p.Level) {
		p.Level++
		ev := p.EV

		// Recalculate attributes
		p.HP = int(float64(p.HP) * (1 + ev))
		p.Attack = int(float64(p.Attack) * (1 + ev))
		p.Defense = int(float64(p.Defense) * (1 + ev))
		p.SpAttack = int(float64(p.SpAttack) * (1 + ev))
		p.SpDefense = int(float64(p.SpDefense) * (1 + ev))
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
		attackType := randomChoice([]string{"normal", "special"})
		dmg := calculateDamage(first, second, attackType)
		second.HP -= dmg
		fmt.Printf("%s attacked %s with %s attack dealing %d damage.\n", first.Name, second.Name, attackType, dmg)

		if second.HP <= 0 {
			fmt.Printf("%s fainted.\n", second.Name)
			break
		}

		attackType = randomChoice([]string{"normal", "special"})
		dmg = calculateDamage(second, first, attackType)
		first.HP -= dmg
		fmt.Printf("%s attacked %s with %s attack dealing %d damage.\n", second.Name, first.Name, attackType, dmg)

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

// Helper function to choose a random element from a slice of strings
func randomChoice(choices []string) string {
	rand.Seed(time.Now().UnixNano())
	return choices[rand.Intn(len(choices))]
}

func main() {
	// Example Pokemon
	bulbasaur := &Pokemon{
		Index:       1,
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
		EV:          0.5,
	}

	charmander := &Pokemon{
		Index:       4,
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
		EV:          0.5,
	}

	squirtle := &Pokemon{
		Index:       7,
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
		EV:          0.5,
	}

	// Level up Bulbasaur
	levelUp(bulbasaur)
	fmt.Printf("%+v\n", bulbasaur)

	// Battle between Charmander and Squirtle
	winner, loser := battle(charmander, squirtle)
	fmt.Printf("Winner: %+v\nLoser: %+v\n", winner, loser)
}
