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
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Planet struct {
	Name               string
	System             *System
	X                  int
	Y                  int
	Z                  int
	Orbit              int
	EconomicEfficiency int
	Inventory          []*Item
	Link               string
	LSN                int
	MiningDifficulty   float64
	Named              bool
	ProductionPenalty  int
}

type PlanetJSON struct {
	X                  int     `json:"x"`
	Y                  int     `json:"y"`
	Z                  int     `json:"z"`
	Orbit              int     `json:"orbit"`
	EconomicEfficiency int     `json:"economic_efficiency"`
	Inventory          []*Item `json:"inventory,omitempty"`
	LSN                int     `json:"lsn"`
	MiningDifficulty   float64 `json:"mining_difficulty"`
	ProductionPenalty  int     `json:"production_penalty,omitempty"`
}

func GetPlanets(root string) (map[string]*Planet, error) {
	input, err := ioutil.ReadFile(filepath.Join(root, "planets.json"))
	if err != nil {
		return nil, err
	}
	data := make(map[string]*PlanetJSON)
	if err = json.Unmarshal(input, &data); err != nil {
		return nil, err
	}
	output := make(map[string]*Planet)
	for k, v := range data {
		p := &Planet{
			Name:               k,
			System:             nil,
			X:                  v.X,
			Y:                  v.Y,
			Z:                  v.Z,
			Orbit:              v.Orbit,
			EconomicEfficiency: v.EconomicEfficiency,
			Inventory:          v.Inventory,
			LSN:                v.LSN,
			MiningDifficulty:   v.MiningDifficulty,
			Named:              k != fmt.Sprintf("%d %d %d %d", v.X, v.Y, v.Z, v.Orbit),
			ProductionPenalty:  v.ProductionPenalty,
		}
		output[p.Name] = p
	}
	return output, nil
}
