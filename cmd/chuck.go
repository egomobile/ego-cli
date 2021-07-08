// This file is part of the ego-cli distribution.
// Copyright (c) Next.e.GO Mobile SE, Aachen, Germany (https://e-go-mobile.com/)
//
// ego-cli is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, version 3.
//
// ego-cli is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/thatisuday/commando"
)

type ChuckNorrisJoke struct {
	Type  string               `json:"type"`
	Value ChuckNorrisJokeValue `json:"value"`
}

type ChuckNorrisJokeValue struct {
	Id         int64    `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}

func chuck_getJoke(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	var jokeOrError = ""

	explicit, err := flags["explicit"].GetBool()
	withExplicit := err == nil && explicit

	var url = "https://api.icndb.com/jokes/random?escape=javascript"

	if !withExplicit {
		url += "&exclude=[explicit]"
	}

	resp, err := http.Get(url)

	if err != nil {
		jokeOrError = fmt.Sprintf("ERROR: %v", err)
	} else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			jokeOrError = fmt.Sprintf("ERROR: %v", err)
		} else {
			var joke ChuckNorrisJoke

			err := json.Unmarshal(body, &joke)
			if err != nil {
				jokeOrError = fmt.Sprintf("ERROR: %v", err)
			} else {
				jokeOrError = joke.Value.Joke
			}
		}
	}

	fmt.Printf("%v", jokeOrError)
	fmt.Println()
}

func Setup_chuck_Command() {
	commando.
		Register("chuck").
		AddFlag("explicit,e", "include jokes with explicit vocabulary", commando.Bool, nil).
		SetShortDescription("random Chuck Norris joke").
		SetDescription("Tries to output a random Chuck Norris joke by using icndb.com service").
		SetAction(chuck_getJoke)
}
