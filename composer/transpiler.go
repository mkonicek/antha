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

func (et *ElementType) Transpile(c *Composer) error {
	baseDir := filepath.FromSlash(path.Join(c.OutDir, "src", et.ImportPath()))

	fSet := token.NewFileSet()
	anthaFiles := compile.NewAnthaFiles()

	for leafPath, content := range et.files {
		fullPath := filepath.Join(baseDir, filepath.FromSlash(leafPath))
		elemDir := filepath.Dir(fullPath)

		if err := os.MkdirAll(elemDir, 0700); err != nil {
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
				et.isAnthaElement = true
				for _, ipt := range antha.ImportReqs {
					if err := et.maybeRewriteImport(c, ipt); err != nil {
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

func (et *ElementType) maybeRewriteImport(c *Composer, ipt *compile.ImportReq) error {
	// we don't expect imports inside antha files to have revisions
	// within them. So, the strategy is:
	// 1. Look for longest matching repo and use that
	// 2. Otherwise (and this is most likely), it's not an import we should be rewriting.

	repoPrefix, repo := c.Workflow.Repositories.LongestMatching(ipt.Path)
	if repo == nil {
		return nil // (2)
	}
	et2 := &ElementType{
		RepositoryPrefix: repoPrefix,
		ElementPath:      ElementPath(strings.TrimPrefix(ipt.Path, string(repoPrefix))),
	}
	if _, err := repo.FetchFiles(et2); err != nil {
		return err
	}
	ipt.Path = et2.ImportPath()
	c.EnsureElementType(et2)

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
