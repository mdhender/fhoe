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

// Package templates implements a handler for template.Template
package templates

import (
	"github.com/mdhender/fhoe/internal/trace"
	"html/template"
	"net/http"
	"path/filepath"
)

func New(layout, root string, tracer trace.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Join(root, layout+".gohtml")
		tracer.Tracef("executing template %q\n", filename)
		t := template.Must(template.ParseFiles(filename))
		var ctx struct {
			Site struct {
				Title string
			}
		}
		ctx.Site.Title = "joe"
		if err := t.Execute(w, &ctx); err != nil {
			tracer.Tracef("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tracer.Tracef("executed  template %q\n", filename)
	}
}
