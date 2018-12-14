package composer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/compile"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type LocatedElement struct {
	Source *ElementSource
	// commit shasum or branch name
	Commit string
	// remaining path to the element directory
	Path        string
	PackageName string
	ImportPath  string
	// files fetched from the element directory mapping name to
	// content. Note that file name (key) is the path relative to the
	// path field.
	Files map[string]string
}

func NewLocatedElement(source *ElementSource, commit, remainingPath string) *LocatedElement {
	return &LocatedElement{
		Source:      source,
		Commit:      commit,
		Path:        remainingPath,
		PackageName: path.Base(remainingPath),
		ImportPath:  path.Join(source.Prefix, commit, remainingPath),
	}
}

func (le *LocatedElement) FetchFiles() error {
	if le.Files != nil {
		return nil
	}
	if err := le.Source.ensureRepo(); err != nil {
		return err
	}

	var commitHash plumbing.Hash
	if branch, err := le.Source.repo.Branch(le.Commit); err == git.ErrBranchNotFound {
		// it's not a branch, so assume it's a commit hash
		commitHash = plumbing.NewHash(le.Commit)
	} else if err != nil {
		return err
	} else if ch, err := le.Source.repo.ResolveRevision(plumbing.Revision(branch.Merge)); err != nil {
		return err
	} else {
		commitHash = *ch
	}

	// now follow that commitHash
	if commit, err := le.Source.repo.CommitObject(commitHash); err != nil {
		return err
	} else if tree, err := le.Source.repo.TreeObject(commit.TreeHash); err != nil {
		return err
	} else {
		results := make(map[string]string)
		iter := tree.Files()
		for {
			if f, err := iter.Next(); err == io.EOF {
				break
			} else if err != nil {
				return err
			} else if !strings.HasPrefix(f.Name, le.Path) {
				continue
			} else if c, err := f.Contents(); err != nil {
				return err
			} else {
				results[strings.TrimPrefix(f.Name, le.Path)] = c
			}
		}
		le.Files = results
		return nil
	}
}

func (le *LocatedElement) Transpile(config *Config) error {
	baseDir := filepath.Join(config.OutDir, le.ImportPath)

	root := compile.NewAnthaRoot(le.ImportPath)

	for leaf, content := range le.Files {
		fullPath := filepath.Join(baseDir, leaf)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}

		if filepath.Ext(leaf) == ".an" {
			if err := processFile(root, dir, fullPath, content); err != nil {
				return err
			}
		} else {
			if err := ioutil.WriteFile(fullPath, []byte(content), 0600); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeAnthaFiles(files *compile.AnthaFiles, baseDir string) error {
	for _, file := range files.Files() {
		outFile := filepath.Join(filepath.Dir(baseDir), filepath.FromSlash(file.Name))
		fmt.Println(outFile, file.Name)
		if err := writeAnthaFile(outFile, file); err != nil {
			return err
		}
	}

	return nil
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

// processFile generates the corresponding go code for an antha file.
func processFile(root *compile.AnthaRoot, outdir, filename, content string) error {
	fileSet := token.NewFileSet() // per process FileSet
	file, adjust, err := parse(fileSet, filename, content, false)
	if err != nil {
		return err
	} else if adjust != nil {
		return errNotAnthaFile
	}

	antha := compile.NewAntha(root)

	if err := antha.Transform(fileSet, file); err != nil {
		return err
	}

	files, err := antha.Generate(fileSet, file)
	if err != nil {
		return err
	}

	return writeAnthaFiles(files, outdir)
}

const (
	parserMode = parser.ParseComments
)

var (
	errNotAnthaFile = errors.New("not antha file")
)

// parse parses src, which was read from filename,
// as an Antha source file or statement list.
func parse(fset *token.FileSet, filename string, src string, stdin bool) (*ast.File, func(orig, src []byte) []byte, error) {
	// Try as whole source file.
	file, err := parser.ParseFile(fset, filename, []byte(src), parserMode)
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
