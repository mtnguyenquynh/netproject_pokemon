/* package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// Pokemon represents the structure of a Pokemon entity.
type Pokemon struct {
	Index       string `json:"index"`
	Name        string `json:"name"`
	Exp         int    `json:"exp"`
	HP          int    `json:"hp"`
	Attack      int    `json:"attack"`
	Defense     int    `json:"defense"`
	SpAttack    int    `json:"sp_attack"`
	SpDefense   int    `json:"sp_defense"`
	Speed       int    `json:"speed"`
	TotalEVs    int    `json:"total_evs"`
}

// FetchPokemonData fetches Pokémon data from a given URL.
func FetchPokemonData(url string) ([]Pokemon, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var pokemons []Pokemon
	doc.Find("table.roundy tbody tr").Each(func(i int, s *goquery.Selection) {
		index := strings.TrimSpace(s.Find("td").First().Text())
		name := strings.TrimSpace(s.Find("a[title*='(Pokémon)']").Text())
		expStr := strings.TrimSpace(s.Find("td").Eq(3).Text())
		hpStr := strings.TrimSpace(s.Find("td").Eq(4).Text())
		attackStr := strings.TrimSpace(s.Find("td").Eq(5).Text())
		defenseStr := strings.TrimSpace(s.Find("td").Eq(6).Text())
		spAttackStr := strings.TrimSpace(s.Find("td").Eq(7).Text())
		spDefenseStr := strings.TrimSpace(s.Find("td").Eq(8).Text())
		speedStr := strings.TrimSpace(s.Find("td").Eq(9).Text())
		totalEVsStr := strings.TrimSpace(s.Find("td").Eq(10).Text())

		// Convert strings to integers
		exp, _ := strconv.Atoi(expStr)
		hp, _ := strconv.Atoi(hpStr)
		attack, _ := strconv.Atoi(attackStr)
		defense, _ := strconv.Atoi(defenseStr)
		spAttack, _ := strconv.Atoi(spAttackStr)
		spDefense, _ := strconv.Atoi(spDefenseStr)
		speed, _ := strconv.Atoi(speedStr)
		totalEVs, _ := strconv.Atoi(totalEVsStr)

		// Print the extracted data for debugging purposes
		fmt.Printf("Extracted: Index = %s, Name = %s, Exp = %d, HP = %d, Attack = %d, Defense = %d, SpAttack = %d, SpDefense = %d, Speed = %d, TotalEVs = %d\n",
			index, name, exp, hp, attack, defense, spAttack, spDefense, speed, totalEVs)

		if index != "" && name != "" {
			pokemon := Pokemon{
				Index:    index,
				Name:     name,
				Exp:      exp,
				HP:       hp,
				Attack:   attack,
				Defense:  defense,
				SpAttack: spAttack,
				SpDefense: spDefense,
				Speed:    speed,
				TotalEVs: totalEVs,
			}
			pokemons = append(pokemons, pokemon)
		}
	})

	return pokemons, nil
}

func main() {
	url := "https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield_(Generation_IX)"
	pokemons, err := FetchPokemonData(url)
	if err != nil {
		log.Fatalf("Error fetching Pokémon data: %v", err)
	}

	file, err := os.Create("pokedex.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(pokemons); err != nil {
		log.Fatalf("Error encoding JSON to file: %v", err)
	}

	log.Println("Pokedex data has been written to pokedex.json")
}
 */

 /* package main

 import (
	 "encoding/json"
	 "fmt"
	 "log"
	 "net/http"
	 "os"
	 "strconv"
	 "strings"
 
	 "github.com/PuerkitoBio/goquery"
 )
 
 // Pokemon represents the structure of a Pokemon entity.
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
 }
 
 // FetchPokemonData fetches Pokémon data from the main list page and individual Pokémon pages.
 func FetchPokemonData(url string) ([]Pokemon, error) {
	 resp, err := http.Get(url)
	 if err != nil {
		 return nil, err
	 }
	 defer resp.Body.Close()
 
	 doc, err := goquery.NewDocumentFromReader(resp.Body)
	 if err != nil {
		 return nil, err
	 }
 
	 var pokemons []Pokemon
	 doc.Find("#monsters-list li").Each(func(i int, s *goquery.Selection) {
		 if i >= 5 { // Limit to the first 5 Pokémon
			 return
		 }
 
		 pokemonName := strings.TrimSpace(s.Find("span").Text())
		 pokemonID := strings.TrimPrefix(s.Find("button").AttrOr("class", ""), "monster-sprite sprite-")
		 pokemonURL := fmt.Sprintf("https://pokedex.org/#/pokemon/%s", pokemonID)
 
		 fmt.Printf("Fetching data for %s (%s)\n", pokemonName, pokemonURL) // Debugging
 
		 pokemon, err := fetchIndividualPokemonData(pokemonURL, pokemonName, pokemonID)
		 if err == nil {
			 pokemons = append(pokemons, pokemon)
			 fmt.Printf("Fetched: %+v\n", pokemon) // Debugging
		 } else {
			 fmt.Printf("Error fetching data for %s: %v\n", pokemonName, err) // Debugging
		 }
	 })
 
	 return pokemons, nil
 }
 
// fetchIndividualPokemonData fetches data for an individual Pokémon.
func fetchIndividualPokemonData(url, name, index string) (Pokemon, error) {
    resp, err := http.Get(url)
    if err != nil {
        return Pokemon{}, err
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return Pokemon{}, err
    }

    // Initialize a new Pokemon struct
    var pokemon Pokemon
    pokemon.Index = index
    pokemon.Name = name
    pokemon.Type = []string{} // Reset Type slice to avoid data carry over

    // Extracting types
    doc.Find(".detail-types .monster-type").Each(func(i int, s *goquery.Selection) {
        pokemon.Type = append(pokemon.Type, strings.TrimSpace(s.Text()))
    })

    // Extracting stats
    doc.Find(".detail-stats-row").Each(func(i int, s *goquery.Selection) {
        statName := strings.TrimSpace(s.Find("span").First().Text())
        statValueStr := strings.TrimSpace(s.Find(".stat-bar-fg").Text())
        statValue, _ := strconv.Atoi(statValueStr)
        switch statName {
        case "HP":
            pokemon.HP = statValue
        case "Attack":
            pokemon.Attack = statValue
        case "Defense":
            pokemon.Defense = statValue
        case "Sp Atk":
            pokemon.SpAttack = statValue
        case "Sp Def":
            pokemon.SpDefense = statValue
        case "Speed":
            pokemon.Speed = statValue
        }
    })

    // Extracting additional info
    pokemon.Description = strings.TrimSpace(doc.Find(".monster-description").Text())
    pokemon.Height = strings.TrimSpace(doc.Find(".monster-minutia span").Eq(0).Text())
    pokemon.Weight = strings.TrimSpace(doc.Find(".monster-minutia span").Eq(1).Text())

    // Sum of stats as Total EVs
    pokemon.TotalEVs = pokemon.HP + pokemon.Attack + pokemon.Defense + pokemon.SpAttack + pokemon.SpDefense + pokemon.Speed

    return pokemon, nil
}

 
 func main() {
	 url := "https://pokedex.org/"
	 pokemons, err := FetchPokemonData(url)
	 if err != nil {
		 log.Fatalf("Error fetching Pokémon data: %v", err)
	 }
 
	 file, err := os.Create("pokedex.json")
	 if err != nil {
		 log.Fatalf("Error creating file: %v", err)
	 }
	 defer file.Close()
 
	 encoder := json.NewEncoder(file)
	 if err := encoder.Encode(pokemons); err != nil {
		 log.Fatalf("Error encoding JSON to file: %v", err)
	 }
 
	 log.Println("Pokedex data has been written to pokedex.json")
 }
 
  */


  package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/PuerkitoBio/goquery"
)

