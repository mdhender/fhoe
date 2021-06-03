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
	"fmt"
	"github.com/mdhender/fhoe/internal/way"
	"github.com/mdhender/fhoe/store"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet

type Site struct {
	Title string
}

func (s *Server) handleGetCluster() http.HandlerFunc {
	type Context struct {
		Site Site
		Page struct {
			Title string
		}
		Game    string
		Systems []*store.System
	}
	layout := "cluster"
	filename := filepath.Join(s.Templates.Root, layout+".gohtml")
	return func(w http.ResponseWriter, r *http.Request) {
		s.t.Tracef("executing template %q\n", filename)
		t := template.Must(template.ParseFiles(filename))
		var ctx Context
		ctx.Site = s.site
		ctx.Page.Title = "Cluster"
		ctx.Game = "Test"
		ctx.Systems = s.ds.SortedSystems()
		if err := t.Execute(w, &ctx); err != nil {
			s.t.Tracef("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.t.Tracef("executed  template %q\n", filename)
	}
}

func (s *Server) handleGetOrder() http.HandlerFunc {
	type Context struct {
		Site Site
		Page struct {
			Title string
		}
		Game   string
	}
	layout := "orders/index"
	filename := filepath.Join(s.Templates.Root, layout+".gohtml")
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx Context
		ctx.Site = s.site
		ctx.Page.Title = fmt.Sprintf("Turn ?")
		ctx.Game = "Test"

		s.t.Tracef("executing template %q\n", filename)
		t, err := template.ParseFiles(filename)
		if err != nil {
			s.t.Tracef("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: should write to a temp buffer, not directly to the response
		if err := t.Execute(w, &ctx); err != nil {
			s.t.Tracef("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.t.Tracef("executed  template %q\n", filename)
	}
}

func (s *Server) handleGetOrdersHelp() http.HandlerFunc {
	type Context struct {
		Site Site
		Page struct {
			Title string
		}
		Game   string
	}
	layout := "orders/help"
	filename := filepath.Join(s.Templates.Root, layout+".gohtml")
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx Context
		ctx.Site = s.site
		ctx.Page.Title = fmt.Sprintf("Turn ?")
		ctx.Game = "Test"

		s.t.Tracef("executing template %q\n", filename)
		t, err := template.ParseFiles(filename)
		if err != nil {
			s.t.Tracef("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: should write to a temp buffer, not directly to the response
		if err := t.Execute(w, &ctx); err != nil {
			s.t.Tracef("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.t.Tracef("executed  template %q\n", filename)
	}
}

func (s *Server) handleGetSystem() http.HandlerFunc {
	type Context struct {
		Site Site
		Page struct {
			Title string
		}
		Game   string
		System *store.System
	}
	layout := "system"
	filename := filepath.Join(s.Templates.Root, layout+".gohtml")
	return func(w http.ResponseWriter, r *http.Request) {
		val := way.Param(r.Context(), "id")
		id, err := strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		system, ok := s.ds.Systems[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var ctx Context
		ctx.Site = s.site
		ctx.Page.Title = fmt.Sprintf("System %d", id)
		ctx.Game = "Test"
		ctx.System = system

		s.t.Tracef("executing template %q\n", filename)
		t, err := template.ParseFiles(filename)
		if err != nil {
			s.t.Tracef("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: should write to a temp buffer, not directly to the response
		if err := t.Execute(w, &ctx); err != nil {
			s.t.Tracef("%+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.t.Tracef("executed  template %q\n", filename)
	}
}

func (s *Server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("%s: not found\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) handleVersion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}
