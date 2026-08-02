package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

/*
// All the Pokemon structs
type Pokemon_result struct {
	Name           string      `json:"name"`
	PokeDex_Number int         `json:"pokedex_number"`
	BaseStats      []BaseStats `json:"base_stats"`
	Prev_Evolution string      `json:"prev_evolution,omitempty"`
	Next_Evolution []string    `json:"next_evoltion,omitempty"`
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

func getBaseStats(doc *goquery.Document) ([]BaseStats, error) {
	docs := findStatTables(doc)
	ReturnStats := []BaseStats{}
	for _, doc := range docs {
		doc.Each(func(i int, s *goquery.Selection) {
			ReturnStats = append(ReturnStats, getBaseStat(s))
		})
	}
	if ReturnStats == nil {
		return []BaseStats{}, fmt.Errorf("No Stats found.")
	}
	return ReturnStats, nil
}

func findStatTables(doc *goquery.Document) []*goquery.Selection {
	tables := []*goquery.Selection{}
	//doc_part := doc.Find(`h4[id="Base_stats"]`).NextAll()
	doc.Find("table").Find("tbody").Each(func(i int, s *goquery.Selection) {
		if href, _ := s.Attr("href"); href == "/wiki/Stat" {
			tables = append(tables, s)
		}
	})
	return tables
}

func getBaseStat(s *goquery.Selection) BaseStats {
	stats := BaseStats{}
	s.Each(func(i int, s *goquery.Selection) {
		switch i {
		case 2:
			stats.Hp = s.First().Text()
		case 3:
			stats.Atk = s.First().Text()
		case 4:
			stats.Def = s.First().Text()
		case 5:
			stats.Sp_Att = s.First().Text()
		case 6:
			stats.Sp_Def = s.First().Text()
		case 7:
			stats.Speed = s.First().Text()
		case 8:
			stats.Total = s.First().Text()
		default:
			return
		}
	})
	return stats
} /*getBaseStat() returns the BaseStats struct filled with the Base Stats contained in the given Selection*/

func findEvo(doc *goquery.Document) *goquery.Selection {
	doc_part := doc.Has(`h3[id="Evolution_data"]`).NextAll()
	evo := doc_part.Find("table").First()
	return evo
} /*findEvo() returns the Selector containing only the Table with the Evolution Data*/

func getEvoLine(doc *goquery.Selection) []string {
	evolutions_Line := []string{}
	tables := []*goquery.Selection{}
	doc.Each(func(_ int, s *goquery.Selection) {
		if s.Length() == 3 {
			tables = append(tables, s.Eq(2))
		}
	})
	if len(tables) == 0 {
		return evolutions_Line
	}
	for _, table := range tables {
		Poke_name := table.First().Text()
		evolutions_Line = append(evolutions_Line, Poke_name)
	}
	return evolutions_Line
} /*getEvoLine() returns the Names of all Members of the Evolution line.*/
