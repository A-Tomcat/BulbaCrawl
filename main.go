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
	Args       []string
	mu         *sync.Mutex
	wg         *sync.WaitGroup
	result     *Result
}

func main() {
	base, err := url.Parse("https://bulbapedia.bulbagarden.net")
	if err != nil {
		log.Fatal(err)
	}
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("No Searchname given.\n")
	}
	searchname := cases.Title(language.English).String(args[0])
	cfg := Config{
		BaseURL:    base,
		SearchName: searchname,
		Args:       args,
		mu:         &sync.Mutex{},
		wg:         &sync.WaitGroup{},
	}
	poke_Link, err := cfg.setSearchName("Pokémon")
	if err != nil {
		fmt.Println("Searchname not found.")
		log.Fatal(err)
	}
	if poke_Link == "Direct_Link" {
		if er := cfg.SpecificCard(); er != nil {
			log.Fatal(er)
		}
		return
	}
	if err := cfg.Pokemon(poke_Link); err != nil {
		log.Fatal(err)
	}
	if err := cfg.Cards(); err != nil {
		log.Fatal(err)
	}

}
