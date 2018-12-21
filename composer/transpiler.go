package composer

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/compile"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

type LocatedElement struct {
	Repository  *Repository
	Element     *ElementSource
	PackageName string
	ImportPath  string
	// files fetched from the element directory mapping name to
	// content. Note that file name (key) is the path relative to the
	// path field.
	Files map[string][]byte
}

func NewLocatedElement(repo *Repository, element *ElementSource) *LocatedElement {
	return &LocatedElement{
		Repository:  repo,
		Element:     element,
		PackageName: path.Base(element.Path),
		ImportPath:  repo.ImportPath(element),
	}
}

func (le *LocatedElement) FetchFiles() error {
	if le.Files != nil {
		return nil
	} else if files, err := le.Repository.repo.FetchFiles(le); err != nil {
		return err
	} else {
		le.Files = files
		return nil
	}
}

func (le *LocatedElement) Transpile(c *Composer) error {
	baseDir := filepath.Join(c.Config.OutDir, "src", le.ImportPath)

	fSet := token.NewFileSet()
	anthaFiles := compile.NewAnthaFiles()

	for leaf, content := range le.Files {
		fullPath := filepath.Join(baseDir, leaf)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}

		if err := ioutil.WriteFile(fullPath, content, 0600); err != nil {
			return err
		}

		if filepath.Ext(fullPath) == ".an" {
			if src, err := parser.ParseFile(fSet, fullPath, content, parser.ParseComments); err != nil {
				return err
			} else if antha, err := compile.NewAntha(fSet, src); err != nil {
				return err
			} else {
				for _, ipt := range antha.ImportReqs {
					if err := le.maybeRewriteImport(c, ipt); err != nil {
						return err
					}
				}
				if err := antha.Transform(anthaFiles); err != nil {
					return err
				}
			}
		}
	}
	return writeAnthaFiles(anthaFiles, baseDir)
}

func (le *LocatedElement) maybeRewriteImport(c *Composer, ipt *compile.ImportReq) error {
	// we don't expect imports inside antha files to have revisions
	// within them. So, the strategy is:
	// 1. if it already matches a class with our own prefix and commit, then we're done
	// 2. else if it matches our own source's prefix, then attempt to use our commit
	// 3. Otherwise (and this is most likely), it's just not an import we should be rewriting.

	if strings.HasPrefix(ipt.Path, le.Repository.ImportPrefix) {
		remainingPath := strings.TrimPrefix(ipt.Path, le.Repository.ImportPrefix)
		elem := &ElementSource{
			RepoId: le.Element.RepoId,
			Branch: le.Element.Branch,
			Commit: le.Element.Commit,
			Path:   remainingPath,
		}
		impPath := le.Repository.ImportPath(elem)

		if le2, found := c.classes[impPath]; found { // (1)
			ipt.Path = le2.ImportPath
			return nil
		} else {
			le2 := NewLocatedElement(le.Repository, elem)
			if err := le2.FetchFiles(); err != nil { // (2)
				// even if no files are found, this should not error (at
				// this point we know this source and commit works). So
				// this error should be returned.
				return err
			} else if len(le2.Files) > 0 {
				ipt.Path = le2.ImportPath
				c.EnsureLocatedElement(le2)
				return nil
			}
		}
	}

	return nil
}

func writeAnthaFiles(files *compile.AnthaFiles, baseDir string) error {
	for _, file := range files.Files() {
		outFile := filepath.Join(filepath.Dir(baseDir), filepath.FromSlash(file.Name))
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
