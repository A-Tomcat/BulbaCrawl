package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
// All the Pokemon structs
type Pokemon_result struct {
	Name           string      `json:"name"`
	PokeDex_Number int         `json:"pokedex_number"`
	BaseStats      []BaseStats `json:"base_stats"`
	Evolution_Line []string `json:"evolution_line"`
}
type BaseStats struct {
	Hp     int `json:"hp"`
	Atk    int `json:"atk"`
	Def    int `json:"def"`
	Sp_Att int `json:"sp_att"`
	Sp_Def int `json:"sp_def"`
	Speed  int `json:"spd"`
	Total  int `json:"total"`
}*/

func getPokemon(doc *goquery.Document) (Pokemon_result, error) {
	stats, err := getBaseStats(doc)
	if err != nil {
		return Pokemon_result{}, err
	}
	evos := getEvoLine(doc)
	raw_name := getName(doc)
	dex := getDexNumber(doc)
	name := strings.Replace(raw_name, "(Pokémon)", "", 1)
	name = strings.TrimSpace(name)
	result := Pokemon_result{
		Name:           name,
		BaseStats:      stats,
		Evolution_Line: evos,
		PokeDex_Number: dex,
	}
	return result, nil
}
func (cfg Config) Pokemon(link string) error {
	Poke_html, err := getHTML(link)
	if err != nil {
		return err
	}
	Poke_doc, err := HtmlToDoc(Poke_html)
	if err != nil {
		return err
	}
	pokemon, err := getPokemon(Poke_doc)
	formatPokemon(pokemon)
	return nil
}

func formatPokemon(p Pokemon_result) {
	fmt.Printf("%s %s \n", p.PokeDex_Number, p.Name)
	if len(p.Evolution_Line) > 1 {
		fmt.Printf("Evolution Line: %s\n", p.Evolution_Line)
	} else {
		fmt.Println("Single Stage Pokémon")
	}
	fmt.Printf("Base Stats:\n -HP: %s\n -Attack: %s\n -Defense: %s\n -Special Attack: %s\n -Special Defense: %s\n -Speed: %s\n  -Total: %s\n\n", p.BaseStats.Hp, p.BaseStats.Atk, p.BaseStats.Def, p.BaseStats.Sp_Att, p.BaseStats.Sp_Def, p.BaseStats.Speed, p.BaseStats.Total)
}

func getName(doc *goquery.Document) string {
	name := doc.Find(`h1[id*="firstHeading"]`).First().First().Text()
	return name
}

func getDexNumber(doc *goquery.Document) string {
	dex_string := doc.Find(`a[title*="National Pokédex"]`).Text()
	var numbers [10]rune
	for i := 0; i < 10; i++ {
		numbers[i] = rune(i)
	}
	nums := regexp.MustCompile("[0-9]+")
	dex := nums.FindAllString(dex_string, -1)
	dexn := ""
	for _, c := range dex {
		dexn = dexn + c
	}
	return dexn
}
func getBaseStats(doc *goquery.Document) (BaseStats, error) {
	docs := findStatTables(doc)
	stats := BaseStatFinder(docs)
	return stats, nil
}

func findStatTables(doc *goquery.Document) *goquery.Selection {
	sel := doc.Find(`a[href="/wiki/Stat"]`).First().Parent().Parent().Parent().Parent()
	return sel
}

func BaseStatFinder(sel *goquery.Selection) BaseStats {
	stats := BaseStats{}
	stats.Hp = sel.Find(`[href="/wiki/HP"]`).Parent().Next().Text()
	stats.Atk = sel.Find(`[href="/wiki/Stat#Attack"]`).Parent().Next().Text()
	stats.Def = sel.Find(`[href="/wiki/Stat#Defense"]`).Parent().Next().Text()
	stats.Sp_Att = sel.Find(`[href="/wiki/Stat#Special_Attack"]`).Parent().Next().Text()
	stats.Sp_Def = sel.Find(`[href="/wiki/Stat#Special_Defense"]`).Parent().Next().Text()
	stats.Speed = sel.Find(`[href="/wiki/Stat#Speed"]`).Parent().Next().Text()
	sel.Find("th").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "Total") {
			stats.Total = s.Last().Text()[6:9]
		}
		return true
	})
	return stats
}
func findEvo(doc *goquery.Document) *goquery.Selection {
	doc_part := doc.Find(`[id="Evolution_data"]`)
	nextdoc := doc_part.Parent().NextUntil(`h3`).Find("table")
	return nextdoc
} /*findEvo() returns the Selector containing only the Table with the Evolution Data*/

func getEvoLine(doc *goquery.Document) []string {
	sel := findEvo(doc)
	evolutions_Line := []string{}
	sel.Find("tr").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "evolved") || strings.Contains(s.Text(), "Evolution") {
			name_sel := s.Next().Find("a").First()
			name := strings.TrimSpace(name_sel.Text())
			evolutions_Line = append(evolutions_Line, name)
		}
	})
	if len(evolutions_Line[0]) == 0 {
		evolutions := []string{}
		evolutions = append(evolutions, evolutions_Line[1:]...)
		return evolutions
	}
	return evolutions_Line
} /*getEvoLine() returns the Names of all Members of the Evolution line.*/
