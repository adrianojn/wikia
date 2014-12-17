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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const wikiaAPI = "http://yugioh.wikia.com/api.php"

var categories = []string{
	"Category:OCG cards",
	"Category:TCG cards",
}

var (
	dbName   = flag.String("db", "cards.cdb", "database file")
	jsonFile = flag.String("json", "wikia.json", "output file")
)

type WikiaResult map[string]struct {
	Title     string
	Revisions []struct {
		Text string `json:"*"`
	}
}

var resultJSON = make(WikiaResult)

var missedCards []string

func main() {
	flag.Parse()

	// download

	cards := getListOfPages()

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

	// save

	out, err := os.Create(*jsonFile)
	catch(err)
	defer out.Close()

	data, err := json.MarshalIndent(&resultJSON, "", " ")
	catch(err)

	io.Copy(out, bytes.NewReader(data))
}

func parseJSON(ids []string) {
	resp, err := http.PostForm(wikiaAPI,
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

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
