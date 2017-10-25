// antha/scanner/errors.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
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

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/antha-lang/antha/antha/token"
)

// An Error in scanning. In an ErrorList, an error is represented by an *Error.
// The position Pos, if valid, points to the beginning of the offending token,
// and the error condition is described by Msg.
//
type Error struct {
	Pos token.Position
	Msg string
}

// Error implements the error interface.
func (e Error) Error() string {
	if e.Pos.Filename != "" || e.Pos.IsValid() {
		// don't print "<unknown position>"
		// TODO(gri) reconsider the semantics of Position.IsValid
		return e.Pos.String() + ": " + e.Msg
	}
	return e.Msg
}

// ErrorList is a list of *Errors.
// The zero value for an ErrorList is an empty ErrorList ready to use.
//
type ErrorList []*Error

// Add adds an Error with given position and error message to an ErrorList.
func (p *ErrorList) Add(pos token.Position, msg string) {
	*p = append(*p, &Error{pos, msg})
}

// Reset resets an ErrorList to no errors.
func (p *ErrorList) Reset() { *p = (*p)[0:0] }

// ErrorList implements the sort Interface.
func (p ErrorList) Len() int      { return len(p) }
func (p ErrorList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p ErrorList) Less(i, j int) bool {
	e := &p[i].Pos
	f := &p[j].Pos
	// Note that it is not sufficient to simply compare file offsets because
	// the offsets do not reflect modified line information (through //line
	// comments).
	if e.Filename < f.Filename {
		return true
	}
	if e.Filename == f.Filename {
		if e.Line < f.Line {
			return true
		}
		if e.Line == f.Line {
			return e.Column < f.Column
		}
	}
	return false
}

// Sort sorts an ErrorList. *Error entries are sorted by position,
// other errors are sorted by error message, and before any *Error
// entry.
//
func (p ErrorList) Sort() {
	sort.Sort(p)
}

// RemoveMultiples sorts an ErrorList and removes all but the first error per line.
func (p *ErrorList) RemoveMultiples() {
	sort.Sort(p)
	var last token.Position // initial last.Line is != any legal error line
	i := 0
	for _, e := range *p {
		if e.Pos.Filename != last.Filename || e.Pos.Line != last.Line {
			last = e.Pos
			(*p)[i] = e
			i++
		}
	}
	(*p) = (*p)[0:i]
}

// An ErrorList implements the error interface.
func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}
	PrintError(os.Stderr, p)
	return fmt.Sprintf("%s (and %d more errors)", p[0], len(p)-1)
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (p ErrorList) Err() error {
	if len(p) == 0 {
		return nil
	}
	return p
}

// PrintError is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an ErrorList. Otherwise
// it prints the err string.
//
func PrintError(w io.Writer, err error) {
	if list, ok := err.(ErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e) // nolint
		}
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err) // nolint
	}
}
