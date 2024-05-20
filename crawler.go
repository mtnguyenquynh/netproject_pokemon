package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Pokemon represents the structure of a Pokemon entity.
type Pokemon struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

// FetchPokemonData fetches Pokémon index and name from a given URL.
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
		index := s.Find("td").First().Text()
		name := s.Find("a[title*='(Pokémon)']").Text()

		// Clean up the extracted data
		index = strings.TrimSpace(index)
		name = strings.TrimSpace(name)

		// Print the extracted data for debugging purposes
		fmt.Printf("Extracted: Index = %s, Name = %s\n", index, name)

		if index != "" && name != "" {
			pokemon := Pokemon{
				Index: index,
				Name:  name,
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
