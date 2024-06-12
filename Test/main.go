package main

import (
	// "math/rand"
	// "time"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nfnt/resize"
)

// type Pokemon struct {
// 	Name     string
// 	Position rl.Vector2
// }

type Pokemon struct {
	Index       string       `json:"index"`
	Name        string       `json:"name"`
	Exp         int          `json:"exp"`
	HP          int          `json:"hp"`
	Attack      int          `json:"attack"`
	Defense     int          `json:"defense"`
	SpAttack    int          `json:"sp_attack"`
	SpDefense   int          `json:"sp_defense"`
	Speed       int          `json:"speed"`
	TotalEVs    int          `json:"total_evs"`
	Type        []string     `json:"type"`
	Description string       `json:"description"`
	Height      string       `json:"height"`
	Weight      string       `json:"weight"`
	ImageURL    string       `json:"image_url"`
	Level       int          `json:"level"`
	AccumExp    int          `json:"accum_exp"`
	Deployable  bool         `json:"deployable"`
	Texture     rl.Texture2D `json:"-"`
	Position    rl.Vector2   `json:"-"`
}

type Player struct {
	Name        string    `json:"name"`
	PokemonList []Pokemon `json:"pokemon_list"`
}

// type Player struct {
// 	Name     string
// 	Position rl.Vector2
// }

// type GameWorld struct {
// 	Players  []*Player
// 	Pokemons []*Pokemon
// }

// func NewGameWorld() *GameWorld {
// 	return &GameWorld{
// 		Players:  []*Player{},
// 		Pokemons: []*Pokemon{},
// 	}
// }

// func (w *GameWorld) SpawnPokemon() {
// 	pokemon := &Pokemon{
// 		Name: "Bulbasaur",
// 		Position: rl.Vector2{
// 			X: float32(rand.Intn(screenWidth)),
// 			Y: float32(rand.Intn(screenHeight)),
// 		},
// 	}
// 	w.Pokemons = append(w.Pokemons, pokemon)
// }

// func (w *GameWorld) AddPlayer(player *Player) {
// 	w.Players = append(w.Players, player)
// }

