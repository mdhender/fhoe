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
	"bytes"
	"github.com/mdhender/fhoe/internal/view"
	"log"
	"net/http"
)

// Responder implements the presentation logic for the application.
// It accepts a context from the action (which is the domain payload),
// sets the HTTP status and headers, and writes the body of the response.
type Responder struct {
	id         string
	yield      string
	statusCode int
	headers    map[string]string
	buf        *bytes.Buffer
	body       []byte
	view       *view.View
}

// NewResponder is
func NewResponder(id string, yield string, v *view.View) (*Responder, error) {
	return &Responder{
		id:      id,
		yield:   yield,
		headers: make(map[string]string),
		buf:     bytes.NewBuffer(nil),
		view:    v,
	}, nil
}

// Render is
func (r *Responder) Render(ctx interface{}) error {
	log.Printf("[responder] render\n")
	r.buf = bytes.NewBuffer(nil)
	t, err := r.view.Load()
	if err != nil {
		log.Printf("[responder] render load: %+v\n", err)
		r.statusCode = http.StatusInternalServerError
		return err
	}
	if err = t.ExecuteTemplate(r.buf, r.view.Name, ctx); err != nil {
		log.Printf("[responder] render execute: %+v\n", err)
		r.statusCode = http.StatusInternalServerError
		return err
	}
	log.Printf("[responder] render %d bytes\n", r.buf.Len())
	r.statusCode = http.StatusOK
	return nil
}

// Body returns a slice of bytes built using the context, or nil if the response body should be empty.
func (r *Responder) Body() []byte {
	return r.buf.Bytes()
}

// Headers returns either nil or a map of key/value pairs
func (r *Responder) Headers() map[string]string {
	return r.headers
}

// ID returns the identifier for this responder.
func (r *Responder) ID() string {
	return r.id
}

// StatusCode returns the http status code for this responder.
func (r *Responder) StatusCode() int {
	return r.statusCode
}
