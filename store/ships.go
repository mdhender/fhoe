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

package store

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Ship struct {
	Name     string
	Code     string
	FTL      bool
	Location struct {
		X, Y, Z int
		Orbit   int
		Landed  bool
	}
	Age      int
	Capacity int
	Cargo    []*Item
}

type ShipJSON struct {
	Code     string  `json:"code"`
	Age      int     `json:"age,omitempty"`
	Capacity int     `json:"capacity,omitempty"`
	Cargo    []*Item `json:"cargo,omitempty"`
	FTL      bool    `json:"ftl,omitempty"`
	Location struct {
		X      int  `json:"x"`
		Y      int  `json:"y"`
		Z      int  `json:"z"`
		Orbit  int  `json:"orbit,omitempty"`
		Landed bool `json:"landed,omitempty"`
	}
}

func GetShips(root string) (map[string]*Ship, error) {
	input, err := ioutil.ReadFile(filepath.Join(root, "ships.json"))
	if err != nil {
		return nil, err
	}
	data := make(map[string]*ShipJSON)
	if err = json.Unmarshal(input, &data); err != nil {
		return nil, err
	}
	output := make(map[string]*Ship)
	for k, v := range data {
		s := &Ship{
			Name:     k,
			Code:     v.Code,
			Age:      v.Age,
			FTL:      v.FTL,
			Capacity: v.Capacity,
		}
		for _, c := range v.Cargo {
			s.Cargo = append(s.Cargo, &Item{Code: c.Code, Quantity: c.Quantity})
		}
		s.Location.X = v.Location.X
		s.Location.Y = v.Location.Y
		s.Location.Z = v.Location.Z
		s.Location.Orbit = v.Location.Orbit
		s.Location.Landed = v.Location.Landed
		output[s.Name] = s
	}
	return output, nil
}
