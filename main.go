package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Config struct {
	BaseURL    *url.URL
	SearchName string
	Category   string
	mu         *sync.Mutex
	wg         *sync.WaitGroup
	result     *Result
}

type Result struct {
	Pokemon struct {
		Name      string `json:"name"`
		BaseStats struct {
			Hp     int `json:"hp"`
			Atk    int `json:"atk"`
			Def    int `json:"def"`
			Sp_Att int `json:"sp_att"`
			Sp_Def int `json:"sp_def"`
			Speed  int `json:"spd"`
			Total  int `json:"total"`
		} `json:"base_stats"`
		Prev_Evolution string   `json:"prev_evolution,omitempty"`
		Next_Evolution []string `json:"next_evoltion,omitempty"`
	} `json:"pokemon,omitempty"`
	TCG struct {
		Name          string `json:"name"`
		NumberOfCards int    `json:"cards_released"`
	} `json:"tcg,omitempty"`
	Move struct {
		Name       string `json:"name"`
		PP         int    `json:"pp"`
		Power      int    `json:"power"`
		Accuracy   int    `json:"accuracy"`
		Priority   int    `json:"priority"`
		LearnedLVL int    `json:"learned_lvl"`
		LearnedTM  int    `json:"learner_tm"`
		LearnedEgg int    `json:"learned_egg"`
	} `json:"move,omitempty"`
}

func main() {
	fmt.Println("Hello Moron.")
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal(`Usage: ./pokecrawl <tcg/pokemon/move> "<name>"`)
	}
	baseURL := os.Getenv("BASEURL")
	base_url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	category := strings.ToLower(args[0])
	switch category {
	case "tcg":
		category = "TCG"
	case "move":
		category = "move"
	case "pokemon":
		category = "pokémon"
	default:
		log.Fatal("Category not unique, Options: <move>, <tcg>, <pokemon>")
	}
	searchname := cases.Title(language.English).String(args[1])
	var wg sync.WaitGroup
	var mutex sync.Mutex
	result := Result{}
	cfg := Config{
		BaseURL:    base_url,
		SearchName: searchname,
		Category:   category,
		mu:         &mutex,
		wg:         &wg,
		result:     &result,
	}
	fmt.Println(cfg)
}
