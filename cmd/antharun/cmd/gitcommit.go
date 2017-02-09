// gitcommit
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func gopath() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(u.HomeDir, "go/src")
}

func GitCommit(path string) (commitID string, err error) {

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if path != "" {
		if found, _ := exists(path); !found {
			return "", fmt.Errorf("path %s not found locally", path)
		}
		err = os.Chdir(path)
		if err != nil {
			return "", err
		}
	}

	cmdName := "git"
	cmdArgs := []string{"rev-parse", "HEAD"}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return commitID, err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			commitID = scanner.Text()
		}
	}()

	err = cmd.Start()
	if err != nil {
		return commitID, err
	}

	err = cmd.Wait()
	if err != nil {
		return commitID, err
	}

	err = os.Chdir(pwd)
	if err != nil {
		return commitID, err
	}

	return commitID, nil
}

var reporttemplate string = ``
