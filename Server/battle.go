// battle.go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

func battle(participant1, participant2 *Participant) (*Participant, *Participant) {
	var surrendered bool
	for {

		// Start the battle
		winner, loser := battleRound(participant1, participant2)
		msg := fmt.Sprintf("\n%s wins the round - %s lost", winner.player.Name, loser.player.Name)

		msgCh <- msg + "#"
		if loser.turn == 0 {
			msg := fmt.Sprintf("\nBATTLE END!!! \n%s has no turns left. %s wins!", loser.player.Name, winner.player.Name)

			msgCh <- msg + "#"
			totalExp := 0
			for _, pokemon := range loser.player.PokemonList {
				totalExp += pokemon.Exp
			}

			// Distribute the total experience to the winning team
			expPerPokemon := totalExp / 3
			for i := range winner.player.PokemonList {
				winner.player.PokemonList[i].Exp += expPerPokemon
			}
			return winner, loser

		}
		// Ask the losing participant to choose another Pokemon
		pokeList := listPokemon(loser.player.PokemonList)

		loser.curPokemon, surrendered = readPokemonFromClient(loser.conn, "\nYou lost the round, Let's choose another Pokemon\n"+pokeList+"Your choice: ", loser.player.PokemonList)
		if surrendered {
			totalExp := 0
			for _, pokemon := range loser.player.PokemonList {
				totalExp += pokemon.Exp
			}

			// Distribute the total experience to the winning team
			expPerPokemon := totalExp / 3
			for i := range winner.player.PokemonList {
				winner.player.PokemonList[i].Exp += expPerPokemon
			}
			return winner, loser
		}
	}
}

func battleRound(participant1, participant2 *Participant) (*Participant, *Participant) {
	var messages []string
	// Announce the current Pokemon
	msg := fmt.Sprintf("---%s chose %s\n%s chose %s\n", participant1.player.Name, participant1.curPokemon.Name, participant2.player.Name, participant2.curPokemon.Name)

	messages = append(messages, msg)
	messages = append(messages, "------------BATTLE REPORT------------\n")

	var winner, loser *Participant
	var attacker, defender *Participant
	if participant1.curPokemon.Speed > participant2.curPokemon.Speed {
		attacker = participant1
		defender = participant2
	} else if participant1.curPokemon.Speed < participant2.curPokemon.Speed {
		attacker = participant2
		defender = participant1
	} else {
		// If the speeds are equal, randomly choose the attacker
		if rand.Intn(2) == 0 {
			attacker = participant1
			defender = participant2
		} else {
			attacker = participant2
			defender = participant1
		}
	}
	for participant1.curPokemon.HP > 0 && participant2.curPokemon.HP > 0 {

		// Player's turn to choose attack type
		attackTypes := []string{"normal", "special"}
		attackType := attackTypes[rand.Intn(len(attackTypes))]

		// Calculate and apply damage
		dmg := calculateDamage(&attacker.curPokemon, &defender.curPokemon, attackType)
		defender.curPokemon.HP -= dmg
		msg := fmt.Sprintf("%s attacked %s with %s attack dealing %d damage.\n", attacker.curPokemon.Name, defender.curPokemon.Name, attackType, dmg)
		messages = append(messages, msg)
		if defender.curPokemon.HP <= 0 {
			msg := fmt.Sprintf("%s fainted.\n", defender.curPokemon.Name)
			defender.turn--
			for i := range defender.player.PokemonList {
				if defender.player.PokemonList[i].Name == defender.curPokemon.Name {
					// update hp to 0
					defender.player.PokemonList[i].HP = 0
					// update deployable
					defender.player.PokemonList[i].Deployable = false
				} else {
					// update hp winner

					attacker.player.PokemonList[i].HP = attacker.curPokemon.HP
				}
			}

			loser = defender
			winner = attacker
			messages = append(messages, msg)
			break
		}

		// Swap attacker and defender for the next round
		attacker, defender = defender, attacker
	}
	msg = fmt.Sprintf("%s has %d HP left.\n", attacker.curPokemon.Name, attacker.curPokemon.HP)
	messages = append(messages, msg)
	// announce the turns
	msg = fmt.Sprintf("%s has %d turns left.\n", participant1.player.Name, participant1.turn)
	messages = append(messages, msg)
	msg = fmt.Sprintf("%s has %d turns left.\n", participant2.player.Name, participant2.turn)
	messages = append(messages, msg)
	messages = append(messages, "------END BATTLE REPORT-----")

	msgCh <- strings.Join(messages, "") + "#"

	return winner, loser
}

func readPokemonFromClient(conn net.Conn, msg string, pokemonList []Pokemon) (Pokemon, bool) {
	// conn.Write([]byte(msg))
	var chosenPokemon Pokemon
	msgChOne <- Message{msg: msg, conn: conn}

	for {
		reader := bufio.NewReader(conn)
		pokemonIndex, err := reader.ReadString('\n')

		if err != nil {
			return Pokemon{}, true
		}
		pokemonIndex = strings.TrimSpace(pokemonIndex)
		index, _ := strconv.Atoi(pokemonIndex)

		// Check if the player wants to surrender
		if index == -1 {
			return Pokemon{}, true
		} else {
			chosenPokemon = pokemonList[index-1]
		}
		// Check if the chosen Pokemon is deployable
		if chosenPokemon.Deployable {
			return chosenPokemon, false
		} else if !chosenPokemon.Deployable {

			// If the Pokemon is not deployable, send a message to the client and ask for another Pokemon
			msg := "This Pokemon lost the ability to fight. Please choose another one.\n#"
			conn.Write([]byte(msg))
		}
	}

}

func listPokemon(pokemonList []Pokemon) string {
	var listOfPokemon []string
	for i, p := range pokemonList {
		listOfPokemon = append(listOfPokemon, fmt.Sprintf("%d. %s\n", i+1, p.Name))
	}
	return strings.Join(listOfPokemon, "") + "#"
}
