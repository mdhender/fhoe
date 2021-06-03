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

package adr

import (
	"html/template"
	"log"
	"net/http"
)

// Action implements the logic needed to wire domains and responders together.
// It extracts information from the request, packages it up for the domain.
// It takes the payload returned from the domain and passes it to the responder.

// Action has a store (which returns data), a responder (which generates
// the response), and a set of templates used to render the response.
type Action struct {
	store     func(ctx interface{}) (interface{}, error)
	responder *Responder
	templates *template.Template
}

// New returns an initialized action dispatcher.
func New(rsp *Responder, store func(ctx interface{}) (interface{}, error)) *Action {
	return &Action{responder: rsp, store: store}
}

// Handle implements the http.Handler interface for the action.
// It is a generic function.
// It uses the store and responder used to create the action to generate the response.
func (a *Action) Handle(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	log.Printf("[action] %-7s %q\n", r.Method, r.URL.Path)

	payload, err := a.store(ctx)
	if err != nil {
		log.Printf("[action] store: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.responder.Render(payload); err != nil {
		log.Printf("[action] responder %q: %+v\n", a.responder.ID(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// sanity check - the responder should have set a valid http status code.
	if http.StatusText(a.responder.StatusCode()) == "" {
		log.Printf("[action] responder %q: failed to set status code\n", a.responder.ID())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// headers are optional; we can use range since it accepts nil maps
	for k, v := range a.responder.Headers() {
		w.Header().Add(k, v)
	}

	// let the responder set the status code
	w.WriteHeader(a.responder.StatusCode())

	// it may be bad manners (and against the ADR dogma) to check for no content,
	// but it makes me sleep better.
	if a.responder.StatusCode() != http.StatusNoContent && a.responder.Body() != nil {
		w.Write(a.responder.Body())
	}

	log.Printf("[action] %s %q: %q: %d\n", r.Method, r.URL.Path, a.responder.ID(), a.responder.StatusCode())
}
