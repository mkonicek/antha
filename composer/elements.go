package composer

import (
	"fmt"
	"io"
	"path"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type ElementSource struct {
	// The git prefix (normally something like a hostname, maybe with a path following)
	Prefix string
	// The directory on the local machine where we can find this repo
	Path string

	repo *git.Repository
}

type ElementSources []ElementSource

// We sort by reverse length of prefix. That way when we do a prefix
// match, we will naturally get the longest prefix first.
func (es ElementSources) Sort() {
	sort.Slice(es, func(i, j int) bool {
		return len(es[i].Prefix) > len(es[j].Prefix)
	})
}

// The element here is as would appear in the workflow. I.e. full git
// repo prefix, a commit or branch name, and then the remaining path
// to a directory. The schema we expect here is
// repo-prefix/commit-sha-or-branch-name/remaining-path. It is not an
// error if no matching ElementSource is found, merely the returned
// LocatedElement will be nil. Error only if the element is
// unparsable.
func (es ElementSources) Match(element string) (*LocatedElement, error) {
	for idx, e := range es {
		if strings.HasPrefix(element, e.Prefix) {
			tail := strings.TrimPrefix(element, e.Prefix)
			tail = strings.TrimPrefix(tail, "/")
			split := strings.SplitN(tail, "/", 2)
			if len(split) != 2 {
				return nil, fmt.Errorf("Invalid path: '%s'. Required at least one / character to separated commit from path", tail)
			}

			return NewLocatedElement(&es[idx], split[0], split[1]), nil
		}
	}
	return nil, nil
}

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

func (le LocatedElement) FetchFiles() error {
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
		return nil
	}
}

func (e *ElementSource) ensureRepo() error {
	if e.repo == nil {
		if repo, err := git.PlainOpen(e.Path); err != nil {
			return err
		} else if _, err := repo.Head(); err != nil { // basically just check it works
			return err
		} else {
			e.repo = repo
		}
	}
	return nil
}
