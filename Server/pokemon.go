package main
import
("strings")

func findPokemon(pokedex []Pokemon, name string) (Pokemon, bool) {
	for _, p := range pokedex {
		if strings.EqualFold(p.Name, name) {
			return p, true
		}
	}
	return Pokemon{}, false
}