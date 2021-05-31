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

package config

import (
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Debug bool
	App   struct {
		Root            string
		Data            string
		TimestampFormat string
	}
	TimestampFormat string
	Server          struct {
		Scheme  string
		Host    string
		Port    string
		Timeout struct {
			Idle  time.Duration
			Read  time.Duration
			Write time.Duration
		}
		TLS struct {
			Serve    bool
			CertFile string
			KeyFile  string
		}
		Salt          string
		Key           string
		TemplatesRoot string
		WebRoot       string
	}
}

// Default returns a default configuration.
// These are the values without loading the environment, configuration file, or command line.
func Default() *Config {
	var cfg Config
	cfg.App.Root = "D:\\GoLand\\fhoe"
	cfg.App.Data = filepath.Join(cfg.App.Root, "testdata")
	cfg.App.TimestampFormat = "2006-01-02T15:04:05.99999999Z"
	cfg.Server.Scheme = "http"
	cfg.Server.Host = "localhost"
	cfg.Server.Port = "3000"
	cfg.Server.Timeout.Idle = 10 * time.Second
	cfg.Server.Timeout.Read = 5 * time.Second
	cfg.Server.Timeout.Write = 10 * time.Second
	cfg.Server.Key = "curry.aka.yrruc"
	cfg.Server.Salt = "pepper"
	cfg.Server.TemplatesRoot = filepath.Join(cfg.App.Root, "templates")
	cfg.Server.WebRoot = filepath.Join(cfg.App.Root, "web")
	return &cfg
}

// Load updates the values in a Config in this order:
//   1. It will load a configuration file if one is given on the
//      command line via the `-config` flag. If provided, the file
//      must contain a valid JSON object.
//   2. Environment variables, using the prefix `CONDUIT_RYER_SERVER`
//   3. Command line flags
func (cfg *Config) Load() error {
	fs := flag.NewFlagSet("Server", flag.ExitOnError)
	debug := fs.Bool("debug", cfg.Debug, "log debug information (optional)")
	appRoot := fs.String("root", cfg.App.Root, "path to treat as root for relative file references")
	appData := fs.String("data", cfg.App.Data, "path containing data files")
	serverScheme := fs.String("scheme", cfg.Server.Scheme, "http scheme, either 'http' or 'https'")
	serverHost := fs.String("host", cfg.Server.Host, "host name (or IP) to listen on")
	serverPort := fs.String("port", cfg.Server.Port, "port to listen on")
	serverKey := fs.String("key", cfg.Server.Key, "set key for signing tokens")
	serverSalt := fs.String("salt", cfg.Server.Salt, "set salt for hashing")
	serverTimeoutIdle := fs.Duration("idle-timeout", cfg.Server.Timeout.Idle, "http idle timeout")
	serverTimeoutRead := fs.Duration("read-timeout", cfg.Server.Timeout.Read, "http read timeout")
	serverTimeoutWrite := fs.Duration("write-timeout", cfg.Server.Timeout.Write, "http write timeout")
	serverTLSServe := fs.Bool("https", cfg.Server.TLS.Serve, "serve https")
	serverTLSCertFile := fs.String("https-cert-file", cfg.Server.Host, "https certificate file")
	serverTLSKeyFile := fs.String("https-key-file", cfg.Server.Host, "https certificate key file")
	serverTemplatesRoot := fs.String("templates", cfg.Server.WebRoot, "path to fetch templates from")
	serverWebRoot := fs.String("web", cfg.Server.WebRoot, "path to serve web assets from")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("FHOE"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.JSONParser)); err != nil {
		return err
	}

	cfg.Debug = *debug
	cfg.App.Root = filepath.Clean(*appRoot)
	cfg.App.Data = filepath.Clean(*appData)
	cfg.Server.Scheme = *serverScheme
	cfg.Server.Host = *serverHost
	cfg.Server.Port = *serverPort
	cfg.Server.Key = *serverKey
	cfg.Server.Salt = *serverSalt
	cfg.Server.Timeout.Idle = *serverTimeoutIdle
	cfg.Server.Timeout.Read = *serverTimeoutRead
	cfg.Server.Timeout.Write = *serverTimeoutWrite
	cfg.Server.TLS.Serve = *serverTLSServe
	cfg.Server.TLS.CertFile = *serverTLSCertFile
	cfg.Server.TLS.KeyFile = *serverTLSKeyFile
	cfg.Server.TemplatesRoot = filepath.Clean(*serverTemplatesRoot)
	cfg.Server.WebRoot = filepath.Clean(*serverWebRoot)

	if cfg.Server.TLS.Serve == true {
		if cfg.Server.TLS.CertFile == "" {
			return fmt.Errorf("must supply certificates file when serving HTTPS")
		}
		if cfg.Server.TLS.KeyFile == "" {
			return fmt.Errorf("must supply certificate key file when serving HTTPS")
		}
	}

	return nil
}
