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
)

func getRulings() {
	pages := getListOfPages()[:10]

	result := make(map[string]string)

	for _, page := range pages {
		id, text := getRuling(page)

		if id != "" {
			result[id] = text
		}
	}

	save(result, *ruling)
}

func getRuling(page string) (cardId, cardText string) {
	resp, err := http.PostForm("http://yugioh.wikia.com/api.php",
		url.Values{
			"action": {"query"},
			"format": {"json"},
			"prop":   {"revisions"},
			"rvprop": {"content"},
			"titles": {"Card Rulings:" + page + "|" + page},
		})
	if err != nil {
		fmt.Println(page, err)
		return
	}
	var data struct {
		Query struct {
			Pages map[string]struct {
				Ns        int
				Revisions []struct {
					Text string `json:"*"`
				}
			}
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(page, err)
		return
	}

	for _, p := range data.Query.Pages {
		if p.Revisions == nil {
			return
		}
		if p.Ns == 102 {
			// TODO: strip text
			cardText = p.Revisions[0].Text

		}
		if p.Ns == 0 {
			cardId = extract(p.Revisions[0].Text, "|number = ")
		}
	}
	return
}
