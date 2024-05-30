/* package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Player struct {
	ID          string         `json:"id"`
	PokemonList []PlayerPokemon `json:"pokemon_list"`
	X           int            `json:"x"`
	Y           int            `json:"y"`
}

type PlayerPokemon struct {
	Pokemon    Pokemon `json:"pokemon"`
	Level      int     `json:"level"`
	Experience int     `json:"experience"`
}

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
}

func loadPlayer(playerID string) (*Player, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/load_player?id=%s", playerID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var player Player
	err = json.NewDecoder(resp.Body).Decode(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func movePlayer(playerID, direction string) error {
	resp, err := http.PostForm(fmt.Sprintf("http://localhost:8080/move"), map[string][]string{
		"player_id": {playerID},
		"direction": {direction},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("move failed: %s", body)
	}
	return nil
}

func catchPokemon(playerID string) error {
	resp, err := http.PostForm(fmt.Sprintf("http://localhost:8080/poke_cat"), map[string][]string{
		"player_id": {playerID},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("catch failed: %s", body)
	}
	return nil
}

func startBattle(playerID string) error {
	resp, err := http.PostForm(fmt.Sprintf("http://localhost:8080/start_battle"), map[string][]string{
		"player_id": {playerID},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Battle response: %s\n", body)
	return nil
}

func main() {
	playerID := "player1"
	if len(os.Args) > 1 {
		playerID = os.Args[1]
	}

	player, err := loadPlayer(playerID)
	if err != nil {
		fmt.Printf("Error loading player: %v\n", err)
		return
	}

	fmt.Printf("Player loaded: %v\n", player)

	var input string
	for {
		fmt.Println("1. Catch Pokémon")
		fmt.Println("2. Start Battle")
		fmt.Print("Choose an option: ")
		fmt.Scanln(&input)

		switch input {
		case "1":
			fmt.Println("Catching Pokémon...")
			for i := 0; i < 120; i++ {
				err := movePlayer(playerID, "up")
				if err != nil {
					fmt.Printf("Error moving player: %v\n", err)
				}

				err = catchPokemon(playerID)
				if err != nil {
					fmt.Printf("Error catching Pokémon: %v\n", err)
				}

				time.Sleep(1 * time.Second)
			}
		case "2":
			fmt.Println("Starting battle...")
			err := startBattle(playerID)
			if err != nil {
				fmt.Printf("Error starting battle: %v\n", err)
			}
		default:
			fmt.Println("Invalid option")
		}
	}
}
 */


// client.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const serverURL = "http://localhost:8080"

type Player struct {
	ID          string
	PokemonList []PlayerPokemon
	X           int
	Y           int
}

type PlayerPokemon struct {
	Pokemon    Pokemon
	Level      int
	Experience int
}

type Pokemon struct {
	Index       string
	Name        string
	Exp         int
	HP          int
	Attack      int
	Defense     int
	SpAttack    int
	SpDefense   int
	Speed       int
	TotalEVs    int
	Type        []string
	Description string
	Height      string
	Weight      string
	ImageURL    string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your player ID: ")
	playerID, _ := reader.ReadString('\n')
	playerID = strings.TrimSpace(playerID)

	player, err := loadPlayer(playerID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Welcome, %s! You are currently at position (%d, %d).\n", player.ID, player.X, player.Y)
	fmt.Println("You have the following Pokémon:")
	for _, p := range player.PokemonList {
		fmt.Printf("- %s (Level %d, %d EXP)\n", p.Pokemon.Name, p.Level, p.Experience)
	}

	for {
		fmt.Println("Options: (1) Catch Pokémon, (2) Move, (3) Battle, (4) Exit")
		fmt.Print("Choose an option: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			fmt.Print("Enter X coordinate to catch Pokémon: ")
			xStr, _ := reader.ReadString('\n')
			x, _ := strconv.Atoi(strings.TrimSpace(xStr))

			fmt.Print("Enter Y coordinate to catch Pokémon: ")
			yStr, _ := reader.ReadString('\n')
			y, _ := strconv.Atoi(strings.TrimSpace(yStr))

			err = catchPokemon(playerID, x, y)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Caught Pokémon!")
			}

		case "2":
			fmt.Println("Move (up, down, left, right): ")
			direction, _ := reader.ReadString('\n')
			direction = strings.TrimSpace(direction)

			err = movePlayer(playerID, direction)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Moved!")
			}

		case "3":
			fmt.Println("Entering battle mode...")
			// Implement battle mode
		case "4":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func loadPlayer(playerID string) (Player, error) {
	resp, err := http.Get(serverURL + "/loadPlayer?id=" + playerID)
	if err != nil {
		return Player{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Creating new player...")
	}

	var player Player
	err = json.NewDecoder(resp.Body).Decode(&player)
	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func catchPokemon(playerID string, x, y int) error {
	resp, err := http.PostForm(serverURL+"/catchPokemon",
		url.Values{
			"playerID": {playerID},
			"x":        {strconv.Itoa(x)},
			"y":        {strconv.Itoa(y)},
		})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to catch Pokémon")
	}

	return nil
}

func movePlayer(playerID, direction string) error {
	resp, err := http.PostForm(serverURL+"/move",
		url.Values{
			"playerID":  {playerID},
			"direction": {direction},
		})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to move player")
	}

	return nil
}
