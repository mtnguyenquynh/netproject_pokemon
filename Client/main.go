package main

import (
	// "math/rand"
	// "time"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nfnt/resize"
)

type PokemonData struct {
	PokedexNumber string       `json:"index"`
	Name          string       `json:"name"`
	BaseHP        int          `json:"hp"`
	BaseAtk       int          `json:"attack"`
	BaseDef       int          `json:"defense"`
	BaseSpAtk     int          `json:"sp_attack"`
	BaseSpDef     int          `json:"sp_defense"`
	BaseSpeed     int          `json:"speed"`
	TotalEVs      int          `json:"total_evs"`
	Type          [2]string    `json:"type"`
	Description   string       `json:"description"`
	Height        string       `json:"height"`
	Weight        string       `json:"weight"`
	ImageURL      string       `json:"image_url"`
	Exp           int          `json:"exp"`
	Moves         []Move       `json:"moves"`
	Texture       rl.Texture2D `json:"-"`
	Position      rl.Vector2   `json:"-"`
	SpawnTime     time.Time    `json:"-"`
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

// type Player struct {
// 	Name        string    `json:"name"`
// 	PokemonList []Pokemon `json:"pokemon_list"`
// }

type Player struct {
	Name        string       `json:"name"`
	PokemonList []Pokemon    `json:"pokemon_list"`
	PlayerSrc   rl.Rectangle `json:"playerSrc"`
	PlayerDest  rl.Rectangle `json:"playerDest"`
}

const (
	screenWidth  = 1000
	screenHeight = 480
)

var (
	running  = true
	bkgColor = rl.NewColor(147, 211, 196, 255)

	grassSprite rl.Texture2D
	tex         rl.Texture2D

	playerSpire  rl.Texture2D
	hillSprite   rl.Texture2D
	fenceSprite  rl.Texture2D
	houseSprite  rl.Texture2D
	waterSprite  rl.Texture2D
	tilledSprite rl.Texture2D

	brushSprite rl.Texture2D

	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerRight, playerLeft bool
	playerFrame                                   int

	frameCount int

	tileDest   rl.Rectangle
	tileSrc    rl.Rectangle
	tileMap    []int
	srcMap     []string
	mapW, mapH int

	playerSpeed float32 = 3

	pokemons     []Pokemon
	pokemons_all []*Pokemon
	player       Player
	otherPlayers []Player

	brushTiles []rl.Vector2

	musicPaused bool
	music       rl.Music

	cam            rl.Camera2D
	conn           net.Conn
	updateInterval = 1 * time.Second
	mu           sync.Mutex
)

func drawScene() {
	// rl.DrawTexture(grassSprite, 100, 50, rl.White)

	for i := 0; i < len(tileMap); i++ {
		if tileMap[i] != 0 {
			tileDest.X = tileDest.Width * float32(i%mapW)
			tileDest.Y = tileDest.Height * float32(i/mapW)

			if srcMap[i] == "g" {
				tex = grassSprite
			}
			if srcMap[i] == "l" {
				tex = hillSprite
			}
			if srcMap[i] == "f" {
				tex = fenceSprite
			}
			if srcMap[i] == "h" {
				tex = houseSprite
			}
			if srcMap[i] == "w" {
				tex = waterSprite
			}
			if srcMap[i] == "t" {
				tex = tilledSprite
			}
			if srcMap[i] == "a" {
				tex = brushSprite
			}

			if srcMap[i] == "h" || srcMap[i] == "f" || srcMap[i] == "a" {
				tileSrc.X = tileSrc.Width * float32((56-1)%int(grassSprite.Width/int32(tileSrc.Width)))
				tileSrc.Y = tileSrc.Height * float32((56-1)/int(grassSprite.Width/int32(tileSrc.Width)))
				rl.DrawTexturePro(grassSprite, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
			}

			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))
			rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)

		}
	}

	for _, pokemon := range pokemons {
		if isBrushTile(pokemon.Position.X, pokemon.Position.Y) {
			// Assuming tileSrc and other necessary variables are properly set
			rl.DrawTexturePro(pokemon.Texture, tileSrc, rl.NewRectangle(pokemon.Position.X, pokemon.Position.Y, 16, 16), rl.NewVector2(8, 8), 0, rl.White)
		}
	}

	rl.DrawTexturePro(playerSpire, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)

	// Draw other players with semi-transparency
	for _, otherPlayer := range otherPlayers {
		if otherPlayer.Name != player.Name {
			otherPlayerDest := rl.NewRectangle(otherPlayer.PlayerDest.X, otherPlayer.PlayerDest.Y, playerDest.Width, playerDest.Height)
			rl.DrawTexturePro(playerSpire, otherPlayer.PlayerSrc, otherPlayerDest, rl.NewVector2(otherPlayerDest.Width, otherPlayerDest.Height), 0, rl.NewColor(255, 255, 255, 128))
			
		}

		
	}
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 1
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 0
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 3
		playerRight = true
	}
	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
	}
}