// Pokemon represents the structure of a Pokemon entity.
type Pokemon struct {
	Index       string   `json:"index"`
	Name        string   `json:"name"`
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
}

// FetchPokemonData fetches Pokémon data from the main list page and individual Pokémon pages.
func FetchPokemonData(ctx context.Context) ([]Pokemon, error) {
	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://pokedex.org/"),
		chromedp.Sleep(5*time.Second), // Wait for the page to load completely
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load main page: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse main page HTML: %v", err)
	}

	var pokemons []Pokemon
	doc.Find("#monsters-list li").Each(func(i int, s *goquery.Selection) {
		if i >= 5 { // Limit to the first 5 Pokémon
			return
		}

		pokemonName := strings.TrimSpace(s.Find("span").Text())
		pokemonID := strings.TrimPrefix(s.Find("button").AttrOr("class", ""), "monster-sprite sprite-")
		pokemonURL := fmt.Sprintf("https://pokedex.org/#/pokemon/%s", pokemonID)

		fmt.Printf("Fetching data for %s (%s)\n", pokemonName, pokemonURL)

		pokemon, err := fetchIndividualPokemonData(ctx, pokemonURL, pokemonName, pokemonID)
		if err == nil {
			pokemons = append(pokemons, pokemon)
			fmt.Printf("Fetched: %+v\n", pokemon)
		} else {
			fmt.Printf("Error fetching data for %s: %v\n", pokemonName, err)
		}
	})

	return pokemons, nil
}

