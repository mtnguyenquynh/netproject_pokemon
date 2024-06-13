package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

// Assume these structs and methods are defined somewhere in your code
type UserInput struct {
    username      string
    action        string
    activePokemon *Pokemon
    team          []*Pokemon
    move          string
    isAI          bool
    gameOver      bool
}

var pokemonList map[string]PokemonData

type Pokemon struct {
    PokemonData
    
    level             int
    hp                int
    atk               int
    def               int
    spatk             int
    spdef             int
    speed             int
    nonVolatileStatus string
    volatileStatus    string
    fainted           bool
}

type PokemonData struct {
    PokedexNumber string   `json:"index"`
    Name          string   `json:"name"`
    Exp           int      `json:"exp"`
    BaseHP        int      `json:"hp"`
    BaseAtk       int      `json:"attack"`
    BaseDef       int      `json:"defense"`
    BaseSpAtk     int      `json:"sp_attack"`
    BaseSpDef     int      `json:"sp_defense"`
    BaseSpeed     int      `json:"speed"`
    TotalEVs      int      `json:"total_evs"`
    Type          [2]string `json:"type"`
    Description   string   `json:"description"`
    Height        string   `json:"height"`
    Weight        string   `json:"weight"`
    Level         int      `json:"level"`
    AccumExp      int      `json:"accum_exp"`
    Moves         []Move   `json:"moves"`
}


func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error connecting to server:", err.Error())
        return
    }
    defer conn.Close()

    reader := bufio.NewReader(conn)

    // Read welcome message
    welcomeMessage, _ := reader.ReadString('\n')
    fmt.Print(welcomeMessage)

    // Read Pokémon list
    pokemonList = make(map[string]PokemonData)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        line = strings.TrimSpace(line)
        if line == "" {
            break
        }
        pokemonNames := strings.Fields(line)
        for _, name := range pokemonNames {
            // For simplicity, we're only populating the names.
            // In a real scenario, you'd want to fully populate the PokemonData.
            pokemonList[name] = PokemonData{Name: name}
        }
    }

    // Debugging: Print out received Pokémon names
    fmt.Println("Received Pokémon List from Server:")
    for name := range pokemonList {
        fmt.Println(name)
    }

    // Example team setup using received Pokémon list
    venusaur := NewPokemon("Venusaur", true)
    charmeleon := NewPokemon("Charmeleon", true)
    wartortle := NewPokemon("Wartortle", true)
    blastoise := NewPokemon("Blastoise", true)
    caterpie := NewPokemon("Caterpie", true)
    bulbasaur := NewPokemon("Bulbasaur", true)

    cynthiasTeam := []*Pokemon{venusaur, charmeleon, wartortle, blastoise, caterpie, bulbasaur}

    myInput := &UserInput{"Ash", "", nil, nil, "", false, false}
    cynthiasInput := &UserInput{"Cynthia", "", venusaur, cynthiasTeam, "", true, false}

    fmt.Println()
    ChooseName(myInput)
    ChooseTeam(myInput)
    myInput.activePokemon = myInput.team[0]

    Battle(myInput, cynthiasInput)

}
