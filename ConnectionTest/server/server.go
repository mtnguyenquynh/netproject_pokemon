package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net"
    "os"
)

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

type Move struct {
    MoveName        string  `json:"name"`
    MoveType        string  `json:"type"`
    AtkType         string  `json:"atk_type"`
    Power           int     `json:"power"`
    Accuracy        int     `json:"accuracy"`
    SecondEffectRate float64 `json:"pp"`
    SecondEffect    string  `json:"description"`
}

func InitData() {
    pokemonList = make(map[string]PokemonData)
    file, err := os.Open("./pokedex.json")
    if err != nil {
        fmt.Println("Error opening pokedex file:", err)
        os.Exit(1)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    pokedex := []PokemonData{}
    if err := decoder.Decode(&pokedex); err != nil {
        fmt.Println("Error decoding pokedex data:", err)
        os.Exit(1)
    }

    for _, pokemon := range pokedex {
        pokemonList[pokemon.Name] = pokemon
    }

    fmt.Println("PokemonList loaded")

    for _, pokemon := range pokemonList {
        fmt.Println(pokemon)
    }
}

func main() {
    InitData()

    listener, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error starting server:", err.Error())
        return
    }
    defer listener.Close()
    fmt.Println("Server is listening on localhost:8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err.Error())
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    writer := bufio.NewWriter(conn)

    fmt.Fprintln(writer, "Welcome to the Pokémon Team Builder!")
    writer.Flush()

    sendPokemonList(writer)
}

func sendPokemonList(writer *bufio.Writer) {
    pokemonListJSON, err := json.Marshal(pokemonList)
    if err != nil {
        fmt.Println("Error marshaling pokemonList:", err)
        return
    }

    writer.WriteString(string(pokemonListJSON) + "\n")
    writer.Flush()
}
