// compile.go: Part of the Antha language
// Copyright (C) 2017 The Antha authors. All rights reserved.
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

package cmd

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/compile"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	parserMode = parser.ParseComments
)

var (
	errNotAnthaFile = errors.New("not antha file")
)

var compileCmd = &cobra.Command{
	Use:   "compile <files or directories>",
	Short: "Compile antha elements",
	RunE:  runCompile,
}

func runCompile(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	outdir := viper.GetString("outdir")
	outPackage := viper.GetString("outputPackage")

	if len(outdir) == 0 {
		return fmt.Errorf("missing outdir")
	}

	// parse every filename or directory passed in as input
	root := compile.NewAnthaRoot(outPackage)

	var errs []error

	for _, path := range args {
		if err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if f.IsDir() {
				return nil
			}

			if !isAnthaFile(f.Name()) {
				return nil
			}

			// Collect errors processing errors
			if err := processFile(root, path, outdir); err != nil {
				errs = append(
					errs,
					fmt.Errorf("error processing file: %s \n Error: %s", path, err),
				)
			}

			return nil
		}); err != nil {
			return err
		}
	}

	for _, err := range errs {
		fmt.Println(err)
	}

	if len(errs) != 0 {
		return fmt.Errorf("some files did not compile")
	}

	files, err := root.Generate()
	if err != nil {
		return err
	}

	if err := writeAnthaFiles(files, outdir); err != nil {
		return err
	}

	return nil
}

// isAnthaFile returns if file matches antha file naming convention
func isAnthaFile(name string) bool {
	return strings.HasSuffix(name, ".an")
}

func writeAnthaFile(outFile string, file *compile.AnthaFile) error {
	dst, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer dst.Close() // nolint: errcheck

	src := file.NewReader()
	defer src.Close() // nolint: errcheck

	_, err = io.Copy(dst, src)
	return err
}

func writeAnthaFiles(files *compile.AnthaFiles, outDir string) error {
	for _, file := range files.Files() {
		outFile := filepath.Join(outDir, filepath.FromSlash(file.Name))
		if err := os.MkdirAll(filepath.Dir(outFile), 0700); err != nil {
			return err
		}

		if err := writeAnthaFile(outFile, file); err != nil {
			return err
		}
	}

	return nil
}

// processFile generates the corresponding go code for an antha file.
func processFile(root *compile.AnthaRoot, filename, outdir string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	fileSet := token.NewFileSet() // per process FileSet
	file, adjust, err := parse(fileSet, filename, src, false)
	if err != nil {
		return err
	} else if adjust != nil {
		return errNotAnthaFile
	}

	h := sha256.New()
	if _, err := io.Copy(h, bytes.NewReader(src)); err != nil {
		return err
	}

	antha := compile.NewAntha(root)
	antha.SourceSHA256 = h.Sum(nil)

	if err := antha.Transform(fileSet, file); err != nil {
		return err
	}

	files, err := antha.Generate(fileSet, file)
	if err != nil {
		return err
	}

	return writeAnthaFiles(files, outdir)
}

// parse parses src, which was read from filename,
// as an Antha source file or statement list.
func parse(fset *token.FileSet, filename string, src []byte, stdin bool) (*ast.File, func(orig, src []byte) []byte, error) {
	// Try as whole source file.
	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err == nil {
		return file, nil, nil
	}
	// If the error is that the source file didn't begin with a
	// package line and this is standard input, fall through to
	// try as a source fragment.  Stop and return on any other error.
	if !stdin || !strings.Contains(err.Error(), "expected 'package'") {
		return nil, nil, err
	}

	// If this is a declaration list, make it a source file
	// by inserting a package clause.
	// Insert using a ;, not a newline, so that the line numbers
	// in psrc match the ones in src.
	psrc := append([]byte("protocol p;"), src...)
	file, err = parser.ParseFile(fset, filename, psrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			// Remove the package clause.
			// Anthafmt has turned the ; into a \n.
			src = src[len("protocol p\n"):]
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	// If the error is that the source file didn't begin with a
	// declaration, fall through to try as a statement list.
	// Stop and return on any other error.
	if !strings.Contains(err.Error(), "expected declaration") {
		return nil, nil, err
	}

	// If this is a statement list, make it a source file
	// by inserting a package clause and turning the list
	// into a function body.  This handles expressions too.
	// Insert using a ;, not a newline, so that the line numbers
	// in fsrc match the ones in src.
	fsrc := append(append([]byte("protocol p; func _() {"), src...), '}')
	file, err = parser.ParseFile(fset, filename, fsrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			// Remove the wrapping.
			// Anthafmt has turned the ; into a \n\n.
			src = src[len("protocol p\n\nfunc _() {"):]
			src = src[:len(src)-len("}\n")]
			// Anthafmt has also indented the function body one level.
			// Remove that indent.
			src = bytes.Replace(src, []byte("\n\t"), []byte("\n"), -1)
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}

	// Failed, and out of options.
	return nil, nil, err
}

// Utility function for matchSpace
func cutSpace(b []byte) (before, middle, after []byte) {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n') {
		i++
	}
	j := len(b)
	for j > 0 && (b[j-1] == ' ' || b[j-1] == '\t' || b[j-1] == '\n') {
		j--
	}
	if i <= j {
		return b[:i], b[i:j], b[j:]
	}
	return nil, nil, b[j:]
}

// matchSpace reformats src to use the same space context as orig.
// 1) If orig begins with blank lines, matchSpace inserts them at the beginning of src.
// 2) matchSpace copies the indentation of the first non-blank line in orig
//    to every non-blank line in src.
// 3) matchSpace copies the trailing space from orig and uses it in place
//   of src's trailing space.
func matchSpace(orig []byte, src []byte) []byte {
	before, _, after := cutSpace(orig)
	i := bytes.LastIndex(before, []byte{'\n'})
	before, indent := before[:i+1], before[i+1:]

	_, src, _ = cutSpace(src)

	var b bytes.Buffer
	b.Write(before)
	for len(src) > 0 {
		line := src
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, src = line[:i+1], line[i+1:]
		} else {
			src = nil
		}
		if len(line) > 0 && line[0] != '\n' { // not blank
			b.Write(indent) // nolint: errcheck
		}
		b.Write(line)
	}
	b.Write(after) // nolint: errcheck

	return b.Bytes()
}

func init() {
	c := compileCmd
	flags := c.Flags()
	RootCmd.AddCommand(c)

	flags.String("outdir", "", "output directory for generated files")
	flags.String("outputPackage", "", "base package name for generated files")
}
