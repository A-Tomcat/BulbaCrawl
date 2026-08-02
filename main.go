package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
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

func main() {
	base, err := url.Parse("https://bulbapedia.bulbagarden.net")
	if err != nil {
		log.Fatal(err)
	}
	//split
	/*
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
	*/
	/*cfg := Config{
		BaseURL:    base_url,
		SearchName: searchname,
		Category:   category,
		mu:         &mutex,
		wg:         &wg,
		result:     &result,
	}
		split end
	*/
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("No Searchname given.\n")
	}
	searchname := cases.Title(language.English).String(args[0])
	cfg := Config{
		BaseURL:    base,
		SearchName: searchname,
		Category:   "pokemon",
		mu:         &sync.Mutex{},
		wg:         &sync.WaitGroup{},
	}
	html, err := cfg.getHTML()
	if err != nil {
		log.Fatal(err)
	}
	doc, err := HtmlToDoc(html)
	if err != nil {
		log.Fatal(err)
	}
	pokemon, err := getPokemon(doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pokemon)
	/*moveResult, err := getMoveFromDoc(doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(moveResult)*/

}
