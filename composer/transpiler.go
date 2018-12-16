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
	Files map[string][]byte
}

func NewLocatedElement(source *ElementSource, commit, remainingPath string) *LocatedElement {
	remainingPath = strings.Trim(remainingPath, "/") + "/"
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
		results := make(map[string][]byte)
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
				results[strings.TrimPrefix(f.Name, le.Path)] = []byte(c)
			}
		}
		le.Files = results
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
	// 3. otherwise attempt to parse it as if it has a revision / branch
	// 4. Otherwise (and this is most likely), it's just not an import we should be rewriting.

	if strings.HasPrefix(ipt.Path, le.Source.Prefix) {
		remainingPath := strings.TrimPrefix(ipt.Path, le.Source.Prefix)
		class := path.Join(le.Source.Prefix, le.Commit, remainingPath)
		if le2, found := c.classes[class]; found { // (1)
			ipt.Path = le2.ImportPath
			return nil
		} else {
			le2 := NewLocatedElement(le.Source, le.Commit, remainingPath)
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

	if le2, err := c.Config.ElementSources.Match(ipt.Path); err != nil || le2 == nil {
		return nil // not an error in this case (format err), so we don't return err
	} else if err := le2.FetchFiles(); err != nil {
		// most likely that the branch or commit doesn't exist (mis-parsing). Also not an error.
		return nil
	} else if len(le2.Files) > 0 {
		ipt.Path = le2.ImportPath
		c.EnsureLocatedElement(le2)
		return nil
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
