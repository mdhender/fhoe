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

package view

import (
	"net/http"
)

// Handler implements the http interface.
// This should be the only part of the view that you need to customize.
// Mostly that means updating the context.
func (v *View) Handler() http.HandlerFunc {
	var context struct{}
	return func(w http.ResponseWriter, r *http.Request) {
		v.Render(w, r, &context)
	}
}
