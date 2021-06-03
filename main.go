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

// Package main implements a web server for the order entry application
package main

import (
	"fmt"
	"github.com/mdhender/fhoe/internal/config"
	"github.com/mdhender/fhoe/internal/trace"
	"github.com/mdhender/fhoe/internal/way"
	"github.com/mdhender/fhoe/store"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC) // force logs to be UTC

	// go depends on the operating system to associate extensions with mime-types.
	// the default works mostly for CSS, but this forces it.
	if err := mime.AddExtensionType(".css", "text/css; charset=utf-8"); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}

	cfg := config.Default()
	err := cfg.Load()
	if err == nil {
		err = run(cfg)
	}
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
}

func run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("missing configuration information")
	}

	s := &Server{
		DtFmt:  cfg.TimestampFormat,
		Router: way.NewRouter(),
		t:      trace.Off(),
	}
	s.Addr = net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	s.IdleTimeout = cfg.Server.Timeout.Idle
	s.ReadTimeout = cfg.Server.Timeout.Read
	s.WriteTimeout = cfg.Server.Timeout.Write
	s.MaxHeaderBytes = 1 << 20 // TODO: make this configurable
	s.Handler = s.Router
	s.Templates.Root = cfg.Server.TemplatesRoot

	if cfg.Debug {
		s.t = trace.New(os.Stderr)
	}
	s.site.Title = "Far Horizons Order Entry"

	ds, err := store.New(cfg.App.Data)
	if err != nil {
		return err
	}
	s.ds = ds
	s.t.Tracef("store: loaded %6d systems\n", len(s.ds.Systems))
	s.t.Tracef("store: loaded %6d planets\n", len(s.ds.Planets))
	s.t.Tracef("store: loaded %6d ships\n", len(s.ds.Ships))

	err = s.Routes(cfg, s.t)
	if err != nil {
		return err
	}

	if cfg.Server.TLS.Serve {
		log.Printf("[main] serving TLS on %s\n", s.Addr)
		return s.ListenAndServeTLS(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile)
	}
	log.Printf("[main] listening on %s\n", s.Addr)
	return s.ListenAndServe()
}

type Server struct {
	http.Server
	DtFmt     string // format string for timestamps in responses
	Router    *way.Router
	Templates struct {
		Root string
	}
	debug bool
	ds    *store.Store
	t     trace.Tracer
	site  Site
}
