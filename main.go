package main

import (
	"fmt"
	tcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
	"github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg/request"
	"html/template"
	"log"
	"net/http"
)

type PokemonData struct {
	Cards []Card `json:"data"`
}

type Card struct {
	ID                     string       `json:"id"`
	Name                   string       `json:"name"`
	Supertype              string       `json:"supertype"`
	Subtypes               []string     `json:"subtypes"`
	Hp                     string       `json:"hp"`
	Types                  []string     `json:"types"`
	EvolvesFrom            string       `json:"evolvesFrom"`
	EvolvesTo              []string     `json:"evolvesTo"`
	Abilities              []Ability    `json:"abilities"`
	Attacks                []Attack     `json:"attacks"`
	Weaknesses             []Weakness   `json:"weaknesses"`
	Set                    []Set        `json:"set"`
	Resistances            []Resistance `json:"resistances"`
	RetreatCost            []string     `json:"retreatCost"`
	ConvertedRetreatCost   int          `json:"convertedRetreatCost"`
	Number                 string       `json:"number"`
	Artist                 string       `json:"artist"`
	Rarity                 string       `json:"rarity"`
	NationalPokedexNumbers []int        `json:"nationalPokedexNumbers"`
	Images                 struct {
		Small string `json:"small"`
		Large string `json:"large"`
	} `json:"images"`
}

type Ability struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Type string `json:"type"`
}

type Attack struct {
	Name                string   `json:"name"`
	Cost                []string `json:"cost"`
	ConvertedEnergyCost int      `json:"convertedEnergyCost"`
	Damage              string   `json:"damage"`
	Text                string   `json:"text"`
}

type Weakness struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Resistance struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Set struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Series       string `json:"series"`
	PrintedTotal int    `json:"printedTotal"`
	Total        int    `json:"total"`
	PtcgoCode    string `json:"ptcgoCode"`
	ReleaseDate  string `json:"releaseDate"`
	UpdatedAt    string `json:"updatedAt"`
	Images       struct {
		Small string `json:"small"`
		Large string `json:"large"`
	} `json:"images"`
}

func main() {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	c := tcg.NewClient("82a5348e-9d27-4af9-9d49-3886047cae6d")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		name := r.URL.Query().Get("name")
		setName := r.URL.Query().Get("set-name")

		var query string
		if name != "" && setName != "" {
			query = fmt.Sprintf("name:\"%s\" set.name:\"%s\"", name, setName)
		} else if name != "" {
			query = fmt.Sprintf("name:\"%s\"", name)
		} else if setName != "" {
			query = fmt.Sprintf("set.name:\"%s\"", setName)
		}

		cards, err := c.GetCards(
			request.Query(query),
			request.PageSize(7),
		)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching data from API: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, cards)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		fmt.Println(query)
		for _, card := range cards {
			fmt.Println(card.Name)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err.Error())
	}

}
