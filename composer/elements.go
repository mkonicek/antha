package composer

import (
	"fmt"
	"io"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type ElementSource struct {
	Prefix string
	Path   string
	repo   *git.Repository
}

type ElementSources []ElementSource

// We sort by reverse length of prefix. That way when we do a prefix
// match, we will naturally get the longest prefix first.
func (es ElementSources) Sort() {
	sort.Slice(es, func(i, j int) bool {
		return len(es[i].Prefix) > len(es[j].Prefix)
	})
}

// The schema we expect here is commit-sha-or-branch-name/remaining-path.
func (es ElementSources) Match(element string) (*ElementSource, string, string, error) {
	for idx, e := range es {
		if strings.HasPrefix(element, e.Prefix) {
			tail := strings.TrimPrefix(element, e.Prefix)
			tail = strings.TrimPrefix(tail, "/")
			split := strings.SplitN(tail, "/", 2)
			if len(split) != 2 {
				return nil, "", "", fmt.Errorf("Invalid path: '%s'. Required at least one / character to separated commit from path", tail)
			}

			return &es[idx], split[0], split[1], nil
		}
	}
	return nil, "", "", nil
}

func (e ElementSource) FetchFiles(revision, path string) (map[string]string, error) {
	if err := e.ensureRepo(); err != nil {
		return nil, err
	}

	var commitHash plumbing.Hash
	if branch, err := e.repo.Branch(revision); err == git.ErrBranchNotFound {
		// it's not a branch, so assume it's a commit hash
		commitHash = plumbing.NewHash(revision)
	} else if err != nil {
		return nil, err
	} else if ch, err := e.repo.ResolveRevision(plumbing.Revision(branch.Merge)); err != nil {
		return nil, err
	} else {
		commitHash = *ch
	}

	// now follow that commitHash
	if commit, err := e.repo.CommitObject(commitHash); err != nil {
		return nil, err
	} else if tree, err := e.repo.TreeObject(commit.TreeHash); err != nil {
		return nil, err
	} else {
		results := make(map[string]string)
		iter := tree.Files()
		for {
			if f, err := iter.Next(); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			} else if !strings.HasPrefix(f.Name, path) {
				continue
			} else if c, err := f.Contents(); err != nil {
				return nil, err
			} else {
				results[f.Name] = c
			}
		}
		return results, nil
	}
}

func (e *ElementSource) ensureRepo() error {
	if e.repo == nil {
		if repo, err := git.PlainOpen(e.Path); err != nil {
			return err
		} else if _, err := repo.Head(); err != nil {
			return err
		} else {
			e.repo = repo
		}
	}
	return nil
}
