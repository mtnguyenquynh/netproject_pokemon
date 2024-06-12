package main
import (
	"net"
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
	Level       int      `json:"level"`
	AccumExp    int      `json:"accum_exp"`
	Deployable  bool     `json:"deployable"`
	MoveList   [4]string
	Fainted    bool
	NonVolatileStatus string
	VolatileStatus    string
}

type Player struct {
	Name        string    `json:"name"`
	PokemonList []Pokemon `json:"pokemon_list"`
}

type Participant struct {
	player     Player
	turn       int
	isWin      bool
	curPokemon Pokemon
	conn       net.Conn
}
type Message struct {
	msg  string
	conn net.Conn
}

var MoveList = map[string]Move{
	"Strength":    {"Strength", "Normal", "atk", 80, 100, 0, "None"},
	"Hyper Voice": {"Hyper Voice", "Normal", "spatk", 90, 100, 0, "None"},
	// Add more moves as needed
}

// Move represents a Pokemon move
type Move struct {
	Name             string
	MoveType         string
	AtkType          string
	Power            int
	Accuracy         int
	SecondEffectRate float64
	SecondEffect     string
}