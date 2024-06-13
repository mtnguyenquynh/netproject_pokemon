package main

import (
	"fmt"
	"math/rand"
	"time"
	"strconv"
)


var WinMessages = []string{"That was excellent. Truly, an outstanding battle. You gave the support your Pokémon needed to maximize their power. And you guided them with certainty to secure victory. You have both passion and calculating coolness. Together, you and your Pokémon can overcome any challenge that may come your way. Those are the impressions I got from our battle. I'm glad I got to take part in the crowning of Sinnoh's new Champion! Come with me. We'll take the lift."}   // win messages
var LoseMessages = []string{"Smell ya later!", "Better luck next time", "Keep training", "Time to soft-reset", "You whited out...", "Come back when you're stronger", "Do or do not...there is no try", "Looks like you're blasting off again"}   // lose messages
// List of non-volatile status conditions
var StatusList = map[string]bool{
    "PSN": true, // Poison
    "FRZ": true, // Freeze
    "BRN": true, // Burn
    "PRZ": true, // Paralysis
}

// List of Pokemon types
var Types = map[string]int{
    "Normal": 0,
    "Fire": 1,
    "Water": 2,
    "Electric": 3,
    "Grass": 4,
    "Ice": 5,
    "Fighting": 6,
    "Poison": 7,
    "Ground": 8,
    "Flying": 9,
    "Psychic": 10,
    "Bug": 11,
    "Rock": 12,
    "Ghost": 13,
    "Dragon": 14,
    "Dark": 15,
    "Steel": 16,
}

// Type multipliers based on type
// KEY is attacking type
// VALUES are multiplier based on defending type
var Matchup = map[string][17]float64{
    "Normal":   {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0.5, 0, 1, 1, 0.5},
    "Fire":     {1, 0.5, 0.5, 1, 2, 2, 1, 1, 1, 1, 1, 2, 0.5, 1, 0.5, 1, 1},
    "Water":    {1, 2, 0.5, 1, 0.5, 1, 1, 1, 2, 1, 1, 1, 2, 1, 0.5, 1, 1},
    "Electric": {1, 1, 2, 0.5, 0.5, 1, 1, 1, 0, 2, 1, 1, 1, 1, 0.5, 1, 1},
    "Grass":    {1, 0.5, 2, 1, 0.5, 1, 1, 0.5, 2, 0.5, 1, 0.5, 2, 1, 0.5, 1, 0.5},
    "Ice":      {1, 0.5, 0.5, 1, 2, 0.5, 1, 1, 2, 2, 1, 1, 1, 1, 2, 1, 0.5},
    "Fighting": {2, 1, 1, 1, 1, 2, 1, 0.5, 1, 0.5, 0.5, 0.5, 2, 0, 1, 2, 2},
    "Poison":   {1, 1, 1, 1, 2, 1, 1, 0.5, 0.5, 1, 1, 1, 0.5, 0.5, 1, 1, 0},
    "Ground":   {1, 2, 1, 2, 0.5, 1, 1, 2, 1, 0, 1, 0.5, 2, 1, 1, 1, 2},
    "Flying":   {1, 1, 1, 0.5, 2, 1, 2, 1, 1, 1, 1, 2, 0.5, 1, 1, 1, 0.5},
    "Psychic":  {1, 1, 1, 1, 1, 1, 2, 2, 1, 1, 0.5, 1, 1, 1, 1, 0, 0.5},
    "Bug":      {1, 0.5, 1, 1, 2, 1, 0.5, 0.5, 1, 0.5, 2, 1, 1, 0.5, 1, 2, 0.5},
    "Rock":     {1, 2, 1, 1, 1, 2, 0.5, 1, 0.5, 2, 1, 2, 1, 1, 1, 1, 0.5},
    "Ghost":    {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 2, 1, 0.5, 0.5},
    "Dragon":   {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 0.5},
    "Dark":     {1, 1, 1, 1, 1, 1, 0.5, 1, 1, 1, 2, 1, 1, 2, 1, 0.5, 0.5},
    "Steel":    {1, 0.5, 0.5, 0.5, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 0.5},
}

