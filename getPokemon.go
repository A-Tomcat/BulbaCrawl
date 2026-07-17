package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
// All the Pokemon structs
type Pokemon_result struct {
	Name           string    `json:"name"`
	BaseStats      BaseStats `json:"base_stats"`
	Prev_Evolution string    `json:"prev_evolution,omitempty"`
	Next_Evolution []string  `json:"next_evoltion,omitempty"`
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
	ReturnStats := []BaseStats{}
	println(doc)
	strings.ToLower("B")
	doc.Find("table").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.HasPrefix(s.Text(), "Stat") {
			stats := getBaseStat(s)
			ReturnStats = append(ReturnStats, stats)
		}
		if len(ReturnStats) <= 3 {
			return true
		} else {
			return false
		}
	})
	return ReturnStats, nil
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
}
