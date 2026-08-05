package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

/*
// All the TCG structs

	type TCG_Result struct {
		Name          string `json:"name"`
		NumberOfCards int    `json:"cards_released"`
		Cards         []Card `json:"cards"`
	}

	type Card struct {
		Name         string        `json:"name"`
		Type         string        `json:"type"`
		Stage        string        `json:"stage"`
		EvolvesFrom  string        `json:"evolves_from,omitempty"`
		Ability      string        `json:"ability,omitempty"`
		Attacks      []CardAttacks `json:"attacks"`
		Resistance   string        `json:"resistance"`
		Weakness     string        `json:"weakness"`
		RetreatCost  int           `json:"retreat_cost"`
		PokedexEntry string        `json:"pokedex_entry"`
	}

	type CardAttacks struct {
		Name   string `json:"name"`
		Damage int    `json:"damage"`
		Eff
*/
func (cfg *Config) getCards(doc *goquery.Document) {
	//name := getTCGName(doc)
	names, links := cfg.getCardLinks(doc)
	cfg.formatAllCards(names, links)

}
func getTCGName(doc *goquery.Document) string {
	name := doc.Find(`h1[id*="firstHeading"]`).First().First().Text()
	return name
}
func (cfg *Config) formatAllCards(names, links []string) {
	fmt.Printf("%s has %d different Cards in the Pokémon TCG:\n", cfg.SearchName, len(names))
	baseLink := "https://bulbapedia.bulbagarden.net"
	for i, name := range names {
		link := baseLink + links[i]
		fmt.Printf(" -%s, %s \n", name, link)
	}
}

func (cfg *Config) getCardLinks(doc *goquery.Document) (names, links []string) {
	selector := `tbody tr[style^="background"]`
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		sel := s.Find(`a[href]`)
		name, _ := sel.Attr("title")
		link, _ := sel.Attr("href")
		if name != "" {
			links = append(links, link)
			names = append(names, name)
		}
	})
	return names, links
}
