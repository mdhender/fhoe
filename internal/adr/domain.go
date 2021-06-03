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

// Domain implements the domain logic.
// It accepts a context from the action, applies business rules to update state (including persisting changes), and returns a payload for the responder to use.

// Domain is the interface for the business logic and data store.
// It's driven by a single function, Payload, that takes a Payload as input and returns a Payload containing the results.
type Domain interface {
	Payload(Payload) Payload
}

// Payload is a key/value store.
// It's used to allow the domain payload to adjust to the requested action.
// I'm sorry that the value is an interface.
// The user will be responsible for mapping it to the expected type using reflection.
// type Payload map[string]interface{}
type Payload struct {
	Site    *Site
	Content interface{}
}

// Site is
type Site struct {
	Title         string
	Slug          string
	Author        string
	CopyrightYear int
}

// Store is
type Store struct {
}

// NewStore returns an initialized store
func NewStore() *Store {
	return &Store{}
}

func (s *Store) HandleIndex(url string) func(ctx interface{}) (interface{}, error) {
	return func(ctx interface{}) (i interface{}, err error) {
		payload := DomainPayload(url)
		return payload, nil
	}
}

// DomainPayload is for testing.
func DomainPayload(url string) Payload {
	var context Payload
	context.Site = &Site{}
	context.Site.Title = "Chas"
	context.Site.Slug = "a writing tool"
	context.Site.Author = "Michael D Henderson"
	context.Site.CopyrightYear = 2020
	return context
}
