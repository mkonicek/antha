package cmd

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// assumes GOPATH is in home directory if not set as environment variable
func gopath() string {
	// if gopath set return gopath
	if p := os.Getenv("GOPATH"); len(p) != 0 {
		return filepath.Join(p, "src")
	}
	// if not set assume under user's home directory
	u, err := user.Current()
	if err != nil {
		return ""
	}

	return filepath.Join(u.HomeDir, "go/src")
}

func gitCommit(path string) (string, error) {
	// nolint: gas
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path
	commitID, err := cmd.Output()
	return strings.TrimSpace(string(commitID)), err
}
