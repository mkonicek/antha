package composer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (rs Repositories) LongestMatching(importPath string) (RepositoryPrefix, *Repository) {
	// Currently, because of the limitations of using GOPATH and not go
	// mod, and hence the limitations imposed by
	// Repositories.validate(), there is no chance that we ever need to
	// discriminate by length: we are guaranteed to find at most one
	// match. This however may become more relevant if/when we move to
	// go modules.
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

func (rs Repositories) CloneRepository(et *ElementType, dir string) error {
	if r, found := rs[et.RepositoryPrefix]; !found {
		return fmt.Errorf("Unknown RepositoryPrefix when FetchingFiles: '%v'", et.RepositoryPrefix)
	} else {
		return r.Clone(filepath.Join(dir, string(et.RepositoryPrefix)))
	}
}

func (r *Repository) Clone(dir string) error {
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	switch {
	case r.Branch == "" && r.Commit == "":
		return r.cloneFromDirectory(dir)
	case r.Commit != "":
		return r.cloneFromGitCommit(dir)
	default:
		return r.cloneFromGitBranch(dir)
	}
}

func (r *Repository) cloneFromDirectory(dir string) error {
	src := filepath.FromSlash(r.Directory)
	return filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if err != nil || !info.Mode().IsRegular() {
			return err
		}
		suffix := strings.TrimPrefix(p, src)
		dst := filepath.Join(dir, suffix)
		if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
			return err
		}
		srcFh, err := os.Open(p)
		if err != nil {
			return err
		}
		defer srcFh.Close()
		dstFh, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0400)
		if err != nil {
			return err
		}
		defer dstFh.Close()
		_, err = io.Copy(dstFh, srcFh)
		return err
	})
}

func (r *Repository) cloneFromGitCommit(dir string) error {
	if err := r.ensureGitRepo(); err != nil {
		return err
	}
	commitHash := plumbing.NewHash(r.Commit)
	if commit, err := r.gitRepo.CommitObject(commitHash); err != nil {
		return err
	} else if tree, err := r.gitRepo.TreeObject(commit.TreeHash); err != nil {
		return err
	} else {
		iter := tree.Files()
		defer iter.Close()
		for {
			if f, err := iter.Next(); err == io.EOF {
				break
			} else if err != nil {
				return err
			} else if c, err := f.Contents(); err != nil {
				return err
			} else {
				dst := filepath.Join(dir, filepath.FromSlash(f.Name))
				if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
					return err
				} else if err := ioutil.WriteFile(dst, []byte(c), 0400); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func (r *Repository) cloneFromGitBranch(dir string) error {
	if err := r.ensureGitRepo(); err != nil {
		return err
	}
	if branch, err := r.gitRepo.Branch(r.Branch); err != nil {
		return err
	} else if ch, err := r.gitRepo.ResolveRevision(plumbing.Revision(branch.Merge)); err != nil {
		return err
	} else {
		// we switch from branch to commit
		r.Commit = ch.String()
		r.Branch = ""
		return r.cloneFromGitCommit(dir)
	}
}

func (r *Repository) ensureGitRepo() error {
	if r.gitRepo == nil {
		if repo, err := git.PlainOpen(filepath.FromSlash(r.Directory)); err != nil {
			return err
		} else if _, err := repo.Head(); err != nil { // basically just check it works
			return err
		} else {
			r.gitRepo = repo
		}
	}
	return nil
}
