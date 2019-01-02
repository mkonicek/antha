package composer

import (
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

func (rs Repositories) LongestMatching(importPath string) (RepositoryPrefix, *Repository) {
	var winningRepo *Repository
	winningPrefix := RepositoryPrefix("")
	for repoPrefix, repo := range rs {
		if strings.HasPrefix(importPath, string(repoPrefix)) && len(repoPrefix) > len(winningPrefix) {
			winningRepo = repo
			winningPrefix = repoPrefix
		}
	}
	return winningPrefix, winningRepo
}

func (rs Repositories) FetchFiles(et *ElementType) (map[string][]byte, error) {
	if et.files != nil {
		return et.files, nil
	} else if r, found := rs[et.RepositoryPrefix]; !found {
		return nil, fmt.Errorf("Unknown RepositoryPrefix when FetchingFiles: '%v'", et.RepositoryPrefix)
	} else {
		return r.FetchFiles(et)
	}
}

func (r *Repository) FetchFiles(et *ElementType) (map[string][]byte, error) {
	var err error
	var files map[string][]byte
	switch {
	case r.Branch == "" && r.Commit == "":
		files, err = r.fetchFilesFromDirectory(et)
	case r.Commit != "":
		files, err = r.fetchFilesFromGitCommit(et)
	default:
		files, err = r.fetchFilesFromGitBranch(et)
	}
	if err != nil {
		return nil, err
	} else {
		et.files = files
		return files, nil
	}
}

func (r *Repository) fetchFilesFromDirectory(et *ElementType) (map[string][]byte, error) {
	d := filepath.FromSlash(path.Join(r.Directory, string(et.ElementPath)))
	results := make(map[string][]byte)
	err := filepath.Walk(d, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if !info.Mode().IsRegular() {
			return nil
		} else if c, err := ioutil.ReadFile(p); err != nil {
			return err
		} else {
			results[filepath.ToSlash(strings.TrimPrefix(p, d))] = c
			return nil
		}
	})
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (r *Repository) fetchFilesFromGitCommit(et *ElementType) (map[string][]byte, error) {
	if err := r.ensureGitRepo(); err != nil {
		return nil, err
	}
	commitHash := plumbing.NewHash(r.Commit)
	if commit, err := r.gitRepo.CommitObject(commitHash); err != nil {
		return nil, err
	} else if tree, err := r.gitRepo.TreeObject(commit.TreeHash); err != nil {
		return nil, err
	} else {
		prefix := strings.Trim(string(et.ElementPath), "/") + "/"
		results := make(map[string][]byte)
		iter := tree.Files()
		for {
			if f, err := iter.Next(); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			} else if !strings.HasPrefix(f.Name, prefix) {
				continue
			} else if c, err := f.Contents(); err != nil {
				return nil, err
			} else {
				results[strings.TrimPrefix(f.Name, prefix)] = []byte(c)
			}
		}
		return results, nil
	}
}

func (r *Repository) fetchFilesFromGitBranch(et *ElementType) (map[string][]byte, error) {
	if err := r.ensureGitRepo(); err != nil {
		return nil, err
	}
	if branch, err := r.gitRepo.Branch(r.Branch); err != nil {
		return nil, err
	} else if ch, err := r.gitRepo.ResolveRevision(plumbing.Revision(branch.Merge)); err != nil {
		return nil, err
	} else {
		r.Commit = ch.String()
		return r.fetchFilesFromGitCommit(et)
	}
}

func (r *Repository) ensureGitRepo() error {
	if r.gitRepo == nil {
		if repo, err := git.PlainOpen(r.Directory); err != nil {
			return err
		} else if _, err := repo.Head(); err != nil { // basically just check it works
			return err
		} else {
			r.gitRepo = repo
		}
	}
	return nil
}
