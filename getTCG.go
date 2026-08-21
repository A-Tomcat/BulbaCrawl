package main

import (
	"fmt"
	"strings"

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
		HP           string        `json:"hp"`
		Ability      []CardAbility `json:"ability,omitempty"`
		Attacks      []CardAttacks `json:"attacks"`
		Resistance   string        `json:"resistance"`
		Weakness     string        `json:"weakness"`
		RetreatCost  int           `json:"retreat_cost"`
		PokedexEntry string        `json:"pokedex_entry"`
	}

	type CardAttacks struct {
		Name   string `json:"name"`
		Damage int    `json:"damage"`
		Effect string `json:"effect,omitempty"`
		Cost   string `json:"cost"`
	}

	type CardAbility struct {
		Name   string `json:"name"`
		Effect string `json:"effect"`
	}
*/

// Here starts the functions needed to recieve all TCG Cards of the specified Pokemon
func (cfg *Config) getCards(doc *goquery.Document) {
	//name := getTCGName(doc)
	names, links := getCardLinks(doc)
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
} //This Prints a formatted version of -<name> <link> for of the Pokémons cards.

func getCardLinks(doc *goquery.Document) (names, links []string) {
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

// Here ends the functions needed to recieve all TCG Cards of the specified Pokemon
func formatSpecificCard(card Card) {
	fmt.Printf("%s - %s - %s HP\n", card.Name, card.Type, card.HP)
	fmt.Println(card.Stage)
	fmt.Printf("Is weak to: %s, resistant to: %s, and has a Retreat Cost of: %d\n", card.Weakness, card.Resistance, card.RetreatCost)
	fmt.Println(card.PokedexEntry)
}
func (cfg *Config) getSpecificCardContent(doc *goquery.Document) (card Card) {
	card.Name = doc.Find("h1").Text()
	card.Stage = doc.Find(`[href="/wiki/Type_(TCG)"]`).Parent().Text()
	stats_sel := getStatsTable(doc)
	card.Type = getType(stats_sel)
	card.Weakness, card.Resistance, card.RetreatCost = getWRR(stats_sel)
	card.Stage = getEvoStage(stats_sel)
	card.HP = getHP(stats_sel)
	card.PokedexEntry = getTCGDexData(doc)
	card.Attacks, card.Ability = getCardAbilityAttacks(doc)
	return card
}

func getCardAbilityAttacks(doc *goquery.Document) (attacks []CardAttacks, abilities []CardAbility) {
	sel := doc.Find(`[id="Card_text"]`).Parent().NextUntil(`h2`).Not("h3").First()
	sel.Each(func(i int, s *goquery.Selection) {
		attr, ok := s.Attr("title")
		if ok == true {
			if attr == "Ability" || attr == "Poké-BODY" || attr == "Poké-POWER" {
			}
			if s.Text() == "Pokémon Power" {
				new_sel := s.Parent().Parent().Parent().Parent()
				abilities = append(abilities, getAbility(new_sel))
			}
			getAttack(sel, attr, attacks)
		}
	})
	return
}

func getAttack(sel *goquery.Selection, attr string, attacks []CardAttacks) {
	attack := CardAttacks{}
	if checkElement(attr) {
		sel = sel.Parent().Parent().Parent()
		cost := make(map[string]int)
		sel.First().Find("[title]").Each(func(i int, s *goquery.Selection) {
			name, _ := s.Attr("title")
			if !checkElement(name) {
				return
			}
			if _, ok := cost[name]; !ok {
				cost[name] = 1
			} else {
				cost[name]++
			}
		})
		attack.Cost = cost
		attack.Damage = sel.Eq(2).Text()
		attack.Name = sel.Eq(1).Text()
		attack.Effect = sel.Siblings().Text()
		ok := false
		for _, at := range attacks {
			if attack.Name == at.Name {
				ok = true
				break
			}
		}
		if !ok {
			attacks = append(attacks, attack)
		}
	}
}

func getAbility(sel *goquery.Selection) CardAbility {
	ability := CardAbility{}
	ability.Name = sel.First().Eq(0).Eq(1).Text()
	ability.Effect = sel.First().Eq(1).Text()
	return ability
}

// Below Works as Intended!

func getTCGDexData(doc *goquery.Document) string {
	sel := doc.Find(`[id="Pokédex_data"]`).Parent().Next()
	if sel.Text() == "" {
		return "No PokéDex Entry on this Pokémon Card."
	}
	raw := sel.Text()
	de_spaced := strings.TrimSpace(raw)
	de_ln := strings.ReplaceAll(de_spaced, "\n\n", "\n")
	de_ln = strings.Replace(de_ln, "\n", "", 1)
	clean := strings.Replace(de_ln, "Pokédex entry", "Pokédex entry:\n", 1)
	return clean
}

func checkElement(attr string) bool {
	if attr == " " || attr == "Grass" || attr == "Fire" || attr == "Water" || attr == "Lightning" || attr == "Fighting" || attr == "Psychic" || attr == "Colorless" || attr == "Darkness" || attr == "Metal" || attr == "Dragon" || attr == "Fairy" {
		return true
	}
	return false
}
func getStatsTable(doc *goquery.Document) (sel *goquery.Selection) {
	sel = doc.Find(`[title="Type (TCG)"]`).Parent().Parent().Parent()
	return
}

func getType(sel *goquery.Selection) (element string) {
	s := sel.Find(`[title="Type (TCG)"]`).Parent().Siblings().Find("[title]")
	element, _ = s.Attr("title")
	return
}

func getWRR(sel *goquery.Selection) (weakness string, resistance string, retreat_cost int) {
	sel.Find(`small`).Each(func(i int, s *goquery.Selection) {
		if s.Text() == "retreat cost" {
			retreat_cost = s.Parent().Find("span[typeof]").Length()
		}
		if s.Text() == "weakness" {
			wk, _ := s.Siblings().Find(`[title]`).Attr("title")
			temp := strings.ReplaceAll(s.Parent().Text(), "weakness", "")
			temp = strings.ReplaceAll(temp, "\n", "")
			weakness = fmt.Sprintf("%s %s", wk, temp)
		}
		if s.Text() == "resistance" {
			rs, _ := s.Siblings().Find(`[title]`).Attr("title")
			temp := strings.ReplaceAll(s.Parent().Text(), "resistance", "")
			temp = strings.ReplaceAll(temp, "\n", "")
			resistance = fmt.Sprintf("%s %s", rs, temp)
		}
	})
	return
}

func getEvoStage(sel *goquery.Selection) string {
	raw := sel.Children().First().Children().Eq(1).Text()
	text := strings.ReplaceAll(raw, "\n\n", "")
	text = strings.ReplaceAll(text, "PokémonEvolves", "Pokémon, Evolves")
	text = strings.TrimLeft(text, "\n")
	return text
}

func getHP(sel *goquery.Selection) (hp string) {
	hp = sel.Find(`[title="HP (TCG)"]`).Parent().Siblings().Text()
	hp = strings.ReplaceAll(hp, "\n", "")
	return hp
}