func update() {
	running = !rl.WindowShouldClose()

	playerSrc.X = playerSrc.Width * float32(playerFrame)

	newPlayerDest := playerDest

	if playerMoving {
		if playerUp {
			newPlayerDest.Y -= playerSpeed
		}
		if playerDown {
			newPlayerDest.Y += playerSpeed
		}
		if playerLeft {
			newPlayerDest.X -= playerSpeed
		}
		if playerRight {
			newPlayerDest.X += playerSpeed
		}

		if frameCount%8 == 1 {
			playerFrame++
		}
	} else if frameCount%45 == 1 {
		playerFrame++
	}

	if isGrassTile(newPlayerDest.X, newPlayerDest.Y) {
		playerDest = newPlayerDest
	}

	for i := 0; i < len(pokemons); i++ {
		pokemon := &pokemons[i] // Use pointer to the current pokemon
	
		if pokemon == nil {
			continue // Skip if the Pokémon is nil (though this should be rare in a normal slice)
		}
	
		pokemonRect := rl.NewRectangle(pokemon.Position.X, pokemon.Position.Y, 16, 16)
	
		if rl.CheckCollisionRecs(playerDest, pokemonRect) {
			player.PokemonList = append(player.PokemonList, *pokemon)
			// Remove the captured pokemon from pokemons slice
			pokemons = append(pokemons[:i], pokemons[i+1:]...)
			i-- // Decrement i because the next pokemon is shifted left
			//fmt.Println(player)
			savePlayer(player)
		}

		for _, otherPlayer := range otherPlayers {
			if otherPlayer.Name != player.Name {
				otherPlayerDest := rl.NewRectangle(otherPlayer.PlayerDest.X, otherPlayer.PlayerDest.Y, playerDest.Width, playerDest.Height)				
				if rl.CheckCollisionRecs(playerDest, otherPlayerDest) && isHouseTile(otherPlayerDest.X, otherPlayerDest.Y){
					battle(player.Name)
				}
			}
	
			
		}
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2)))

	playerMoving = false
	playerUp, playerDown, playerRight, playerLeft = false, false, false, false

	if time.Now().Second()%int(updateInterval.Seconds()) == 0 {
		player.PlayerSrc = playerSrc
		player.PlayerDest = playerDest
		updatePlayerPosition(player)

		// Fetch positions of other players from the server
		fetchAllPlayers()
	}

	if time.Now().Second() == 0 || len(pokemons) == 0 {
		spawnNewPokemon()
	}

	// Despawn Pokémon after 5 minutes
	despawnOldPokemon()
}

func spawnNewPokemon() {
	mu.Lock()
	defer mu.Unlock()

	// fetchPokemonData()

	if len(brushTiles) == 0 || len(pokemons_all) == 0 {
		return // No brush tiles or pokemons available
	}

	// Randomly pick a brush tile and a Pokémon
	randomIndex := rand.Intn(len(brushTiles))
	pokemonIndex := rand.Intn(len(pokemons_all))
	if pokemons_all[pokemonIndex] == nil {
		fmt.Println("Selected nil Pokémon, skipping spawn.")
		return // Skip if selected Pokémon is nil
	}

	newPokemon := *pokemons_all[pokemonIndex]
	newPokemon.Position = brushTiles[randomIndex]
	newPokemon.SpawnTime = time.Now()

	var err error
	newPokemon.Texture, err = downloadTexture(newPokemon.ImageURL)
	if err != nil {
		fmt.Println("Error downloading texture:", err)
		return
	}

	pokemons = append(pokemons, newPokemon)
}

func despawnOldPokemon() {
    currentTime := time.Now()
    // Create a new slice to store active Pokemon
    var activePokemons []Pokemon

    // Iterate over the pokemons slice
    for _, pokemon := range pokemons {
        // Check if the pokemon is active
        if pokemon.IsActive(currentTime) {
            activePokemons = append(activePokemons, pokemon)
        }
    }

    // Update the pokemons slice to contain only activePokemons
    pokemons = activePokemons
}

