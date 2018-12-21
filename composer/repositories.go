package composer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type RepoId string

type Repository struct {
	id           RepoId
	ImportPrefix string `json:"ImportPrefix"`
	File         string `json:"File"`
	Git          string `json:"Git"`

	repo repo
}

type repo interface {
	FetchFiles(*LocatedElement) (map[string][]byte, error)
}

type Repositories map[RepoId]*Repository

func (rs *Repositories) UnmarshalJSON(bs []byte) error {
	if string(bs) == "null" {
		return nil
	}
	rs2 := make(map[RepoId]*Repository)
	if err := json.Unmarshal(bs, &rs2); err != nil {
		return err
	}
	for id, repo := range rs2 {
		repo.id = RepoId(id)
		if fileEmpty, gitEmpty := repo.File == "", repo.Git == ""; fileEmpty == gitEmpty {
			return fmt.Errorf("Exactly one of File and Git must be specified in repository '%s'", id)
		} else if fileEmpty {
			repo.repo = newGitRepo(repo.Git)
		} else {
			repo.repo = newFileRepo(repo.File)
		}
	}
	*rs = Repositories(rs2)
	return nil
}

func (r *Repository) ImportPath(element *ElementSource) string {
	if r.File != "" {
		return path.Join(r.ImportPrefix, element.Path)
	} else {
		return path.Join(r.ImportPrefix, element.CommitOrBranch(), element.Path)
	}
}

type gitRepo struct {
	dir  string
	repo *git.Repository
}

func newGitRepo(dir string) *gitRepo {
	return &gitRepo{
		dir: dir,
	}
}

func (g *gitRepo) ensureRepo() error {
	if g.repo == nil {
		if repo, err := git.PlainOpen(g.dir); err != nil {
			return err
		} else if _, err := repo.Head(); err != nil { // basically just check it works
			return err
		} else {
			g.repo = repo
		}
	}
	return nil
}

func (g *gitRepo) FetchFiles(le *LocatedElement) (map[string][]byte, error) {
	if bEmpty, cEmpty := le.Element.Branch == "", le.Element.Commit == ""; bEmpty == cEmpty {
		return nil, fmt.Errorf("Exactly one of Commit and Branch must be specified for git (Commit: '%s', Branch: '%s')", le.Element.Commit, le.Element.Branch)
	}

	if err := g.ensureRepo(); err != nil {
		return nil, err
	}

	var commitHash plumbing.Hash
	if le.Element.Commit != "" {
		commitHash = plumbing.NewHash(le.Element.Commit)
	} else {
		if branch, err := g.repo.Branch(le.Element.Branch); err != nil {
			// it's not a branch, so assume it's a commit hash
			return nil, err
		} else if ch, err := g.repo.ResolveRevision(plumbing.Revision(branch.Merge)); err != nil {
			return nil, err
		} else {
			commitHash = *ch
		}
	}

	// now follow that commitHash
	if commit, err := g.repo.CommitObject(commitHash); err != nil {
		return nil, err
	} else if tree, err := g.repo.TreeObject(commit.TreeHash); err != nil {
		return nil, err
	} else {
		p := strings.Trim(le.Element.Path, "/") + "/"
		results := make(map[string][]byte)
		iter := tree.Files()
		for {
			if f, err := iter.Next(); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			} else if !strings.HasPrefix(f.Name, p) {
				continue
			} else if c, err := f.Contents(); err != nil {
				return nil, err
			} else {
				results[strings.TrimPrefix(f.Name, p)] = []byte(c)
			}
		}
		return results, nil
	}
}

type fileRepo struct {
	dir string
}

func newFileRepo(dir string) *fileRepo {
	return &fileRepo{
		dir: dir,
	}
}

func (f *fileRepo) FetchFiles(le *LocatedElement) (map[string][]byte, error) {
	if le.Element.Branch != "" || le.Element.Commit != "" {
		return nil, fmt.Errorf("Neither Commit nor Branch can be specified for file (Commit: '%s', Branch: '%s')", le.Element.Commit, le.Element.Branch)
	}

	p := path.Join(le.Repository.File, strings.Trim(le.Element.Path, "/")) + "/"
	results := make(map[string][]byte)
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if !info.Mode().IsRegular() {
			return nil
		} else if c, err := ioutil.ReadFile(path); err != nil {
			return err
		} else {
			results[strings.TrimPrefix(path, p)] = c
			return nil
		}
	})
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}
