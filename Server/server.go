package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	participants []Participant
	conns        []net.Conn
	connCh       = make(chan net.Conn)
	closeCh      = make(chan Participant)
	msgCh        = make(chan string)
	msgChOne     = make(chan Message)
	starters     = []string{"Charmander", "Bulbasaur", "Squirtle"}
	mu           sync.Mutex
	writeMu      sync.Mutex
)

func main() {
	server, err := net.Listen("tcp", ":3012")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server started")
	// Load the Pokedex
	file, _ := os.Open("./crawler/pokedex.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	pokedex := []Pokemon{}
	_ = decoder.Decode(&pokedex)
	fmt.Println("Pokedex loaded")

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Fatal(err)
			}
			mu.Lock()
			conns = append(conns, conn)
			mu.Unlock()
			connCh <- conn
		}
	}()
	go func() {
		for {
			if len(participants) == 2 {
				winner, loser := battle(&participants[0], &participants[1])
				saveWinner(winner.player)
				msg := fmt.Sprintf("\n%s wins the battle - %s lost\n", winner.player.Name, loser.player.Name)
				msgCh <- msg + "#"
				// remove all connections
				for _, p := range participants {
					closeCh <- p
				}
			}
		}
	}()
	for {
		select {
		case conn := <-connCh:
			go onMessage(conn, pokedex)

		case msg := <-msgCh:
			fmt.Print(msg)
			publishMsgAll(msg)

		case participant := <-closeCh:
			fmt.Printf("%s exit\n", participant.player.Name)
			removeParticipant(participant)
		case msg := <-msgChOne:
			fmt.Print(msg.msg)
			publishMsgOne(msg.conn, msg.msg)
		}
	}

}

func removeParticipant(participant Participant) {

	for i := range participants {
		if participants[i].player.Name == participant.player.Name {
			participants = append(participants[:i], participants[i+1:]...)
			break
		}
	}
	// Remove from conns
	for i, conn := range conns {
		if conn == participant.conn {
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}


func publishMsgOne(conn net.Conn, msg string) error {
	writeMu.Lock()
	defer writeMu.Unlock()

	msgBytes := []byte(msg)
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(msgBytes)))

	if _, err := conn.Write(lenBytes); err != nil {
		return err
	}
	if _, err := conn.Write(msgBytes); err != nil {
		return err
	}
	return nil
}

func publishMsgAll(msg string) error {
	writeMu.Lock()
	defer writeMu.Unlock()

	for i := range conns {
		msgBytes := []byte(msg)
		lenBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBytes, uint32(len(msgBytes)))

		if _, err := conns[i].Write(lenBytes); err != nil {
			return err
		}
		if _, err := conns[i].Write(msgBytes); err != nil {
			return err
		}
	}
	return nil
}

func onMessage(conn net.Conn, pokedex []Pokemon) {
	fmt.Println("new client")
	reader := bufio.NewReader(conn)
	playerName, err := reader.ReadString('\n')
	playerName = strings.TrimSpace(playerName)
	if err != nil {
		return
	}
	fmt.Println(playerName)
	player, found := findPlayer(playerName)
	if !found {
		publishMsgOne(conn, "Player does not exist. Created a new player.\n#")
		player = createPlayer(pokedex, playerName)
	}
	// request the player to choose a Pokemon
	// make all the Pokemon deployable
	for i := range player.PokemonList {
		player.PokemonList[i].Deployable = true
	}
	msg := listPokemon(player.PokemonList)

	chosenPokemon, _ := readPokemonFromClient(conn, msg[:len(msg)-1]+"Choose a pokemon: #", player.PokemonList)

	// Add the player to the list of participants
	mu.Lock()
	participants = append(participants, Participant{
		player:     player,
		turn:       3,
		isWin:      false,
		curPokemon: chosenPokemon,
		conn:       conn,
	})
	mu.Unlock()
	fmt.Println(len(participants))

}


func calculateDamage(attacker, defender *Pokemon, attackType string) int {
	if attackType == "normal" {
		return max(attacker.Attack-defender.Defense, 0)
	}
	if attackType == "special" {
		elementalMultiplier := 1.75
		return max(int(float64(attacker.SpAttack)*elementalMultiplier)-defender.SpDefense, 0)
	}
	return 0
}

// Helper function to get the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func saveWinner(player Player) {
	// Load the existing players from the JSON file
	file, _ := os.Open("./crawler/players.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	existingPlayers := []Player{}
	_ = decoder.Decode(&existingPlayers)

	// Find the winner in the list and update their data
	for i, existingPlayer := range existingPlayers {
		if existingPlayer.Name == player.Name {
			existingPlayers[i] = player
			break
		}
	}

	// Save the updated list of players to the JSON file
	file, _ = os.Create("./crawler/players.json")
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Set indent to 4 spaces
	_ = encoder.Encode(existingPlayers)

	fmt.Printf("Winner %s saved\n", player.Name)
}