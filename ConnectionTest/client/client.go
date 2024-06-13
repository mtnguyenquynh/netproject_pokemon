package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

type UserInput struct {
    username string
    action string
    activePokemon *Pokemon
    team []*Pokemon
    move string
    isAI bool
    gameOver bool
}


// Define structs to hold the data from pokedex.json
type PokemonData struct {
    PokedexNumber  string   `json:"index"`
    Name           string   `json:"name"`
    Exp            int      `json:"exp"`
    BaseHP         int      `json:"hp"`
    BaseAtk        int      `json:"attack"`
    BaseDef        int      `json:"defense"`
    BaseSpAtk      int      `json:"sp_attack"`
    BaseSpDef      int      `json:"sp_defense"`
    BaseSpeed      int      `json:"speed"`
    TotalEVs       int      `json:"total_evs"`
    Type           [2]string `json:"type"`
    Description    string   `json:"description"`
    Height         string   `json:"height"`
    Weight         string   `json:"weight"`
    Level          int      `json:"level"`
    AccumExp       int      `json:"accum_exp"`
    Moves         []Move  `json:"moves"`
}

type Pokemon struct {
    // Shared among individuals
    PokemonData

    // Specific per individual
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

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error connecting:", err.Error())
        return
    }
    defer conn.Close()

    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(conn)

    // Read server prompts
    fmt.Print(readMessage(conn))

    // Input user's name
    fmt.Print("Enter your name: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)
    fmt.Fprintf(conn, name+"\n")
    writer.Flush()

    // Receive and handle messages from the server
    for {
        message := readMessage(conn)
        if message == "" {
            break
        }

        if strings.Contains(message, "Type a number and press ENTER to choose an option") {
            // Handle the team selection prompt
            fmt.Print(message)
            number, _ := reader.ReadString('\n')
            fmt.Fprintf(conn, number+"\n")
            writer.Flush()
            } else if strings.Contains(message, "Choose your (") {
             
            // Read and display all Pokémon names sent by the server
            for {
                additionalMessage := readMessage(conn)
                fmt.Print(additionalMessage)
                if strings.Contains(additionalMessage, "Choose your (") {
                    break
                }
            }

            pokemon, _ := reader.ReadString('\n')
            pokemon = strings.TrimSpace(pokemon) // Trim the newline character

            writer.Flush()
        } else {
            // Print other messages from the server
            fmt.Print(message)
        }
    }
}

func readMessage(conn net.Conn) string {
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        return ""
    }
    return message
}