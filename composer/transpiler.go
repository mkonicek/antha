package composer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/compile"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
	"github.com/antha-lang/antha/workflow"
)

type TranspilableElementType struct {
	*workflow.ElementType
	// An element can import code from elsewhere within its repository,
	// and we cope with that just fine, creating new instances of
	// ElementType to track this. But need to keep an eye out for
	// whether the files we find within such a package/directory
	// actually contain .an files, because we can only assume certain
	// functions (eg RegisterLineMap) exist for packages which really
	// contain elements.
	transpiler *compile.Antha
}

func NewTranspilableElementType(et *workflow.ElementType) *TranspilableElementType {
	return &TranspilableElementType{
		ElementType: et,
	}
}

func (tet TranspilableElementType) IsAnthaElement() bool {
	return tet.transpiler != nil
}

func (tet *TranspilableElementType) Transpile(c *Composer) error {
	baseDir := filepath.Join(c.OutDir, "src", filepath.FromSlash(tet.ImportPath()))

	anthaFiles := compile.NewAnthaFiles()
	anthaFound := false
	err := filepath.Walk(baseDir, func(p string, info os.FileInfo, err error) error {
		if err != nil || !info.Mode().IsRegular() || !workflow.IsAnthaFile(p) {
			return err
		} else if anthaFound {
			return fmt.Errorf("Multiple .an files found in %v", baseDir)
		}
		anthaFound = true
		c.Logger.Log("transpiling", tet.ImportPath())
		if bs, err := ioutil.ReadFile(p); err != nil {
			return err
		} else if err := tet.CreateTranspiler(bs, p); err != nil {
			return err
		} else {
			for _, ipt := range tet.transpiler.ImportReqs {
				if err := tet.maybeRewriteImport(c, ipt); err != nil {
					return err
				}
			}
			return tet.transpiler.Transform(anthaFiles)
		}
	})
	if err != nil {
		return err
	} else {
		return writeAnthaFiles(anthaFiles, filepath.Dir(baseDir))
	}
}

func (tet *TranspilableElementType) Meta(bs []byte, path string) (*compile.Meta, error) {
	if err := tet.CreateTranspiler(bs, path); err != nil {
		return nil, err
	} else {
		return tet.transpiler.Meta()
	}
}

// path is deliberately separate from bs because path might be some
// symbolic name unrelated to the actual source of the file (eg think
// some git repo). I.e. we can't just do an ioutil.ReadFile on path.
func (tet *TranspilableElementType) CreateTranspiler(bs []byte, path string) error {
	fSet := token.NewFileSet()
	if src, err := parser.ParseFile(fSet, path, bs, parser.ParseComments); err != nil {
		return err
	} else if antha, err := compile.NewAntha(fSet, src); err != nil {
		return err
	} else {
		tet.transpiler = antha
		return nil
	}
}

func (tet *TranspilableElementType) maybeRewriteImport(c *Composer, ipt *compile.ImportReq) error {
	// we don't expect imports inside antha files to have revisions
	// within them (or specify repositories in any non-standard way). So, the strategy is:
	// 1. Look for longest matching repo and use that
	// 2. Otherwise (and this is most likely), it's not an import we should be rewriting.

	repoPrefix, repo := c.Workflow.Repositories.LongestMatching(ipt.Path)
	if repo == nil {
		return nil // (2)
	}
	tet2 := NewTranspilableElementType(&workflow.ElementType{
		RepositoryPrefix: repoPrefix,
		ElementPath:      workflow.ElementPath(strings.TrimPrefix(ipt.Path, string(repoPrefix))),
	})
	if err := repo.Clone(filepath.Join(c.OutDir, "src")); err != nil {
		return err
	}
	ipt.Path = tet2.ImportPath()
	c.EnsureElementType(tet2)

	return nil
}

func writeAnthaFiles(files *compile.AnthaFiles, baseDir string) error {
	for _, file := range files.Files() {
		outFile := filepath.Join(baseDir, filepath.FromSlash(file.Name))
		if err := writeAnthaFile(outFile, file); err != nil {
			return err
		}
	}

	return nil
}

func writeAnthaFile(outFile string, file *compile.AnthaFile) error {
	dst, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400)
	if err != nil {
		return err
	}
	defer dst.Close() // nolint: errcheck

	src := file.NewReader()
	defer src.Close() // nolint: errcheck

	_, err = io.Copy(dst, src)
	return err
}
