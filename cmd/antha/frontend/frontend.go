// frontend.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package frontend

import (
	"fmt"
	"io"
	"os"

	"github.com/antha-lang/antha/microArch/logger"
)

// A Frontend is a deprecated object to maintain loging state
type Frontend struct {
	shutdowns []func() error
}

// Shutdown closes any open connections
func (a *Frontend) Shutdown() (err error) {
	for _, fn := range a.shutdowns {
		if e := fn(); e != nil {
			err = e
		}
	}
	return
}

type middleware struct {
	out io.Writer
}

func (a *middleware) Log(lvl logger.LogLevel, ts int64, source, msg string, extra ...interface{}) {
	if lvl == logger.TRACK {
		return
	}
	fmt.Fprint(a.out, msg)      // nolint
	fmt.Fprint(a.out, extra...) // nolint
	fmt.Fprint(a.out, "\n")     // nolint
}

func (a *middleware) Measure(ts int64, source, msg string, extra ...interface{}) {}

func (a *middleware) Sensor(ts int64, source, msg string, extra ...interface{}) {}

func (a *middleware) Data(ts int64, data interface{}, extra ...interface{}) {}

// New creates a new frontend
func New() (*Frontend, error) {
	mw := &middleware{out: os.Stdout}
	logger.RegisterMiddleware(mw)

	ret := &Frontend{}
	ret.shutdowns = append(ret.shutdowns, func() error {
		return nil
	})

	return ret, nil
}
