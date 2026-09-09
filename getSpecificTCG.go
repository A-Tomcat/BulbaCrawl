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
func formatSpecificCard(card Card) {
	fmt.Printf("%s - %s - %s HP\n", card.Name, card.Type, card.HP)
	fmt.Println(card.Stage)
	fmt.Printf("Is weak to: %s, resistant to: %s, and has a Retreat Cost of: %d\n", card.Weakness, card.Resistance, card.RetreatCost)
	for _, a := range card.Ability {
		fmt.Printf("-%s:\n", a.Name)
		fmt.Printf("  %s\n", a.Effect)
	}
	fmt.Println(card.Attacks)
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
	card.getCardAbilityAttacks(doc)
	return card
}

func (cfg Config) SpecificCard() error {
	link, err := cfg.LinkGiven()
	if err != nil {
		return err
	}
	TCG_html, err := getHTML(link)
	if err != nil {
		return err
	}
	TCG_doc, err := HtmlToDoc(TCG_html)
	if err != nil {
		return err
	}
	card := cfg.getSpecificCardContent(TCG_doc)
	formatSpecificCard(card)
	return nil
}

func (card *Card) getCardAbilityAttacks(doc *goquery.Document) {
	sel := doc.Find(`[id="Card_text"]`).Parent().NextUntil(`h2`).Not("h3").First().Children()
	sel.Find("[title]").Each(func(i int, s *goquery.Selection) {
		attr, _ := s.Attr("title")
		card.swapAbilityAttack(sel, attr)
	})
}

/*
func (card *Card) getAttack(sel *goquery.Selection, attr string) {
	attacks := card.Attacks
	attack := CardAttacks{
		Cost: make(map[string]int),
	}

	attack.Name = "sel.Eq(1).Text()"

	attack.Damage = "sel.Eq(2).Text()"
	attack.Effect = "sel.Siblings().Text()"
	println(sel.Length())
	println(sel.Eq(1).Text())

	//Get Cost:
	cost := make(map[string]int)
	id := -1
	for i, at := range attacks {
		if attack.Name == at.Name {
			id = i
			break
		}
	}
	if id < 0 {
		attack.Cost[attr] = 1
	} else if _, ok := attacks[attack.Name]; !ok {

	}

	if _, ok := attack.Cost[attr]; !ok {
		attack.Cost[attr] = 1
	}
	sel.Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("title")
		if !checkElement(name) {
			return
		}
		if name == " " {
			name = "Free"
		}
		if _, ok := cost[name]; !ok {
			cost[name] = 1
		} else {
			cost[name]++
		}
	})
	attack.Cost = cost
	ok := false
	for _, at := range card.Attacks {
		if attack.Name == at.Name {
			ok = true
			break
		}
	}
	if !ok {
		card.Attacks = append(card.Attacks, attack)
	}
}*/

// Below Works as Intended!

func (card *Card) swapAbilityAttack(sel *goquery.Selection, attr string) {
	if attr == "Ability" || attr == "Poké-BODY" || attr == "Poké-POWER" || strings.Contains(sel.Text(), "Pokémon Power") {
		card.getAbility(sel)
	} else if checkElement(attr) {
		return
		//card.getAttack(sel, attr)
	}
}

func (card *Card) getAbility(sel *goquery.Selection) {
	ability := CardAbility{}
	//fmt.Println(sel.Text())
	s := sel.Parent().Parent().Parent().Parent().Children()
	name, _ := s.Find(`span[class="explain"]`).Attr("title")
	notnew := true
	for _, x := range card.Ability {
		if name == x.Name {
			notnew = false
			break
		}
	}
	if !notnew {
		return
	}
	ability.Name = name
	e := s.Find(`span[class="explain"]`).Parent().Parent().Parent().Next().First()
	ability.Effect = e.Text()
	card.Ability = append(card.Ability, ability)
}

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
			temp = strings.TrimSpace(temp)
			resistance = fmt.Sprintf("%s %s", rs, temp)
		}
	})
	return
}

func getEvoStage(sel *goquery.Selection) string {
	raw := sel.Children().First().Children().Eq(1).Text()
	text := strings.ReplaceAll(raw, "\n\n", "")
	text = strings.ReplaceAll(text, "PokémonEvolves", "Pokémon, Evolves")
	text = strings.Trim(text, "\n")
	return text
}

func getHP(sel *goquery.Selection) (hp string) {
	hp = sel.Find(`[title="HP (TCG)"]`).Parent().Siblings().Text()
	hp = strings.ReplaceAll(hp, "\n", "")
	return hp
}
