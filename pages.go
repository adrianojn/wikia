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
	"net/http"
	"net/url"
)

type PageResult struct {
	Query    struct{ Categorymembers []struct{ Title string } }
	Continue struct {
		Categorymembers struct{ Cmcontinue string }
	} `json:"query-continue"`
}

func getListOfPages() []string {
	var pages []string
	var next string
	for _, category := range categories {
		for {
			resp, err := http.PostForm(wikiaAPI,
				url.Values{
					"action":     {"query"},
					"format":     {"json"},
					"list":       {"categorymembers"},
					"cmlimit":    {"500"},
					"cmtitle":    {category},
					"cmcontinue": {next},
				})
			catch(err)

			var page PageResult
			err = json.NewDecoder(resp.Body).Decode(&page)
			catch(err)
			resp.Body.Close()

			for _, p := range page.Query.Categorymembers {
				pages = append(pages, p.Title)
			}

			next = page.Continue.Categorymembers.Cmcontinue
			if next == "" {
				break
			}
		}
	}
	return pages
}
