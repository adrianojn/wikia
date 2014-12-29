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
	"io"
	"os"
	"regexp"
)

var regex = regexp.MustCompile

func save(data interface{}, fileName string) {
	out, err := os.Create(fileName)
	catch(err)
	defer out.Close()

	jsonData, err := json.MarshalIndent(data, "", " ")
	catch(err)

	_, err = io.Copy(out, bytes.NewReader(jsonData))
	catch(err)
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
