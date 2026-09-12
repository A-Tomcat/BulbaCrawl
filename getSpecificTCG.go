package main

import (
	"fmt"
	"sort"
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
	fmt.Printf("Is weak to: %s, resistant to: %s, and has a Retreat Cost of: %d\n", strings.TrimSpace(card.Weakness), strings.TrimSpace(card.Resistance), card.RetreatCost)
	//Abilities
	for _, a := range card.Ability {
		fmt.Printf("-%s:\n", a.Name)
		fmt.Printf("%s\n", a.Effect)
	}
	//Attacks
	for _, attack := range card.Attacks {
		fmt.Printf("-%s:\n", attack.Name)
		damage := attack.Damage
		if damage == "" {
			damage = "0"
		}
		fmt.Print("Requires following Energie Cards to activate: ")
		types := sortMap(attack.Cost)
		for _, element := range types {
			if element == "Free" {
				fmt.Println("None")
				break
			}
			fmt.Printf("%d %s", attack.Cost[element], element)
		}
		fmt.Printf("\nDoes %s Damage to the opposing active Pokémon\n", damage)
		if attack.Effect != "" {
			fmt.Printf("%s\n", attack.Effect)
		}
	}
	fmt.Println()
	fmt.Println(card.PokedexEntry)
}

func sortMap(m map[string]int) (keyList []string) {
	for key := range m {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	return
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
	sel := doc.Find(`[id="Card_text"]`).Parent().NextAllFiltered(`div.roundy`).First()
	sel.Children().Each(func(i int, s *goquery.Selection) {
		if s.Find(`a[title="Poké-POWER"], a[title="Poké-BODY"], a[title="Ability"]`).Length() > 0 || strings.Contains(s.Text(), "Pokémon Power") {
			card.getAbility(s)
			return
		}
		card.getAttack(s)
	})
}

func (card *Card) getAttack(sel *goquery.Selection) {
	//Getting base Children to work with
	children := sel.Children()
	//Header is a Child element that possesses the Cost, Name and Damage an Attack does
	header := children.First().Children()
	//Check if the Sel actually contains an attack, in this case by checking for the Translation
	//Because we check if it's an Ability before going into this Function it should be definitiv.
	if s := header.Find(`div[lang] span.explain`); s.Length() == 0 {
		return
	}
	//Get Attack Name
	name_sel := header.Eq(1)
	clone := name_sel.Clone()
	clone.Find(`div[lang]`).Remove()
	at_name := clone.Text()
	//Get Attack Damage, if it does any
	at_dmg := header.Eq(2).Text()
	//Get Attack Cost
	attack_Cost := make(map[string]int)
	cost_s := header.First().Find(`img.mw-file-element`)
	cost_s.Each(func(i int, s *goquery.Selection) {
		element, ok := s.Attr(`alt`)
		if !ok || element == "" {
			return
		}
		if element == "\u00a0" {
			element = "Free"
		}
		if _, ok := attack_Cost[element]; !ok {
			attack_Cost[element] = 1
		} else {
			attack_Cost[element]++
		}
	})
	//Taking secondary Child to get Effect, if the Attack does more then just Damage
	effect := children.Eq(1).Text()
	attack := CardAttacks{
		Name:   at_name,
		Effect: effect,
		Damage: at_dmg,
		Cost:   attack_Cost,
	}
	card.Attacks = append(card.Attacks, attack)
}

func (card *Card) getAbility(sel *goquery.Selection) {
	ability := CardAbility{}
	children := sel.Children()
	if children.Length() < 2 {
		return
	}
	header := children.First()
	name := header.Find(`div[lang] span.explain`).Parent().Parent()
	name_clone := name.Clone()
	name_clone.Find(`div[lang]`).Remove()
	name_text := name_clone.Text()
	ability.Name = name_text

	effect_s := children.Eq(1)
	ability.Effect = effect_s.Text()
	card.Ability = append(card.Ability, ability)
}

func getTCGDexData(doc *goquery.Document) string {
	sel := doc.Find(`[id="Pokédex_data"]`).Parent().Next()
	if sel.Text() == "" {
		return "No PokéDex Entry on this Pokémon Card."
	}
	clone := sel.Clone()
	clone.Find(`div[style="margin-top: 0.125rem"]`).Children().Eq(2).Remove()
	raw := clone.Text()
	de_spaced := strings.TrimSpace(raw)
	de_ln := strings.ReplaceAll(de_spaced, "\n\n", "\n")
	de_ln = strings.Replace(de_ln, "\n", "", 1)
	clean := strings.Replace(de_ln, "Pokédex entry", "Pokédex entry:\n", 1)
	clean = strings.ReplaceAll(clean, "\nNo.\n", "\nNo.: ")
	clean = strings.ReplaceAll(clean, "\nHeight\n", ", Height: ")
	clean = strings.ReplaceAll(clean, "\nWeight\n", ", Weight: ")
	return clean
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
