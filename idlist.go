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

var resultJSON map[string]struct {
	Title     string
	Revisions []struct {
		Text string `json:"*"`
	}
}

func main() {
	flag.Parse()

	// load

	db, err := sql.Open("sqlite3", *dbName)
	catch(err)
	defer db.Close()

	rows, err := db.Query("select id from datas")
	catch(err)

	ids := make([]string, 0, 10000)
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		catch(err)
		ids = append(ids, fmt.Sprintf("%08s", id))
	}
	catch(rows.Err())

	// parse

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

	err = json.Unmarshal(rawData.Query.Pages, &resultJSON)
	catch(err)
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
