// fhoe - Far Horizons order entry
// Copyright (c) 2021 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type Game struct {
	root    string
	Name    string
	Systems map[int]*System
	Sorted  []*System
	Planets map[string]*Planet
}

type System struct {
	ID      int
	X       int
	Y       int
	Z       int
	Planets []*Planet
}

type Planet struct {
	Name               string
	System             *System
	Orbit              int
	EconomicEfficiency int
	LSN                int
	ProductionPenalty  int
}

func Read(root string) (*Game, error) {
	g := Game{root: root, Name: "test", Systems: make(map[int]*System)}
	systems, err := loadGalaxyListing(filepath.Join(root, "galaxy.list"))
	if err != nil {
		return nil, err
	}
	g.Sorted = systems
	for _, s := range g.Sorted {
		g.Systems[s.ID] = s
	}
	planets, err := loadPlanets(filepath.Join(root, "planets.json"), systems)
	if err != nil {
		return nil, err
	}
	g.Planets = planets
	// ensure that planets are sorted by orbit number
	for _, s := range g.Sorted {
		for i := 0; i < len(s.Planets); i++ {
			for j := i + 1; j < len(s.Planets); j++ {
				if s.Planets[i].Orbit < s.Planets[j].Orbit {
					s.Planets[i], s.Planets[j] = s.Planets[j], s.Planets[i]
				}
			}
		}
	}
	//fmt.Println("{")
	//for _, s := range g.Sorted {
	//	fmt.Printf("  %q: {%q: %d, %q: %d, %q: %d},\n", fmt.Sprintf("%d", s.ID), "x", s.X, "y", s.Y, "z", s.Z)
	//}
	//fmt.Println("}")

	return &g, nil
}

func (g *Game) Write() error {
	panic("!")
}

// input looks like
//   x = 1\ty = 22\tz = 31\tstellar type =  O4
//   The galaxy has a radius of 25 parsecs.
//   It contains 15 dwarf stars, 13 degenerate stars, 120 main sequence stars,
//       and 14 giant stars, for a total of 162 stars.
func loadGalaxyListing(filename string) ([]*System, error) {
	var s []*System
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	for i, line := range bytes.Split(b, []byte{'\n'}) {
		// expect `x = 1\ty = 22\tz = 31\tstellar type =  O4`
		if f := split(line); len(f) == 13 {
			if !(f[0] == "x" && f[3] == "y" && f[6] == "z") {
				return nil, fmt.Errorf("galaxy.list: line %d: unexpected input", i+1)
			} else if x, err := strconv.Atoi(f[2]); err != nil {
				return nil, fmt.Errorf("galaxy.list: line %d: x: %+v", i+1, err)
			} else if y, err := strconv.Atoi(f[5]); err != nil {
				return nil, fmt.Errorf("galaxy.list: line %d: y: %+v", i+1, err)
			} else if z, err := strconv.Atoi(f[8]); err != nil {
				return nil, fmt.Errorf("galaxy.list: line %d: z: %+v", i+1, err)
			} else {
				s = append(s, &System{ID: len(s) + 1, X: x, Y: y, Z: z})
			}
		}
	}
	return s, nil
}

func loadPlanets(filename string, systems []*System) (map[string]*Planet, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data := make(map[string]*struct {
		X                  int    `json:"x"`
		Y                  int    `json:"y"`
		Z                  int    `json:"z"`
		Orbit              int    `json:"orbit"`
		Name               string `json:"name"`
		EconomicEfficiency int    `json:"economic_efficiency"`
		LSN                int    `json:"lsn"`
		ProductionPenalty  int    `json:"production_penalty"`
	})
	if err = json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	planets := make(map[string]*Planet)
	for name, p := range data {
		var planet *Planet
		for _, s := range systems {
			if s.X == p.X && s.Y == p.Y && s.Z == p.Z {
				planet = &Planet{
					EconomicEfficiency: p.EconomicEfficiency,
					LSN:                p.LSN,
					Name:               name,
					Orbit:              p.Orbit,
					ProductionPenalty:  p.ProductionPenalty,
					System:             s,
				}
				s.Planets = append(s.Planets, planet)
				break
			}
		}
		if planet == nil {
			return nil, fmt.Errorf("no such system x: %d y: %d z: %d", p.X, p.Y, p.Z)
		}
		planets[name] = planet
	}
	return planets, nil
}

func split(b []byte) (fields []string) {
	var field []byte
	for len(b) != 0 {
		switch b[0] {
		case '\t', ' ':
			for len(b) != 0 && (b[0] == '\t' || b[0] == ' ') {
				b = b[1:]
			}
			fields = append(fields, string(field))
			field = nil
		case '\r':
			b = nil
		default:
			field = append(field, b[0])
			b = b[1:]
		}
	}
	if field != nil {
		fields = append(fields, string(field))
	}
	return fields
}