// fetchIndividualPokemonData fetches data for an individual Pokémon.
func fetchIndividualPokemonData(ctx context.Context, url, name, index string) (Pokemon, error) {
	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second), // Wait for the page to load completely
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to load Pokémon page: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to parse Pokémon page HTML: %v", err)
	}

	var pokemon Pokemon
	pokemon.Index = index
	pokemon.Name = name
	pokemon.Type = []string{} // Reset Type slice to avoid data carry over

	doc.Find(".detail-types .monster-type").Each(func(i int, s *goquery.Selection) {
		pokemon.Type = append(pokemon.Type, strings.TrimSpace(s.Text()))
	})

	doc.Find(".detail-stats-row").Each(func(i int, s *goquery.Selection) {
		statName := strings.TrimSpace(s.Find("span").First().Text())
		statValueStr := strings.TrimSpace(s.Find(".stat-bar-fg").Text())
		statValue, _ := strconv.Atoi(statValueStr)
		switch statName {
		case "HP":
			pokemon.HP = statValue
		case "Attack":
			pokemon.Attack = statValue
		case "Defense":
			pokemon.Defense = statValue
		case "Sp Atk":
			pokemon.SpAttack = statValue
		case "Sp Def":
			pokemon.SpDefense = statValue
		case "Speed":
			pokemon.Speed = statValue
		}
	})

	pokemon.Description = strings.TrimSpace(doc.Find(".monster-description").Text())
	pokemon.Height = strings.TrimSpace(doc.Find(".monster-minutia span").Eq(0).Text())
	pokemon.Weight = strings.TrimSpace(doc.Find(".monster-minutia span").Eq(1).Text())

	pokemon.TotalEVs = pokemon.HP + pokemon.Attack + pokemon.Defense + pokemon.SpAttack + pokemon.SpDefense + pokemon.Speed

	return pokemon, nil
}

func main() {
	// Create context
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Flag("disable-dev-shm-usage", true),
	}
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Run the data fetching
	pokemons, err := FetchPokemonData(ctx)
	if err != nil {
		log.Fatalf("Error fetching Pokémon data: %v", err)
	}

	file, err := os.Create("pokedex.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(pokemons); err != nil {
		log.Fatalf("Error encoding JSON to file: %v", err)
	}

	log.Println("Pokedex data has been written to pokedex.json")
}
