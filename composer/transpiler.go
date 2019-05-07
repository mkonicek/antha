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
	"github.com/antha-lang/antha/utils"
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
	Transpiler *compile.Antha
}

func NewTranspilableElementType(et *workflow.ElementType) *TranspilableElementType {
	return &TranspilableElementType{
		ElementType: et,
	}
}

func (tet TranspilableElementType) IsAnthaElement() bool {
	return tet.Transpiler != nil
}

// This only works based on the elements already having been cloned to
// the c.OutDir. I.e. this does not work Repositories (eg git etc)
func (tet *TranspilableElementType) TranspileFromFS(cb *ComposerBase, wf *workflow.Workflow) error {
	baseDir := filepath.Join(cb.OutDir, "src", filepath.FromSlash(tet.ImportPath()))

	anthaFound := false
	anthaFiles := compile.NewAnthaFiles()
	err := filepath.Walk(baseDir, func(p string, info os.FileInfo, err error) error {
		if err != nil || !info.Mode().IsRegular() || !workflow.IsAnthaFile(p) {
			return err
		} else if anthaFound {
			return fmt.Errorf("Multiple .an files found in %v", baseDir)
		}
		anthaFound = true
		cb.Logger.Log("transpiling", tet.ImportPath())
		if elemBs, err := ioutil.ReadFile(p); err != nil {
			return err
		} else if metaBs, err := ioutil.ReadFile(filepath.Join(filepath.Dir(p), "metadata.json")); err != nil && !os.IsNotExist(err) {
			return err
		} else if antha, err := tet.EnsureTranspiler(p, elemBs, metaBs); err != nil {
			return err
		} else {
			for _, ipt := range antha.ImportReqs {
				if err := tet.maybeRewriteImport(cb, wf, ipt); err != nil {
					return err
				}
			}
			return antha.Transform(anthaFiles)
		}
	})
	if err != nil {
		return err
	} else {
		return writeAnthaFiles(anthaFiles, filepath.Dir(baseDir))
	}
}

// Path should be the path in filepath format, to the antha element
// file. It does not have to be absolute. It is deliberately separate
// from bs because path might be some symbolic name unrelated to the
// actual source of the file (eg think some git repo). I.e. we can't
// just do an ioutil.ReadFile on path.
func (tet *TranspilableElementType) EnsureTranspiler(path string, elemBs, metaBs []byte) (*compile.Antha, error) {
	if tet.Transpiler == nil {
		fSet := token.NewFileSet()
		if src, err := parser.ParseFile(fSet, path, elemBs, parser.ParseComments); err != nil {
			return nil, err
		} else if antha, err := compile.NewAntha(fSet, src, metaBs); err != nil {
			return nil, err
		} else {
			tet.Transpiler = antha
		}
	}
	return tet.Transpiler, nil
}

func (tet *TranspilableElementType) maybeRewriteImport(cb *ComposerBase, wf *workflow.Workflow, ipt *compile.ImportReq) error {
	// we don't expect imports inside antha files to have revisions
	// within them (or specify repositories in any non-standard way). So, the strategy is:
	// 1. Look for longest matching repo and use that
	// 2. Otherwise (and this is most likely), it's not an import we should be rewriting.

	repoName, repo := wf.Repositories.LongestMatching(ipt.Path)
	if repo == nil {
		return nil // (2)
	}
	tet2 := NewTranspilableElementType(&workflow.ElementType{
		RepositoryName: repoName,
		ElementPath:    workflow.ElementPath(strings.TrimPrefix(ipt.Path, string(repoName))),
	})
	ipt.Path = tet2.ImportPath()
	cb.ensureElementType(tet2)

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
	dst, err := utils.CreateFile(outFile, utils.ReadWrite)
	if err != nil {
		return err
	}
	defer dst.Close() // nolint: errcheck

	src := file.NewReader()
	defer src.Close() // nolint: errcheck

	_, err = io.Copy(dst, src)
	return err
}
