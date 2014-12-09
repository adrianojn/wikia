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
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const wikiaAPI = "http://yugioh.wikia.com/api.php"

var (
	dbName   = flag.String("db", "cards.cdb", "name of the database file")
	jsonFile = flag.String("json", "result.json", "name of output file")
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

	// load

	db, err := sql.Open("sqlite3", *dbName)
	catch(err)
	defer db.Close()

	rows, err := db.Query("select id, name from texts")
	catch(err)

	ids := make([]string, 0, 10000)
	cardMap := make(map[string]string)

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		catch(err)

		id = fmt.Sprintf("%08s", id)
		ids = append(ids, id)
		cardMap[id] = name
	}
	catch(rows.Err())
	defer rows.Close()

	// download and parse

	const step = 50
	var size = len(ids)

	for i := 0; i < size; i += step {
		fmt.Println(i, "of", size)
		if i+step > size {
			parseJSON(ids[i:])
		} else {
			parseJSON(ids[i : i+step])
		}
	}

	// repeat for missed cards using card name

	size = len(missedCards)
	missed := make([]string, 0, size)

	for _, c := range missedCards {
		missed = append(missed, cardMap[c])
	}

	for i := 0; i < size; i += step {
		fmt.Println("retrying", i, "of", size)
		if i+step > size {
			parseJSON(missed[i:])
		} else {
			parseJSON(missed[i : i+step])
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
		if c.Revisions == nil {
			missedCards = append(missedCards, c.Title)
		} else {
			resultJSON[id] = c
		}
	}
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
