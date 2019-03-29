package workflow

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (rs Repositories) LongestMatching(importPath string) (RepositoryName, *Repository) {
	// Currently, because of the limitations of using GOPATH and not go
	// mod, and hence the limitations imposed by
	// Repositories.validate(), there is no chance that we ever need to
	// discriminate by length: we are guaranteed to find at most one
	// match. This however may become more relevant if/when we move to
	// go modules.
	//
	// This code says: if we have an importPath of
	// foo/bar/baz/Aliquot_Liquid and we have repositories with
	// prefixes foo/bar and foo/bar/baz then we choose to use the
	// longest repo prefix (foo/bar/baz).
	var winningRepo *Repository
	winningName := RepositoryName("")
	for repoName, repo := range rs {
		if strings.HasPrefix(importPath, string(repoName)) && len(repoName) > len(winningName) {
			winningRepo = repo
			winningName = repoName
		}
	}
	return winningName, winningRepo
}

func (rs Repositories) Clone(dir string) error {
	for prefix, repo := range rs {
		if err := repo.Clone(filepath.Join(dir, string(prefix))); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Clone(dir string) error {
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		// if the dir already exists, we assume we've already cloned it
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return r.Walk(copier(dir))
}

func copier(dir string) func(f *File) error {
	return func(f *File) error {
		if f == nil || !f.IsRegular {
			return nil
		}
		dst := filepath.Join(dir, f.Name)
		if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
			return err
		}
		srcFh, err := f.Contents()
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
	}
}

type TreeWalker func(*File) error

type File struct {
	Name      string // relative to the root of the walk, *always* in local filepath, never absolute
	IsRegular bool
	Contents  func() (io.ReadCloser, error)
}

func (r *Repository) Walk(fun TreeWalker) error {
	switch {
	case r.Branch == "" && r.Commit == "":
		return r.walkFromDirectory(fun)
	case r.Commit != "":
		return r.walkFromGitCommit(fun)
	default:
		return r.walkFromGitBranch(fun)
	}
}

const pathSepStr = string(os.PathSeparator) // os.PathSeparator is a char, which is less useful

func (r *Repository) walkFromDirectory(fun TreeWalker) error {
	src := filepath.Clean(filepath.FromSlash(r.Directory))
	var f File
	return filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		f.Name = strings.TrimPrefix(strings.TrimPrefix(p, src), pathSepStr)
		f.IsRegular = info.Mode().IsRegular()
		f.Contents = func() (io.ReadCloser, error) {
			return os.Open(p)
		}
		return fun(&f)
	})
}

func (r *Repository) walkFromGitCommit(fun TreeWalker) error {
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
		var f File
		for {
			if gf, err := iter.Next(); err == io.EOF {
				break
			} else if err != nil {
				return err
			} else {
				f.Name = strings.TrimPrefix(filepath.FromSlash(gf.Name), pathSepStr)
				f.IsRegular = gf.Mode.IsRegular()
				f.Contents = func() (io.ReadCloser, error) {
					if c, err := gf.Contents(); err != nil {
						return nil, err
					} else {
						return ioutil.NopCloser(bytes.NewBuffer([]byte(c))), nil
					}
				}
				if err := fun(&f); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func (r *Repository) walkFromGitBranch(fun TreeWalker) error {
	if err := r.ensureGitRepo(); err != nil {
		return err
	} else {
		// Sadly, a branch name is problematic: in a fully checked out
		// repo, the plain branch name can work. In a fresh full clone you
		// need to add a `refs/remotes/origin` prefix, and in a bare
		// clone, you need to add `refs/heads` prefix. WHY GIT? WHY?!!
		var ch *plumbing.Hash
		for _, prefix := range []string{"", "refs/remotes/origin/", "refs/heads/"} {
			if ch, err = r.gitRepo.ResolveRevision(plumbing.Revision(prefix + r.Branch)); err == nil {
				break
			}
		}
		if err != nil {
			return err
		} else {
			// we switch from branch to commit
			r.Commit = ch.String()
			r.Branch = ""
			return r.walkFromGitCommit(fun)
		}
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

func IsAnthaFile(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".an")
}
func IsAnthaMetadata(path string) bool {
	return filepath.Base(path) == "metadata.json"
}

type ElementTypeMap map[ElementTypeName]ElementType

func (r *Repository) FindAllElementTypes(repoName RepositoryName) (ElementTypeMap, error) {
	etm := make(ElementTypeMap)

	err := r.Walk(func(f *File) error {
		if !IsAnthaFile(f.Name) {
			return nil
		}

		dir := filepath.Dir(f.Name)
		ename := filepath.Base(dir)
		etm[ElementTypeName(ename)] = ElementType{
			ElementPath:    ElementPath(dir),
			RepositoryName: repoName,
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return etm, nil
	}
}

type ElementTypesByRepository map[*Repository]ElementTypeMap

func (rs Repositories) FindAllElementTypes() (ElementTypesByRepository, error) {
	types := make(ElementTypesByRepository)
	for repoName, rep := range rs {
		rmap, err := rep.FindAllElementTypes(repoName)
		if err != nil {
			return nil, err
		}
		types[rep] = rmap
	}

	return types, nil
}
