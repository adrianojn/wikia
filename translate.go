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
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var data map[string]struct {
	Title     string
	Revisions []struct {
		Text string `json:"*"`
	}
}

var db *sql.DB

func tranlate() {
	// load

	f, err := os.Open(config.Db)
	catch(err)
	defer f.Close()

	json.NewDecoder(f).Decode(&data)
	catch(err)

	db, err = sql.Open("sqlite3", config.Cdb)
	catch(err)
	defer db.Close()

	// parse

	for _, card := range data {
		if card.Revisions == nil {
			continue
		}
		text := card.Revisions[0].Text
		id := strings.TrimLeft(extract(text, config.Number), "0")

		var name string
		if *mainWiki || (config.Name != "") {
			name = strip(extract(text, config.Name))
		} else {
			name = card.Title
		}
		lore := strip(extract(text, config.Text))

		if config.Pendulum != "" {
			pendulum := strip(extract(text, config.Pendulum))
			if pendulum != "" {
				lore = "[Pendulum Effect: ]" + pendulum + "\r\n[Monster Effect: ]" + lore
			}
		}

		dbUpdate(id, name, lore)
	}
}

func extract(source, prefix string) string {
	for _, s := range strings.Split(source, "\n") {
		if strings.HasPrefix(s, prefix) {
			return strings.TrimPrefix(s, prefix)
		}
	}
	return ""
}

const updateQuery = `
UPDATE texts SET name=?, desc=? WHERE id IN
  (SELECT id FROM datas WHERE id=? OR
    (alias=? AND id > alias AND id - alias < 10));`

const updatePartialQuery = `
UPDATE texts SET name=? WHERE id IN
  (SELECT id FROM datas WHERE id=? OR
    (alias=? AND id > alias AND id - alias < 10));`

func dbUpdate(id, name, lore string) {
	if id == "" {
		return
	}
	if name == "" {
		fmt.Println("incomplete", id)
		return
	}
	if lore == "" {
		fmt.Println("incomplete", id, name)
		_, err := db.Exec(updatePartialQuery, name, id, id)

		if err != nil {
			fmt.Println(err)
		}
		return
	}

	_, err := db.Exec(updateQuery, name, lore, id, id)
	if err != nil {
		fmt.Println(err)
	}
}

var (
	htmlRegex = regex(`<.+?>`)
	rubyRegex = regex(`\{\{.+?\}\}`)
	wikiRegex = regex(`\[\[.+?\]\]`)
)

func strip(src string) string {
	s := htmlRegex.ReplaceAllString(src, "\n")
	s = wikiRegex.ReplaceAllStringFunc(s, submatch)
	return rubyRegex.ReplaceAllStringFunc(s, submatchRuby)
}

func submatch(s string) string {
	i := strings.Index(s, "|")
	if i < 0 {
		return s[2 : len(s)-2]
	}
	return s[i+1 : len(s)-2]
}

func submatchRuby(s string) string {
	a := strings.Index(s, "|")
	b := strings.LastIndex(s, "|")
	return s[a+1 : b]
}
