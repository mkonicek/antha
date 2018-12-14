package composer

import (
	"fmt"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
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