// IsActive checks if a pokemon is active based on current time and spawn time
func (p Pokemon) IsActive(currentTime time.Time) bool {
    return !p.SpawnTime.IsZero() && currentTime.Sub(p.SpawnTime) <= 5*time.Minute
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	rl.BeginMode2D(cam)

	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func loadMap(fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading map file", err)
		return
	}
	remNewLines := strings.ReplaceAll(string(file), "\r", "")
	remNewLines = strings.ReplaceAll(remNewLines, "\n", " ")
	sliced := strings.Split(remNewLines, " ")
	mapW = -1
	mapH = -1
	for i := 0; i < len(sliced); i++ {
		s, _ := strconv.ParseInt(sliced[i], 10, 64)
		m := int(s)
		if mapW == -1 {
			mapW = m
		} else if mapH == -1 {
			mapH = m
		} else if i < mapW*mapH+2 {
			tileMap = append(tileMap, m)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}
	if len(tileMap) > mapW*mapH {
		tileMap = tileMap[:len(tileMap)-1]
	}

	brushTiles = nil
	for i := 0; i < len(tileMap); i++ {
		if srcMap[i] == "a" { // 'a' represents brush tiles
			x := float32((i % mapW) * int(tileDest.Width))
			y := float32((i / mapW) * int(tileDest.Height))
			brushTiles = append(brushTiles, rl.NewVector2(x, y))
		}
	}

	fmt.Println("Parsed tileMap:", tileMap)
	fmt.Println("Parsed srcMap:", srcMap)
	fmt.Println("Brush tiles:", brushTiles)
}

func init() {
	playerName, err := input_field()
	if err != nil {
		fmt.Println("Input field closed. Exiting...")
		os.Exit(1)
	}

	player.Name = playerName

	conn, err = net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	// defer conn.Close()
	fmt.Println("Success connect to server")

	rl.InitWindow(screenWidth, screenHeight, "Pokémon Game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("res/Tilesets/Grass.png")
	hillSprite = rl.LoadTexture("res/Tilesets/Hills.png")
	fenceSprite = rl.LoadTexture("res/Tilesets/Fences.png")
	houseSprite = rl.LoadTexture("res/Tilesets/Wooden House.png")
	waterSprite = rl.LoadTexture("res/Tilesets/Water.png")
	tilledSprite = rl.LoadTexture("res/Tilesets/Tilled Dirt.png")

	brushSprite = rl.LoadTexture("res/Objects/Basic_Grass_Biom_things.png")

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerSpire = rl.LoadTexture("res/Characters/BasicCharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 60, 60)

	player, err = fetchPlayer(playerName)
	if err != nil {
		fmt.Println("Error fetching player data:", err)
		return
	}

	// Initialize playerSrc and playerDest
	// playerSrc = player.PlayerSrc
	// playerDest = player.PlayerDest

	pokemons_all, err = fetchPokemonData()
	if err != nil {
		fmt.Println("Error fetching pokemons:", err)
		return
	}

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("res/theme.mp3")
	musicPaused = true
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2))), 0.0, 1.5)

	cam.Zoom = 3
	loadMap("three.map")

	loadPokemonImages()
	updatePokemonPositions()

	// spawnPokemonsOnMap()

}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSpire)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func main() {

	for running {
		input()
		update()
		render()
	}
	quit()
}



func updatePlayerPosition(player Player) error {
	writer := bufio.NewWriter(conn)

	playerData, err := json.Marshal(player)
	if err != nil {
		return err
	}

	fmt.Fprintln(writer, "UPDATE_POSITION")
	writer.Flush()
	fmt.Fprintln(writer, string(playerData))
	writer.Flush()

	return nil
}

func updatePokemonPositions() {
	// Shuffle the pokemons_all slice
	rand.Shuffle(len(pokemons_all), func(i, j int) {
		pokemons_all[i], pokemons_all[j] = pokemons_all[j], pokemons_all[i]
	})

	// Ensure we only place a maximum of three Pokémon
	maxPokemons := 2
	if len(pokemons_all) > maxPokemons {
		pokemons = make([]Pokemon, maxPokemons) // Create a new slice of Pokemon structs
		for i := 0; i < maxPokemons; i++ {
			pokemons[i] = *pokemons_all[i] // Dereference and copy each Pokemon struct
		}
	} else {
		pokemons = make([]Pokemon, len(pokemons_all)) // Create a new slice of Pokemon structs
		for i := 0; i < len(pokemons_all); i++ {
			pokemons[i] = *pokemons_all[i] // Dereference and copy each Pokemon struct
		}
	}

	// Assign positions to the pokemons slice
	for i := 0; i < len(pokemons); i++ {
		// Randomly pick a brush tile index
		randomIndex := rand.Intn(len(brushTiles))
		pokemons[i].Position = brushTiles[randomIndex]
	}
}

