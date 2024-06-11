package main

import (
	// "math/rand"
	// "time"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
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
	Texture  rl.Texture2D
	Position rl.Vector2
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
		rl.DrawTexturePro(pokemon.Texture, tileSrc, rl.NewRectangle(pokemon.Position.X, pokemon.Position.Y, 12, 12), rl.NewVector2(48, 48), 0, rl.White)
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

func update() {
	running = !rl.WindowShouldClose()

	playerSrc.X = playerSrc.Width * float32(playerFrame)

	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}
		if frameCount%8 == 1 {
			playerFrame++
		}
		// playerSrc.X = playerSrc.Width * float32(playerFrame)
	} else if frameCount%45 == 1 {
		playerFrame++
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

	url := "https://archives.bulbagarden.net/media/upload/thumb/f/fb/0001Bulbasaur.png/70px-0001Bulbasaur.png"
	pokemonTexture, err := downloadTexture(url)
	//pokemonTexture, err := rl.LoadImageFromTexture(rl.LoadImageFromMemoryBase64(pokemonImageBase64))
	if err != nil {
		fmt.Println("Error loading pokemon texture:", err)
		os.Exit(1)
	}

	// Create a Pokémon object with the loaded texture
	pokemons = append(pokemons, &Pokemon{
		Texture:  pokemonTexture,
		Position: rl.NewVector2(100, 100), // Adjust the position as needed
	})

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

func downloadTexture(url string) (rl.Texture2D, error) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	img, _, _ := image.Decode(resp.Body)

	resizedImg := resize.Resize(uint(16), uint(16), img, resize.Lanczos3)

	rlImg := rl.NewImageFromImage(resizedImg)
	texture := rl.LoadTextureFromImage(rlImg)

	return texture, nil
}


// func main() {
//     rand.Seed(time.Now().UnixNano())

//     // Initialize game world
//     world := NewGameWorld()

//     // Initialize Raylib
//     rl.InitWindow(screenWidth, screenHeight, "Pokémon Game")
// 	rl.SetExitKey(0)
//     defer rl.CloseWindow()
//     rl.SetTargetFPS(60)

//     // Add a player
//     player := &Player{
//         Name:     "Ash",
//         Position: rl.Vector2{X: screenWidth / 2, Y: screenHeight / 2},
//     }
//     world.AddPlayer(player)

//     // Spawn initial Pokémons
//     for i := 0; i < 10; i++ {
//         world.SpawnPokemon()
//     }

//     // Main game loop
//     for !rl.WindowShouldClose() {
//         // Update game world
//         if rl.IsKeyDown(rl.KeyRight) {
//             world.MovePlayer(player, 5, 0)
//         }
//         if rl.IsKeyDown(rl.KeyLeft) {
//             world.MovePlayer(player, -5, 0)
//         }
//         if rl.IsKeyDown(rl.KeyUp) {
//             world.MovePlayer(player, 0, -5)
//         }
//         if rl.IsKeyDown(rl.KeyDown) {
//             world.MovePlayer(player, 0, 5)
//         }

//         // Draw
//         rl.BeginDrawing()
//         rl.ClearBackground(bkgColor)

//         // Draw player
//         rl.DrawCircleV(player.Position, 20, rl.Blue)

//         // Draw Pokémon
//         for _, pokemon := range world.Pokemons {
//             rl.DrawCircleV(pokemon.Position, 10, rl.Green)
//         }

//         rl.EndDrawing()
//     }
// }
