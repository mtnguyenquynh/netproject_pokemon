package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Pokemon struct {
	Index       string   `json:"index"`
	Name        string   `json:"name"`
	Exp         int      `json:"exp"`
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
	ImageURL    string   `json:"image_url"`
	Stats PokemonStats `json:"stats"`
}

type PokemonTypes struct {
	Type1 string `json:"type1"`
	Type2 string `json:"type2"`
}

type PokemonStats struct {
	HP     string `json:"HP"`
	Attack string `json:"Attack"`
	Defense string `json:"Defense"`
	Speed  string `json:"Speed"`
	SpAtk  string `json:"Sp Atk"`
	SpDef  string `json:"Sp Def"`
}

type Player struct {
	Name     string
	Pokemons []Pokemon
	Active   int
}

type Battle struct {
	Player1 Player
	Player2 Player
}

func loadPokemons(fileName string) ([]Pokemon, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var pokemons []Pokemon
	err = json.Unmarshal(file, &pokemons)
	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

func (p *Pokemon) GetStat(stat string) int {
	var value string
	switch stat {
	case "HP":
		value = p.Stats.HP
	case "Attack":
		value = p.Stats.Attack
	case "Defense":
		value = p.Stats.Defense
	case "Speed":
		value = p.Stats.Speed
	case "Sp Atk":
		value = p.Stats.SpAtk
	case "Sp Def":
		value = p.Stats.SpDef
	}
	intValue, _ := strconv.Atoi(value)
	return intValue
}

func (p *Pokemon) SetStat(stat string, value int) {
	switch stat {
	case "HP":
		p.Stats.HP = strconv.Itoa(value)
	}
}

func (b *Battle) Fight(attacker, defender *Pokemon) {
	rand.Seed(time.Now().UnixNano())
	attackType := rand.Intn(2) // 0 for normal, 1 for special

	var damage int
	if attackType == 0 {
		damage = (attacker.GetStat("Attack") - defender.GetStat("Defense")) * 1
		fmt.Printf("%s uses a normal attack on %s for %d damage!\n", attacker.Name, defender.Name, damage)
	} else {
		damage = (attacker.GetStat("Sp Atk") - defender.GetStat("Sp Def")) * 1
		fmt.Printf("%s uses a special attack on %s for %d damage!\n", attacker.Name, defender.Name, damage)
	}

	if damage < 0 {
		damage = 5
	}

	newHP := defender.GetStat("HP") - damage
	defender.SetStat("HP", newHP)
	fmt.Printf("%s's remaining HP: %d\n", defender.Name, newHP)
}

func (b *Battle) BattleTurn() {
	p1 := &b.Player1.Pokemons[b.Player1.Active]
	p2 := &b.Player2.Pokemons[b.Player2.Active]

	if p1.GetStat("Speed") > p2.GetStat("Speed") {
		b.Fight(p1, p2)
		if p2.GetStat("HP") <= 0 {
			fmt.Printf("%s is knocked out!\n", p2.Name)
			b.Player2.Pokemons = append(b.Player2.Pokemons[:b.Player2.Active], b.Player2.Pokemons[b.Player2.Active+1:]...)
			if len(b.Player2.Pokemons) > 0 {
				b.Player2.Active = 0 // switch to next Pokémon
			}
			return
		}
		b.Fight(p2, p1)
		if p1.GetStat("HP") <= 0 {
			fmt.Printf("%s is knocked out!\n", p1.Name)
			b.Player1.Pokemons = append(b.Player1.Pokemons[:b.Player1.Active], b.Player1.Pokemons[b.Player1.Active+1:]...)
			if len(b.Player1.Pokemons) > 0 {
				b.Player1.Active = 0 // switch to next Pokémon
			}
		}
	} else {
		b.Fight(p2, p1)
		if p1.GetStat("HP") <= 0 {
			fmt.Printf("%s is knocked out!\n", p1.Name)
			b.Player1.Pokemons = append(b.Player1.Pokemons[:b.Player1.Active], b.Player1.Pokemons[b.Player1.Active+1:]...)
			if len(b.Player1.Pokemons) > 0 {
				b.Player1.Active = 0 // switch to next Pokémon
			}
			return
		}
		b.Fight(p1, p2)
		if p2.GetStat("HP") <= 0 {
			fmt.Printf("%s is knocked out!\n", p2.Name)
			b.Player2.Pokemons = append(b.Player2.Pokemons[:b.Player2.Active], b.Player2.Pokemons[b.Player2.Active+1:]...)
			if len(b.Player2.Pokemons) > 0 {
				b.Player2.Active = 0 // switch to next Pokémon
			}
		}
	}
}

func main() {
	pokemons, err := loadPokemons("pokedex.json")
	if err != nil {
		log.Fatalf("Failed to load pokemons: %v", err)
	}

	if len(pokemons) < 6 {
		log.Fatalf("Not enough pokemons to form two teams")
	}

	// Shuffle the list of pokemons to randomize the selection
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pokemons), func(i, j int) {
		pokemons[i], pokemons[j] = pokemons[j], pokemons[i]
	})

	// Assigning three random pokemons to each player
	p1Pokemons := pokemons[:3]
	p2Pokemons := pokemons[3:6]

	p1 := Player{
		Name:     "Player 1",
		Pokemons: p1Pokemons,
		Active:   0,
	}

	p2 := Player{
		Name:     "Player 2",
		Pokemons: p2Pokemons,
		Active:   0,
	}

	printChosenPokemons(p1)
	printChosenPokemons(p2)

	battle := Battle{
		Player1: p1,
		Player2: p2,
	}

	for {
		battle.BattleTurn()

		// Check if any player has won
		if len(battle.Player1.Pokemons) == 0 {
			fmt.Printf("%s wins the battle!\n", battle.Player2.Name)
			break
		} else if len(battle.Player2.Pokemons) == 0 {
			fmt.Printf("%s wins the battle!\n", battle.Player1.Name)
			break
		}
	}

	fmt.Println("The battle has ended.")
}

func printChosenPokemons(player Player) {
	fmt.Printf("%s chose the following Pokémon:\n", player.Name)
	for _, pokemon := range player.Pokemons {
		fmt.Printf("- %s\n", pokemon.Name)
	}
}