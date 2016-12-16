// gitcommit
package cmd

import (
	"bufio"
	"os/exec"
)

func GitCommit() (commitID string, err error) {
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
	return commitID, nil
}
