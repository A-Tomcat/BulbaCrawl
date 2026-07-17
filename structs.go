package main

type Result struct {
	Pokemon Pokemon_result `json:"pokemon,omitempty"`
	TCG     TCG_Result     `json:"tcg,omitempty"`
	Move    MoveResult     `json:"move,omitempty"`
}

// All the Pokemon structs
type Pokemon_result struct {
	Name           string      `json:"name"`
	BaseStats      []BaseStats `json:"base_stats"`
	Prev_Evolution string      `json:"prev_evolution,omitempty"`
	Next_Evolution []string    `json:"next_evoltion,omitempty"`
}
type BaseStats struct {
	Hp     string `json:"hp"`
	Atk    string `json:"atk"`
	Def    string `json:"def"`
	Sp_Att string `json:"sp_att"`
	Sp_Def string `json:"sp_def"`
	Speed  string `json:"spd"`
	Total  string `json:"total"`
}

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
	Effect string `json:"effect,omitempty"`
}

// All the Move structs
type MoveResult struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Type     string `json:"type"`
	PP       string `json:"pp"`
	Power    string `json:"power"`
	Accuracy string `json:"accuracy"`
	/* Maybe add those
	LearnByLVL int    `json:"learn_by_lvl"`
	LearnByTM  int    `json:"learn_by_tm"`
	LearnByEgg int    `json:"learn_by_egg"`
	*/
}
