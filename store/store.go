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

import "fmt"

func New(root string) (*Store, error) {
	ds := &Store{}
	var err error
	if ds.Planets, err = GetPlanets(root); err != nil {
		return nil, err
	}
	if ds.Ships, err = GetShips(root); err != nil {
		return nil, err
	}
	if ds.Systems, err = GetSystems(root); err != nil {
		return nil, err
	}

	// ensure that systems are sorted by ID
	for _, s := range ds.Systems {
		ds.Sorted.Systems = append(ds.Sorted.Systems, s)
	}
	for i := 0; i < len(ds.Sorted.Systems); i++ {
		for j := i + 1; j < len(ds.Sorted.Systems); j++ {
			if ds.Sorted.Systems[j].ID < ds.Sorted.Systems[i].ID {
				ds.Sorted.Systems[i], ds.Sorted.Systems[j] = ds.Sorted.Systems[j], ds.Sorted.Systems[i]
			}
		}
	}

	// link planets to the owning system
	for _, p := range ds.Planets {
		for _, s := range ds.Sorted.Systems {
			if s.X == p.X && s.Y == p.Y && s.Z == p.Z {
				p.System = s
				break
			}
		}
		if p.System == nil {
			return nil, fmt.Errorf("planet: %q: no such system x: %d y: %d z: %d", p.Name, p.X, p.Y, p.Z)
		}
		p.System.Planets = append(p.System.Planets, p)
	}

	// ensure that planets are sorted by orbit number
	for _, s := range ds.Sorted.Systems {
		for i := 0; i < len(s.Planets); i++ {
			for j := i + 1; j < len(s.Planets); j++ {
				if s.Planets[j].Orbit < s.Planets[i].Orbit {
					s.Planets[i], s.Planets[j] = s.Planets[j], s.Planets[i]
				}
			}
		}
	}

	// add links to systems, planets
	for _, s := range ds.Sorted.Systems {
		s.Link = fmt.Sprintf("/system/%d", s.ID)
	}
	for _, p := range ds.Planets {
		if p.Named {
			p.Link = fmt.Sprintf("/pl/%s", p.Name)
		} else {
			p.Link = fmt.Sprintf("/system/%d/orbit/%d", p.System.ID, p.Orbit)
		}
	}

	return ds, nil
}

type Store struct {
	Planets map[string]*Planet
	Ships   map[string]*Ship
	Systems map[int]*System
	Sorted  struct {
		Systems []*System
	}
}

type Item struct {
	Code     string `json:"code"`
	Quantity int    `json:"qty"`
}

func (ds *Store) SortedSystems() []*System {
	return ds.Sorted.Systems
}
