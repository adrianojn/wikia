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
	"regexp"
	"strings"
)

func getRulings() {
	pages := getListOfPages()
	result := make(map[string]string)
	size := len(pages)

	for i, page := range pages {
		if i%100 == 0 {
			fmt.Println(i, "of", size)
		}
		id, text := getRuling(page)

		if (id != "") && (text != "") {
			result[id] = stripRulingText(text)
		}
	}

	save(result, *ruling)
}

func getRuling(page string) (cardId, cardText string) {
	resp, err := http.PostForm("http://yugioh.wikia.com/api.php",
		url.Values{
			"action":            {"query"},
			"format":            {"json"},
			"prop":              {"revisions"},
			"rvprop":            {"content"},
			"rvexpandtemplates": {"1"},
			"titles":            {"Card Rulings:" + page + "|" + page},
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
			cardText = p.Revisions[0].Text
		}
		if p.Ns == 0 {
			cardId = extractNumber(p.Revisions[0].Text)
		}
	}
	return
}

var removalList = []*regexp.Regexp{
	regex(`(?m)^\{.*`),
	regex(`(?m)^<.*`),
	regex(`(?m)^'.*`),
	regex(`(?m)^\[.*`),
	regex(`(?m)^\|.*`),
	regex(`=+?`),
	regex(`\* ?`),
	regex(`'''`),
	regex(`''`),
	regex(`<ref>?.+?</ref>`),
	regex(`<ref.*?/>`),
	regex(`<sup.+?</sup>`),
	regex(`<span.+?</nowiki>[^<].*?">`),
	regex(`</span>`),
	regex(`Notes`),
	regex(`References`),
}

func stripRulingText(text string) string {
	s := text[strings.Index(text, "=="):]
	for _, re := range removalList {
		s = re.ReplaceAllString(s, "")
	}
	s = wikiRegex.ReplaceAllStringFunc(s, submatch)
	return strings.TrimSpace(s)
}

var numberRegex = regex(`Card Number::(\d+)\|`)

func extractNumber(source string) string {
	s := numberRegex.FindStringSubmatch(source)
	if len(s) > 1 {
		return s[1]
	}
	return ""
}