func input_field() (string, error) {
	rl.InitWindow(800, 450, "Text Input Example")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	const maxInputChars = 20

	var inputText string

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyBackspace) && len(inputText) > 0 {
			inputText = inputText[:len(inputText)-1]
		} else {
			key := rl.GetCharPressed()
			if key >= 32 && key <= 126 && len(inputText) < maxInputChars {
				inputText += string(rune(key))
			}
		}

		if rl.IsKeyPressed(rl.KeyEnter) && len(inputText) > 0 {
			return inputText, nil
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawRectangle(10, 50, 400, 40, rl.LightGray)
		rl.DrawRectangleLines(10, 50, 400, 40, rl.Black)
		rl.DrawText("Enter Text: "+inputText, 20, 60, 20, rl.DarkGray)
		rl.DrawText("Type to input text. Backspace to delete. Press Enter to save.", 10, 10, 20, rl.DarkGray)

		rl.EndDrawing()
	}

	return "", errors.New("Input field closed without entering a name")
}

func savePlayer(player Player) error {
	writer := bufio.NewWriter(conn)

	playerData, err := json.Marshal(player)
	if err != nil {
		return err
	}

	fmt.Fprintln(writer, "SAVE_PLAYER")
	writer.Flush()
	fmt.Fprintln(writer, string(playerData))
	writer.Flush()

	return nil
}

func fetchPlayer(name string) (Player, error) {

	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	reader := bufio.NewReader(conn)

	// Send GET_PLAYER command
	fmt.Fprintf(writer, "GET_PLAYER\n")

	writer.Flush()
	// Send player's name
	fmt.Fprintf(writer, "%s\n", name)

	writer.Flush()

	// Read server response
	response, err := reader.ReadString('\n')
	if err != nil {
		return Player{}, fmt.Errorf("failed to read server response: %v", err)
	}

	// Unmarshal response into Player struct
	var player Player
	if err := json.Unmarshal([]byte(response), &player); err != nil {
		return Player{}, fmt.Errorf("failed to unmarshal player data: %v", err)
	}
	return player, nil
}

func fetchPokemonData() ([]*Pokemon, error) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	fmt.Fprintln(writer, "GET_POKEMON")
	writer.Flush()

	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	var pokemons []*Pokemon
	err = json.Unmarshal([]byte(response), &pokemons)
	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

func fetchAllPlayers() {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	fmt.Fprintln(writer, "GET_ALL_PLAYERS")
	writer.Flush()

	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal([]byte(response), &otherPlayers)

}

func downloadTexture(url string) (rl.Texture2D, error) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	img, _, _ := image.Decode(resp.Body)

	resizedImg := resize.Resize(uint(16), uint(16), img, resize.Lanczos3)

	rlImg := rl.NewImageFromImage(resizedImg)
	texture := rl.LoadTextureFromImage(rlImg)

	return texture, nil
}
func loadPokemonImages() {
	for i := 0; i < len(pokemons_all); i++ {
		pokemons_all[i].Texture, _ = downloadTexture(pokemons_all[i].ImageURL)
		if pokemons_all[i].Texture.ID <= 0 {
			fmt.Println("Error loading pokemon texture", pokemons_all[i].ImageURL)
		}
	}
}

func isGrassTile(x, y float32) bool {
	// Calculate the tile index
	tileX := int(x / tileDest.Width)
	tileY := int(y / tileDest.Height)
	tileIndex := tileY*mapW + tileX

	// Check if the tileIndex is within the bounds of the tileMap
	if tileIndex < 0 || tileIndex >= len(tileMap) {
		return false
	}
	if tileMap[tileIndex] == 12 || tileMap[tileIndex] == 2 {
		return false
	}
	if srcMap[tileIndex] == "h" && (tileMap[tileIndex] == 11 || tileMap[tileIndex] == 16) {
		return true
	}

	// Return true if the tile is a grass tile
	return srcMap[tileIndex] == "g" || srcMap[tileIndex] == "a"
}

func isBrushTile(x, y float32) bool {
	// Calculate the tile index
	tileX := int(x / tileDest.Width)
	tileY := int(y / tileDest.Height)
	tileIndex := tileY*mapW + tileX

	// Check if the tileIndex is within the bounds of the tileMap
	if tileIndex < 0 || tileIndex >= len(tileMap) {
		return false
	}

	// Return true if the tile is a brush tile
	return srcMap[tileIndex] == "a"
}

func isHouseTile(x, y float32) bool {
	// Calculate the tile index
	tileX := int(x / tileDest.Width)
	tileY := int(y / tileDest.Height)
	tileIndex := tileY*mapW + tileX

	// Check if the tileIndex is within the bounds of the tileMap
	if tileIndex < 0 || tileIndex >= len(tileMap) {
		return false
	}

	// Return true if the tile is a brush tile
	return srcMap[tileIndex] == "h" && tileMap[tileIndex] == 16
}