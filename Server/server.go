package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

// type Pokemon struct {
// 	Index       string       `json:"index"`
// 	Name        string       `json:"name"`
// 	HP          int          `json:"hp"`
// 	Attack      int          `json:"attack"`
// 	Defense     int          `json:"defense"`
// 	SpAttack    int          `json:"sp_attack"`
// 	SpDefense   int          `json:"sp_defense"`
// 	Speed       int          `json:"speed"`
// 	TotalEVs    int          `json:"total_evs"`
// 	Type        []string     `json:"type"`
// 	Description string       `json:"description"`
// 	Height      string       `json:"height"`
// 	Weight      string       `json:"weight"`
// 	ImageURL    string       `json:"image_url"`
// 	Exp         int          `json:"exp"`
// 	Moves       []Move       `json:"moves"`
// 	Texture     rl.Texture2D `json:"-"`
// 	Position    rl.Vector2   `json:"-"`
// }

// Move represents the structure of a move.
// type Move struct {
// 	Name        string  `json:"name"`
// 	Type        string  `json:"type"`
// 	AtkType     string  `json:"atk_type"`
// 	Power       int     `json:"power"`
// 	Accuracy    float64 `json:"accuracy"`
// 	PP          float64 `json:"pp"`
// 	Description string  `json:"description"`
// }

type Player struct {
	Name        string    `json:"name"`
	PokemonList []Pokemon `json:"pokemon_list"`
	Position    struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"position"`
}

var (
	playerDataMutex sync.Mutex
	players         []Player
	playersMutex    sync.Mutex
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	writer := bufio.NewWriter(conn)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		message = message[:len(message)-1] // Remove newline character
		fmt.Println(message)
		switch {
		case message == "GET_PLAYER":
			name, _ := bufio.NewReader(conn).ReadString('\n')
			name = name[:len(name)-1] // Remove newline character
			getPlayer(writer, name)

		case message == "SAVE_PLAYER":
			playerData, _ := bufio.NewReader(conn).ReadString('\n')
			playerData = playerData[:len(playerData)-1] // Remove newline character
			savePlayer(writer, playerData)

		case message == "GET_POKEMON":
			getPokemon(writer)

		case message == "UPDATE_POSITION":
			playerData, _ := bufio.NewReader(conn).ReadString('\n')
			playerData = playerData[:len(playerData)-1] // Remove newline character
			updatePlayerPosition(writer, playerData)

		case message == "GET_ALL_PLAYERS":
			getAllPlayers(writer)

		// case message == "BATTLE":
		// 	game()

		default:
			fmt.Fprintf(writer, "Unknown command\n")
			writer.Flush()
		}
	}
}

func getPlayer(writer *bufio.Writer, name string) {
	file, err := os.OpenFile("./crawler/player.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintln(writer, "Error opening player.json")
		writer.Flush()
		return
	}
	defer file.Close()

	var players []Player
	err = json.NewDecoder(file).Decode(&players)
	if err != nil && err != io.EOF {
		fmt.Fprintln(writer, "Error decoding player data:", err.Error())
		writer.Flush()
		return
	}

	for _, player := range players {
		if player.Name == name {
			response, _ := json.Marshal(player)
			fmt.Fprintln(writer, string(response))
			writer.Flush()
			return
		}
	}

	// Create a new player if not found
	newPlayer := Player{Name: name}
	players = append(players, newPlayer)

	// Move the file pointer to the beginning to overwrite the file
	file.Seek(0, 0)
	file.Truncate(0)

	err = json.NewEncoder(file).Encode(players)
	if err != nil {
		fmt.Fprintln(writer, "Error saving new player data:", err.Error())
		writer.Flush()
		return
	}

	response, _ := json.Marshal(newPlayer)
	fmt.Fprintln(writer, string(response))
	writer.Flush()
}

func savePlayer(writer *bufio.Writer, playerData string) {
	playerDataMutex.Lock()
	defer playerDataMutex.Unlock()

	var player Player
	err := json.Unmarshal([]byte(playerData), &player)
	if err != nil {
		fmt.Fprintln(writer, "Invalid player data")
		writer.Flush()
		return
	}

	file, err := os.OpenFile("./crawler/player.json", os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintln(writer, "Error opening player.json")
		writer.Flush()
		return
	}
	defer file.Close()

	var players []Player
	err = json.NewDecoder(file).Decode(&players)
	if err != nil {
		fmt.Fprintln(writer, "Error decoding player data")
		writer.Flush()
		return
	}

	found := false
	for i, p := range players {
		if p.Name == player.Name {
			players[i] = player
			found = true
			break
		}
	}

	if !found {
		players = append(players, player)
	}
	playerJSON, err := json.MarshalIndent(players, "", "    ")
	if err != nil {
		fmt.Fprintln(writer, "Error encoding player data")
		writer.Flush()
		return
	}

	file.Seek(0, 0)
	file.Write(playerJSON)
	file.Truncate(int64(len(playerJSON)))

	// fmt.Fprintln(writer, "Player data updated in player.json")
	writer.Flush()
}

func getPokemon(writer *bufio.Writer) {
	file, err := os.Open("./crawler/pokedex.json")
	if err != nil {
		fmt.Fprintln(writer, "Error opening pokedex.json")
		writer.Flush()
		return
	}
	defer file.Close()

	var pokemons []Pokemon
	err = json.NewDecoder(file).Decode(&pokemons)
	if err != nil {
		fmt.Fprintln(writer, "Error decoding pokedex data")
		writer.Flush()
		return
	}

	response, _ := json.Marshal(pokemons)
	fmt.Fprintln(writer, string(response))
	writer.Flush()
}

func updatePlayerPosition(writer *bufio.Writer, playerData string) {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	var updatedPlayer Player
	err := json.Unmarshal([]byte(playerData), &updatedPlayer)
	if err != nil {
		fmt.Fprintln(writer, "Invalid player data")
		writer.Flush()
		return
	}

	found := false
	for i, player := range players {
		if player.Name == updatedPlayer.Name {
			players[i].Position = updatedPlayer.Position
			found = true
			break
		}
	}

	if !found {
		players = append(players, updatedPlayer)
	}

	fmt.Fprintln(writer, "Player position updated")
	writer.Flush()
}

func getAllPlayers(writer *bufio.Writer) {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	response, _ := json.Marshal(players)
	fmt.Fprintln(writer, string(response))
	writer.Flush()
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func game() {
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
