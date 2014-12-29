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
	"flag"
	"os"
)

var config struct {
	// download

	Api        string
	Categories []string
	Db         string

	// translate

	Cdb      string
	Name     string
	Number   string
	Pendulum string
	Text     string
}

var (
	configFile = flag.String("config", "config.json", "configuration file")
	lang       = flag.String("lang", "en", "output language")
	mainWiki   = flag.Bool("main", false, "always use data from English wikia")
	translate  = flag.Bool("translate", false, "translate the ygopro database")
	update     = flag.Bool("update", false, "update the wikia database")
	ruling     = flag.String("ruling", "", "download cards rulings (specify the filename)")
)

func main() {
	flag.Parse()

	file, err := os.Open(*configFile)
	catch(err)
	defer file.Close()

	var rawConfig map[string]json.RawMessage
	err = json.NewDecoder(file).Decode(&rawConfig)
	catch(err)

	err = json.Unmarshal(rawConfig[*lang], &config)
	catch(err)

	if *mainWiki {
		config.Number = "|number = "
		config.Name = "|" + *lang + "_name = "
		config.Text = "|" + *lang + "_lore = "
		config.Pendulum = "|" + *lang + "_pendulum_effect = "
	}

	if *ruling != "" {
		getRulings()
	}
	if *update {
		wikia()
	}
	if *translate {
		tranlate()
	}
}
