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

package main

import (
	"github.com/mdhender/fhoe/internal/config"
	"github.com/mdhender/fhoe/internal/trace"
	"net/http"
)

// Routes initializes all routes exposed by the Server.
// Routes are taken from https://github.com/gothinkster/realworld/blob/9686244365bf5681e27e2e9ea59a4d905d8080db/api/swagger.json
func (s *Server) Routes(cfg *config.Config, t trace.Tracer) error {
	// assign routes handled by the api
	for _, route := range []struct {
		pattern string
		method  string
		handler http.HandlerFunc
	}{
		{"/api/status/live", "GET", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}},
		{"/api/status/ready", "GET", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}},
		{"/api/version", "GET", s.handleVersion()},
		{"/cluster", "GET", s.handleGetCluster()},
		{"/orders", "GET", s.handleGetOrder()},
		{"/orders/help", "GET", s.handleGetOrdersHelp()},
		{"/system/:id", "GET", s.handleGetSystem()},
	} {
		s.Router.HandleFunc(route.method, route.pattern, route.handler)
	}

	s.Router.NotFound = http.FileServer(http.Dir(cfg.Server.WebRoot))

	return nil
}
