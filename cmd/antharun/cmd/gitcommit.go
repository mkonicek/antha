// gitcommit
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

type SSHCommander struct {
	User string
	IP   string
}

func (s *SSHCommander) Command(cmd ...string) *exec.Cmd {
	arg := append(
		[]string{
		//	fmt.Sprintf("%s@%s", s.User, s.IP),
		},
		cmd...,
	)
	return exec.Command("", arg...) //exec.Command("ssh", arg...)
}

func GitCommit() (commitID string, err error) {

	/*
		//commander := SSHCommander{"root", "50.112.213.24"}
		commander := SSHCommander{"", ""}

		cmd := []string{
			"git",
			"rev-parse",
			"HEAD",
		}

		// am I doing this automation thing right?
		if cmderr := commander.Command(cmd...); cmderr != nil {
			fmt.Fprintln(os.Stderr, "There was an error running SSH command: ", cmderr)
			os.Exit(1)

			_, err = cmderr.Output()
			return err
		}
		return err
	*/

	cmdName := "git"
	cmdArgs := []string{"rev-parse", "HEAD"}
	//cmdName := "cd"
	//cmdArgs := []string{"$GOPATH/src/github.com/antha-lang/antha"}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		//os.Exit(1)
		return commitID, err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			//	fmt.Printf("git commit | %s\n", scanner.Text())
			commitID = scanner.Text()
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		//os.Exit(1)
		return commitID, err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		//os.Exit(1)
		return commitID, err
	}
	return commitID, nil
}