var MoveList = map[string]Move{
    "Strength": Move{"Strength", "Normal", "atk", 80, 100, 0, "None"},
    "Hyper Voice": Move{"Hyper Voice", "Normal", "spatk", 90, 100, 0, "None"},
    "Flamethrower": Move{"Flamethrower", "Fire", "spatk", 95, 100, 0.1, "BRN"},
    "Flame Wheel": Move{"Flame Wheel", "Fire", "atk", 60, 100, 0.1, "BRN"},
    "Surf": Move{"Surf", "Water", "spatk", 95, 100, 0, "None"},
    "Waterfall": Move{"Waterfall", "Water", "atk", 80, 100, 0.2, "flinch"},
    "Thunderbolt": Move{"Thunderbolt", "Electric", "spatk", 95, 100, 0.1, "PRZ"},
    "Spark": Move{"Spark", "Electric", "atk", 65, 100, 0.3, "PRZ"},
    "Energy Ball": Move{"Energy Ball", "Grass", "spatk", 80, 100, 0.1, "spdef"},
    "Seed Bomb": Move{"Seed Bomb", "Grass", "atk", 80, 100, 0, "None"},
    "Leaf Blade": Move{"Leaf Blade", "Grass", "atk", 90, 100, 2, "crit"},
    "Ice Beam": Move{"Ice Beam", "Ice", "spatk", 95, 100, 0.1, "FRZ"},
    "Avalanche": Move{"Avalanche", "Ice", "atk", 80, 100, 0, "None"},
    "Aura Sphere": Move{"Aura Sphere", "Fighting", "spatk", 90, 100, 0, "None"},
    "Brick Break": Move{"Brick Break", "Fighting", "atk", 75, 100, 0, "None"},
    "Sludge Bomb": Move{"Sludge Bomb", "Poison", "spatk", 90, 100, 0.3, "PSN"},
    "Poison Jab": Move{"Poison Jab", "Poison", "atk", 90, 100, 0.3, "PSN"},
    "Earthquake": Move{"Earthquake", "Ground", "atk", 100, 100, 0, "None"},
    "Earth Power": Move{"Earth Power", "Ground", "spatk", 90, 100, 10, "spdef"},
    "Air Slash": Move{"Air Slash", "Flying", "spatk", 75, 95, 0.3, "flinch"},
    "Aerial Ace": Move{"Aerial Ace", "Flying", "atk", 60, 100, 0, "None"},
    "Psychic": Move{"Psychic", "Psychic", "spatk", 90, 100, 0.1, "spdef"},
    "Psycho Cut": Move{"Psycho Cut", "Psychic", "atk", 90, 100, 0, "None"},
    "Extrasensory": Move{"Extrasensory", "Psychic", "spatk", 80, 100, 0.1, "flinch"},
    "X-scissor": Move{"X-scissor", "Bug", "atk", 80, 100, 0, "None"},
    "Bug Buzz": Move{"Bug Buzz", "Bug", "spatk", 90, 100, 0.1, "spdef"},
    "Rock Slide": Move{"Rock Slide", "Rock", "atk", 75, 90, 0.3, "flinch"},
    "Stone Edge": Move{"Stone Edge", "Rock", "atk", 100, 80, 2, "crit"},
    "Power Gem": Move{"Power Gem", "Rock", "spatk", 80, 100, 0, "None"},
    "Shadow Ball": Move{"Shadow Ball", "Ghost", "spatk", 80, 100, 0.2, "spdef"},
    "Shadow Punch": Move{"Shadow Punch", "Ghost", "atk", 60, 100, 0, "None"},
    "Dragon Pulse": Move{"Dragon Pulse", "Dragon", "spatk", 90, 100, 0, "None"},
    "Dragon Claw": Move{"Dragon Claw", "Dragon", "atk", 80, 100, 0, "None"},
    "Dark Pulse": Move{"Dark Pulse", "Dark", "spatk", 80, 100, 0.2, "flinch"},
    "Night Slash": Move{"Night Slash", "Dark", "atk", 70, 100, 2, "crit"},
    "Crunch": Move{"Crunch", "Dark", "atk", 80, 100, 0.2, "def"},
    "Iron Head": Move{"Iron Head", "Steel", "atk", 80, 100, 0.3, "flinch"},
    "Flash Cannon": Move{"Flash Cannon", "Steel", "spatk", 80, 100, 0.1, "spdef"},
}


// Move represents the structure of a move.
type Move struct {
    MoveName            string  `json:"name"`
    MoveType        string  `json:"type"`
    AtkType         string  `json:"atk_type"`
    Power           int     `json:"power"`
    Accuracy        int     `json:"accuracy"`
    SecondEffectRate float64 `json:"pp"`
    SecondEffect    string  `json:"description"`
}

var Moves map[string]Move


// the turn function for player
// returns a list of messages
// checks for status hindrance, accuracy
func AttackTurn(attacker *Pokemon, defender *Pokemon, move *Move) ([]string) {

	// check for volatile status
	canAttack, msg := CanAttackWithVolatileStatus(attacker) 
	if !canAttack {
		return []string{msg}
	}

	// check for non-volatile status
	canAttack, msg = CanAttackWithNonVolatileStatus(attacker) 
	if !canAttack {
		return []string{msg}
	}

	// initialize slice of messages to return
	messages := []string{attacker.Name + " used " + move.MoveName}

	// accuracy check
    attackLands, message := AccuracyCheck(attacker, move)
    if !attackLands {
        return append(messages, message)
    }

	// move lands, calculate damage dealt
	damage, msgs := DamageCalc(attacker, defender, move)
	messages = append(messages, msgs...)

	// apply damage to target
	if defender.hp <= damage {   // target faints
		messages = append(messages, defender.Name + " lost " + strconv.Itoa(defender.hp) + " health")
		defender.hp = 0
		defender.fainted = true
		messages = append(messages, defender.Name + " fainted!")
		return messages
	} else {   // target survives
		defender.hp -= damage
		messages = append(messages, defender.Name + " lost " + strconv.Itoa(damage) + " health (" + strconv.Itoa(defender.hp) + " hp left)")
	}
            
	// apply statuses
	msg = ApplyVolatileStatus(defender, move)
	if msg != ""{ 
		messages = append(messages, msg) 
	}
	msg = ApplyNonVolatileStatus(defender, move)
	if msg != "" { 
		messages = append(messages, msg) 
	}

	return messages
}

