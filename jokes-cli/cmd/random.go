/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "get a single dad joke",
	Long:  `get a random single dad joke and display it in the terminal`,
	Run: func(cmd *cobra.Command, args []string) {

		jokeTerm, _ := cmd.Flags().GetString("term")

		if jokeTerm != "" {
			jokes := getRandomJokeWithTerm(jokeTerm)
			randomizeJokeList(len(jokes), jokes)
			return
		}

		getRandomJoke()
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.PersistentFlags().String("term", "", "A search term for a dad joke.")
	randomCmd.PersistentFlags().BoolP("output", "o", false, "indicates if output will be generates")
}

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func getRandomJoke() {
	url := "https://icanhazdadjoke.com/"

	responseBytes := getJokeData(url)

	joke := Joke{}
	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		log.Fatal("could not parse bytes", err)
	}

	fmt.Println(joke.Joke)
}

func getJokeData(url string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		log.Fatal("could not make request", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "jokes CLI")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal("could not get response", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("could not read response body", err)
	}

	return responseBytes
}

type SearchResult struct {
	Results    json.RawMessage `json:"results"`
	SearchTerm string          `json:"search_term"`
	Status     int             `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

func getRandomJokeWithTerm(jokeTerm string) []Joke {
	_, jokes := getJokeDataWithTerm(jokeTerm)
	return jokes
}

func getJokeDataWithTerm(jokeTerm string) (int, []Joke) {
	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeTerm)
	resposeBytes := getJokeData(url)
	jokeList := SearchResult{}

	if err := json.Unmarshal(resposeBytes, &jokeList); err != nil {
		log.Fatal("could not parse bytes")
	}

	jokes := []Joke{}

	if err := json.Unmarshal(jokeList.Results, &jokes); err != nil {
		log.Fatal("could not parse bytes")
	}

	return jokeList.TotalJokes, jokes

}

func randomizeJokeList(length int, jokeList []Joke) {
	rand.Seed(time.Now().Unix())
	min := 0
	max := length - 1

	if length <= 0 {
		log.Fatal("no joke found")
	}

	randomIndex := min + rand.Intn(max-min)
	fmt.Println(jokeList[randomIndex].Joke)

}
