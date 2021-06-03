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

package controller

import (
	"log"
	"net/http"
	"path"
	"strings"
)

// Handle implements the http interface.
// This should be the only part of the controller that you need to customize for the ADR.
// Mostly that means wiring in the route.
func (ctl *Controller) Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("[controller] %-7s %q\n", r.Method, r.URL.Path)

	// shiftPath splits off the first component of p, which will be cleaned of
	// relative components before processing. head will never contain a slash and
	// tail will always be a rooted path without trailing slash.
	shiftPath := func(p string) (head, tail string) {
		p = path.Clean("/" + p)
		i := strings.Index(p[1:], "/") + 1
		if i <= 0 {
			return p[1:], "/"
		}
		return p[1:i], p[i:]
	}

	originalURL := r.URL.Path
	route, rest := shiftPath(r.URL.Path)

	log.Printf("[controller] %-7s %q %s %s\n", r.Method, r.URL.Path, route, rest)

	switch route {
	case "", "index.html":
		route, rest = shiftPath(rest)
		switch route {
		case "":
			if act, ok := ctl.actions["/"]; ok {
				var context struct{}
				act.Handle(w, r, context)
				return
			}
		}
	}

	log.Printf("[controller] %-7s %q: not matched\n", r.Method, originalURL)
	w.WriteHeader(http.StatusNotFound)
}