// returns a more complex message based on the outcome of the battle
func PostBattleMessage(winner *UserInput, loser *UserInput, wonBattle bool) {
	fmt.Println(winner.username, "defeated", loser.username)
	if wonBattle {   // you won the battle
		fmt.Println(WinMessages[rand.Intn(len(WinMessages))], "\n")
	} else {        // you lost the battle
		fmt.Println(LoseMessages[rand.Intn(len(LoseMessages))], "\n")
	}
}

// Printing helper function
func PrintMessages(msgs []string) {
	if len(msgs) == 0 {
		return
	}
	for _, x := range msgs {
		fmt.Println(x)
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
}

// runs a user's input after it has been collected for the turn
func HalfTurn(attackerInput *UserInput, defenderInput *UserInput) (*UserInput, *UserInput, bool) {
	var msg []string
	if attackerInput.action == "attack" {
		move := MoveList[attackerInput.move]
		msg = AttackTurn(attackerInput.activePokemon, defenderInput.activePokemon, &move)
	}
	PrintMessages(msg)

	// if defender faints, then send out a new Pokemon
	defenderLives := true
	if defenderInput.activePokemon.fainted {
		defenderLives= false
		if defenderInput.isAI {
			defenderInput = ReplaceFaintedPokemonAI(defenderInput, attackerInput)
		} else {
			defenderInput = ReplaceFaintedPokemon(defenderInput)
		}
	}
	return attackerInput, defenderInput, defenderLives
}

// wrapper function for everything that happens in one turn
// returns a bunch of messages for each player
func WholeTurn(userOneInput *UserInput, userTwoInput *UserInput) bool {

	var canAttack bool

	fmt.Println("[[ NEW TURN ]] What will you do?\n")

	// turn order
	userOneInput, userTwoInput = TurnOrder(userOneInput, userTwoInput)

	// TODO cannot flinch when moving first (probably should be refactored in future)
	if userOneInput.activePokemon.nonVolatileStatus == "flinch" { 
		userOneInput.activePokemon.nonVolatileStatus = "" 
	}

	// We want to collect the input from both users before anything happens in the turn
	// faster user attacks or switches (AI only attacks)
	if userOneInput.isAI {
		userOneInput = ChooseActionAI(userOneInput, userTwoInput)
	} else {
		userOneInput = ChooseAction(userOneInput)
	}

	// slower user attacks or switches (AI only attacks)
	if userTwoInput.isAI {
		userTwoInput = ChooseActionAI(userTwoInput, userOneInput)
	} else {
		userTwoInput = ChooseAction(userTwoInput)
	}

	// After collecting the input from both players, the turn can proceed
	// Start with faster pokemon attacking (or switching)
	userOneInput, userTwoInput, canAttack = HalfTurn(userOneInput, userTwoInput)

	// slower pokemon does not get a turn if fainted
	if userTwoInput.gameOver { 
		fmt.Println(userTwoInput.username, "is out of usable Pokemon...")
		fmt.Println(userTwoInput.username, "whited out!\n") 
		return true
	} else if !canAttack {
		return false
	}
	
	// slower pokemon can attack (or switch)
	userTwoInput, userOneInput, _ = HalfTurn(userTwoInput, userOneInput)
	
	if userOneInput.gameOver {
		fmt.Println(userOneInput.username, "is out of usable Pokemon...")
		fmt.Println(userOneInput.username, "whited out!\n") 
		return true
	} 
	return false
}

// wrapper function for a whole 6v6 singles battle
func Battle(userOneInput *UserInput, userTwoInput *UserInput) {
	fmt.Println("[[ BATTLE ]] Starting a battle\n")
	
	fmt.Println(userOneInput.username, "sent out", userOneInput.activePokemon.Name)
	fmt.Println(userTwoInput.username, "sent out", userTwoInput.activePokemon.Name, "\n")

	var gameOver bool
	for {
		gameOver = WholeTurn(userOneInput, userTwoInput)
		if gameOver {
			if userOneInput.gameOver {
				PostBattleMessage(userTwoInput, userOneInput, false)
			} else {
				PostBattleMessage(userOneInput, userTwoInput, true)
			}
			break
		}
	}
}


