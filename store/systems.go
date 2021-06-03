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
	"strconv"
)

type System struct {
	ID      int
	X, Y, Z int
	Link    string
	Planets []*Planet
}

type SystemJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"Z"`
}

func GetSystems(root string) (map[int]*System, error) {
	input, err := ioutil.ReadFile(filepath.Join(root, "systems.json"))
	if err != nil {
		return nil, err
	}
	data := make(map[string]*SystemJSON)
	if err = json.Unmarshal(input, &data); err != nil {
		return nil, err
	}
	output := make(map[int]*System)
	for k, v := range data {
		id, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("store: systems: %q: %+v", k, err)
		}
		output[id] = &System{
			ID: id,
			X:  v.X,
			Y:  v.Y,
			Z:  v.Z,
		}
	}
	return output, nil
}