// func (w *GameWorld) MovePlayer(player *Player, dx, dy float32) {
// 	player.Position.X += dx
// 	player.Position.Y += dy
// }

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

	pokemons []*Pokemon
	player Player

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
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
		if pokemon == nil {
			continue // Skip if the Pokémon is nil
		}
		if isBrushTile(pokemon.Position.X, pokemon.Position.Y) {
			// Use the same tileSrc to draw Pokémon
			tileSrc.X = 0
			tileSrc.Y = 0
			rl.DrawTexturePro(pokemon.Texture, tileSrc, rl.NewRectangle(pokemon.Position.X, pokemon.Position.Y, 16, 16), rl.NewVector2(8, 8), 0, rl.White)
		}
	}

	rl.DrawTexturePro(playerSpire, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
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

// func update() {
// 	running = !rl.WindowShouldClose()

// 	playerSrc.X = playerSrc.Width * float32(playerFrame)

// 	if playerMoving {
// 		if playerUp {
// 			playerDest.Y -= playerSpeed
// 		}
// 		if playerDown {
// 			playerDest.Y += playerSpeed
// 		}
// 		if playerLeft {
// 			playerDest.X -= playerSpeed
// 		}
// 		if playerRight {
// 			playerDest.X += playerSpeed
// 		}
// 		if frameCount%8 == 1 {
// 			playerFrame++
// 		}
// 		// playerSrc.X = playerSrc.Width * float32(playerFrame)
// 	} else if frameCount%45 == 1 {
// 		playerFrame++
// 	}

// 	frameCount++
// 	if playerFrame > 3 {
// 		playerFrame = 0
// 	}
// 	if !playerMoving && playerFrame > 1 {
// 		playerFrame = 0
// 	}

// 	playerSrc.Y = playerSrc.Height * float32(playerDir)

// 	rl.UpdateMusicStream(music)
// 	if musicPaused {
// 		rl.PauseMusicStream(music)
// 	} else {
// 		rl.ResumeMusicStream(music)
// 	}

// 	cam.Target = rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2)))

// 	playerMoving = false
// 	playerUp, playerDown, playerRight, playerLeft = false, false, false, false
// }

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

	// Check if the new position is valid (on grass tile)
	if isGrassTile(newPlayerDest.X, newPlayerDest.Y) {
		playerDest = newPlayerDest
	}

    for i, pokemon := range pokemons {
        if pokemon == nil {
            continue
        }

        if rl.CheckCollisionRecs(playerDest, rl.NewRectangle(pokemon.Position.X, pokemon.Position.Y, 16, 16)) {
            player.PokemonList = append(player.PokemonList, *pokemon)
            pokemons[i] = nil
			fmt.Println(player)
			savePlayer(player)
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
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	rl.BeginMode2D(cam)

	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func loadMap(mapFile string) {
	file, err := ioutil.ReadFile(mapFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	remNewLines := strings.ReplaceAll(string(file), "\r", "")
	remNewLines = strings.ReplaceAll(remNewLines, "\n", " ")
	sliced := strings.Split(remNewLines, " ")
	mapW = -1
	mapH = -1
	for i := 0; i < len(sliced); i++ {
		s, _ := strconv.ParseInt(sliced[i], 10, 64)
		// if err != nil {
		// 	fmt.Printf("Error parsing string '%s' at index %d: %v\n", sliced[i], i, err)
		// 	continue
		// }
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
	fmt.Println("Parsed tileMap:", tileMap)
	fmt.Println("src tileMap:", srcMap)
	if len(tileMap) > mapW*mapH {
		tileMap = tileMap[:len(tileMap)-1]
	}

	// mapW = 5
	// mapH = 5
	// for i := 0; i<(mapW*mapH); i++ {
	// 	tileMap = append(tileMap, 13)
	// }
}

func init() {
	playerName, err := input_field()
	if err != nil {
		// Thoát chương trình nếu không nhập tên người chơi
		fmt.Println("Input field closed. Exiting...")
		os.Exit(1)
	}

	// Cập nhật tên người chơi trong player
	player.Name = playerName

	rl.InitWindow(screenWidth, screenHeight, "Pokémon Game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)


	// fmt.Print("Enter your name: ")
	// var playerName string
	// fmt.Scanln(&playerName)
	// fmt.Println("Hello,", playerName)

	// player.Name = playerName

	grassSprite = rl.LoadTexture("res/Tilesets/Grass.png")
	hillSprite = rl.LoadTexture("res/Tilesets/Hills.png")
	fenceSprite = rl.LoadTexture("res/Tilesets/Fences.png")
	houseSprite = rl.LoadTexture("res/Tilesets/Wooden House.png")
	waterSprite = rl.LoadTexture("res/Tilesets/Water.png")
	tilledSprite = rl.LoadTexture("res/Tilesets/Tilled Dirt.png")

	brushSprite = rl.LoadTexture("res/Objects/Basic_Grass_Biom_things.png")

	url := "https://archives.bulbagarden.net/media/upload/thumb/f/fb/0001Bulbasaur.png/70px-0001Bulbasaur.png"
	pokemonTexture, err := downloadTexture(url)
	if err != nil {
		fmt.Println("Error loading pokemon texture:", err)
		os.Exit(1)
	}

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerSpire = rl.LoadTexture("res/Characters/BasicCharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 60, 60)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("res/theme.mp3")
	musicPaused = true
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2))), 0.0, 1.5)

	cam.Zoom = 3
	loadMap("three.map")

	// Spawn Pokémon at valid positions
	brushPositions := getBrushTilePositions()
	pokemons = append(pokemons, spawnPokemon(pokemonTexture, brushPositions))


	squirtleTexture, err := downloadTexture("https://archives.bulbagarden.net/media/upload/thumb/5/54/0007Squirtle.png/70px-0007Squirtle.png")
	if err != nil {
		fmt.Println("Error loading squirtle texture:", err)
		os.Exit(1)
	}
	// Create a new Pokémon object for Squirtle
	squirtle := &Pokemon{
		Index:       "7",
		Name:        "Squirtle",
		Exp:         63,
		HP:          44,
		Attack:      48,
		Defense:     65,
		SpAttack:    50,
		SpDefense:   64,
		Speed:       43,
		TotalEVs:    314,
		Type:        []string{"water"},
		Description: "It shelters itself in its shell then strikes back with spouts of water at every opportunity.",
		Height:      "0.5 m",
		Weight:      "9 kg",
		ImageURL:    "https://archives.bulbagarden.net/media/upload/thumb/5/54/0007Squirtle.png/70px-0007Squirtle.png",
		Level:       0,
		AccumExp:    0,
		Deployable:  false,
		Texture:     squirtleTexture,
	}

	// Spawn Squirtle at a valid position
	brushPositions = getBrushTilePositions()
	squirtle.Position = brushPositions[rand.Intn(len(brushPositions))]

	// Add Squirtle to the pokemons list
	pokemons = append(pokemons, squirtle)

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

	// Trả về lỗi nếu cửa sổ đã đóng mà không nhập tên
	return "", errors.New("Input field closed without entering a name")
}

// func savePlayer(player Player) {
// 	// Encode player data into JSON format
// 	playerJSON, err := json.MarshalIndent(player, "", "    ")
// 	if err != nil {
// 		fmt.Println("Error encoding player data:", err)
// 		return
// 	}

// 	// Write JSON data to file
// 	file, err := os.Create("player.json")
// 	if err != nil {
// 		fmt.Println("Error creating player.json:", err)
// 		return
// 	}
// 	defer file.Close()

// 	_, err = file.Write(playerJSON)
// 	if err != nil {
// 		fmt.Println("Error writing to player.json:", err)
// 		return
// 	}

// 	fmt.Println("Player data saved to player.json")
// }

func savePlayer(player Player) {
	// Read existing player data from file
	file, err := os.OpenFile("player.json", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening player.json:", err)
		return
	}
	defer file.Close()

	// Decode existing player data from JSON
	var existingPlayers []Player
	err = json.NewDecoder(file).Decode(&existingPlayers)
	if err != nil {
		fmt.Println("Error decoding player data:", err)
		return
	}

	// Update existing player data or append new player data
	found := false
	for i, p := range existingPlayers {
		if p.Name == player.Name {
			existingPlayers[i] = player
			found = true
			break
		}
	}

	if !found {
		existingPlayers = append(existingPlayers, player)
	}

	// Encode updated player data into JSON format
	playerJSON, err := json.MarshalIndent(existingPlayers, "", "    ")
	if err != nil {
		fmt.Println("Error encoding player data:", err)
		return
	}

	// Seek to the beginning of the file to overwrite existing data
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("Error seeking player.json:", err)
		return
	}

	// Write updated JSON data back to the file
	_, err = file.Write(playerJSON)
	if err != nil {
		fmt.Println("Error writing to player.json:", err)
		return
	}

	// Truncate the file to remove any remaining data
	err = file.Truncate(int64(len(playerJSON)))
	if err != nil {
		fmt.Println("Error truncating player.json:", err)
		return
	}

	fmt.Println("Player data updated in player.json")
}

func getBrushTilePositions() []rl.Vector2 {
	positions := []rl.Vector2{}
	for i := 0; i < len(tileMap); i++ {
		if srcMap[i] == "a" {
			tileX := (i % mapW) * int(tileDest.Width)
			tileY := (i / mapW) * int(tileDest.Height)
			positions = append(positions, rl.NewVector2(float32(tileX), float32(tileY)))
		}
	}
	fmt.Println(positions)
	return positions
}
func spawnPokemon(texture rl.Texture2D, positions []rl.Vector2) *Pokemon {
	if len(positions) == 0 {
		return nil
	}
	index := rand.Intn(len(positions))
	position := positions[index]
	return &Pokemon{
		Texture:  texture,
		Position: position,
	}
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

