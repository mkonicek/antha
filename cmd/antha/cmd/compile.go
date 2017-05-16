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

package cmd

import (
	"bytes"
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

// execution variables
var (
	fileSet = token.NewFileSet() // per process FileSet
)

// parameters to control code formatting
const (
	tabWidth    = 8
	printerMode = compile.UseSpaces | compile.TabIndent
	parserMode  = parser.ParseComments
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile antha element",
	RunE:  runCompile,
}

func runCompile(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	o := output{
		OutDir: viper.GetString("outdir"),
	}
	if err := o.Init(); err != nil {
		return err
	}
	defer o.Close()

	// try to parse standard input if no files or directories were passed in
	if len(args) == 0 {
		if err := processFile(processFileOptions{
			Filename: "-",
			In:       os.Stdin,
			Stdin:    true,
			OutDir:   o.Dir(),
		}); err != nil {
			return err
		}
	}

	// parse every filename or directory passed in as input
	for _, path := range args {
		switch dir, err := os.Stat(path); {
		case err != nil:
			return err
		case dir.IsDir():
			filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
				// Ignore previous errors
				if isAnthaFile(f) {
					// TODO this might be an issue since we have to analyse the contents in
					// order to establish whether more than one component exist
					err = processFile(processFileOptions{
						Filename: path,
						OutDir:   o.Dir(),
					})
					if err != nil {
						return err
					}
				}
				return err
			})
		default:
			if err := processFile(processFileOptions{
				Filename: path,
				OutDir:   o.Dir(),
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// Utility function to check file extension
func isAnthaFile(f os.FileInfo) bool {
	// ignore non-Antha or Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".an")
}

// Remove files from dir with suffix
func removeFiles(dir, suffix string) error {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, v := range fis {
		if !strings.HasSuffix(v.Name(), suffix) {
			continue
		}
		if err := os.RemoveAll(filepath.Join(dir, v.Name())); err != nil {
			return err
		}
	}
	return nil
}

type output struct {
	OutDir  string
	outName string
	dir     string
}

func (a *output) Dir() string {
	if len(a.OutDir) == 0 {
		return a.dir
	}
	return a.OutDir
}

func (a *output) Close() error {
	if len(a.dir) == 0 {
		return nil
	}
	return os.RemoveAll(a.dir)
}

func (a *output) Init() error {
	if len(a.OutDir) == 0 {
		n, err := ioutil.TempDir("", "antha")
		if err != nil {
			return err
		}
		a.dir = n
		return nil
	}

	p, err := filepath.Abs(a.OutDir)
	if err != nil {
		return err
	}
	a.outName = filepath.Base(p)
	if err := removeFiles(p, "_.go"); err != nil {
		return err
	}
	return nil
}

// Write out a file. Makes sure that output directory exists as well.
func write(fname string, bs []byte) error {
	if err := os.MkdirAll(filepath.Dir(fname), 0777); err != nil {
		return err
	}
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewBuffer(bs))
	return err
}

type processFileOptions struct {
	Filename string
	In       io.Reader
	Stdin    bool
	OutDir   string // empty string means output to same directory as Filename
}

// If in == nil, the source is the contents of the file with the given
// filename.
func processFile(opt processFileOptions) error {
	if opt.In == nil {
		f, err := os.Open(opt.Filename)
		if err != nil {
			return err
		}
		defer f.Close()
		opt.In = f
	}

	src, err := ioutil.ReadAll(opt.In)
	if err != nil {
		return err
	}

	file, adjust, err := parse(fileSet, opt.Filename, src, opt.Stdin)
	if err != nil {
		return err
	}

	if file.Tok != token.PROTOCOL {
		return fmt.Errorf("%s is not a valid Antha file", opt.Filename)
	}
	// Extract protocol name
	compName := file.Name.Name

	var buf bytes.Buffer
	compiler := &compile.Config{Mode: printerMode, Tabwidth: tabWidth, Package: "main"}
	if err := compiler.Fprint(&buf, fileSet, file); err != nil {
		return err
	}
	res := buf.Bytes()
	if adjust != nil {
		res = adjust(src, res)
	}

	dir := opt.OutDir
	if len(dir) == 0 {
		dir = filepath.Dir(opt.Filename)
	}

	outFile := filepath.Join(dir, fmt.Sprintf("%s_.go", compName))
	if err := write(outFile, res); err != nil {
		return err
	}

	return err
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
			b.Write(indent)
		}
		b.Write(line)
	}
	b.Write(after)
	return b.Bytes()
}

func init() {
	c := compileCmd
	flags := c.Flags()
	RootCmd.AddCommand(c)

	flags.String("outdir", "", "output directory for generated files")
}
