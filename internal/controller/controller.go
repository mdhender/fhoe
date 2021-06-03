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
	"github.com/mdhender/fhoe/internal/adr"
	"github.com/mdhender/fhoe/internal/config"
	"github.com/mdhender/fhoe/internal/view"
	"log"
	"net/http"
)

type Controller struct {
	store       *adr.Store
	actions     map[string]*adr.Action
	indexAction *adr.Action
}

// New returns
func New(cfg *config.Config, store *adr.Store) (*Controller, error) {
	ctl := Controller{
		actions: make(map[string]*adr.Action),
	}

	for _, route := range []struct {
		url, tmplt string
		fn         func(url string) func(ctx interface{}) (interface{}, error)
	}{
		{"/", "/page/index.gohtml", store.HandleIndex},
		//{"/assumptions", "/page/assumptions.gohtml", store.HandleAssumptions},
	} {
		v, err := view.New("site", cfg.Server.TemplatesRoot+"/site.gohtml", route.tmplt, cfg.Server.TemplatesRoot)
		if err != nil {
			log.Printf("[controller] view: %+v\n", err)
			return nil, err
		}
		rsp, err := adr.NewResponder(route.url, route.tmplt, v)
		if err != nil {
			return nil, err
		}
		ctl.actions[route.url] = adr.New(rsp, route.fn(route.url))
	}

	return &ctl, nil
}

// ServeHTTP to satisfy the http.Handler interface. It forwards all requests
// to our custom handle function.
func (ctl *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctl.Handle(w, r)
}

// Handler implements the http interface by forwarding all requests to custom
// handle function.
func (ctl *Controller) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctl.Handle(w, r)
	})
}
