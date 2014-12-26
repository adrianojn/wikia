// Copyright (C) 2014 Adriano Soares
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type WikiaResult map[string]struct {
	Title     string
	Revisions []struct {
		Text string `json:"*"`
	}
}

var resultJSON = make(WikiaResult)

var missedCards []string

func wikia() {
	cards := getListOfPages()

	// download

	const step = 50
	var size = len(cards)

	for i := 0; i < size; i += step {
		fmt.Println(i, "of", size)
		if i+step > size {
			parseJSON(cards[i:])
		} else {
			parseJSON(cards[i : i+step])
		}
	}

	save(resultJSON, config.Db)
}

func parseJSON(ids []string) {
	resp, err := http.PostForm(config.Api,
		url.Values{
			"action":    {"query"},
			"format":    {"json"},
			"redirects": {"1"},
			"prop":      {"revisions"},
			"rvprop":    {"content"},
			"titles":    {strings.Join(ids, "|")},
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var rawData struct {
		Query struct{ Pages json.RawMessage }
	}
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	catch(err)

	var cards WikiaResult
	err = json.Unmarshal(rawData.Query.Pages, &cards)
	catch(err)

	for id, c := range cards {
		resultJSON[id] = c
	}
}
