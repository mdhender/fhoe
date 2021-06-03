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

// Package trace implements an interface to help with debugging.
package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface that describes an object capable of tracing events throughout the code.
type Tracer interface {
	Trace(a ...interface{})
	Tracef(format string, a ...interface{})
	Write(p []byte) (int, error)
}

type tracer struct {
	w io.Writer
}

func New(w io.Writer) Tracer {
	return &tracer{w: w}
}

func Off() Tracer {
	return &nilTracer{}
}

func (t *tracer) Trace(a ...interface{}) {
	if t == nil {
		return
	} else if _, err := fmt.Fprint(t.w, a...); err != nil {
		panic(err)
	} else if _, err := fmt.Fprintln(t.w); err != nil {
		panic(err)
	}
}

func (t *tracer) Tracef(format string, a ...interface{}) {
	if t == nil {
		return
	} else if _, err := fmt.Fprintf(t.w, format, a...); err != nil {
		panic(err)
	}
}

func (t *tracer) Write(p []byte) (int, error) {
	if t == nil {
		return 0, nil
	} else {
		return t.w.Write(p)
	}
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

func (t *nilTracer) Tracef(format string, a ...interface{}) {}

func (t *nilTracer) Write(p []byte) (int, error) { return 0, nil }
