package composer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func (c *Composer) CompileWorkflow() error {
	c.Logger.Log("progress", "compiling workflow")
	if err := c.goGenerate(); err != nil {
		return err
	} else if err := c.goBuild(); err != nil {
		return err
	} else {
		return c.cleanOutDir()
	}
}

func (c *Composer) goGenerate() error {
	cmd := exec.Command("go", "generate", "-x")
	cmd.Dir = filepath.Join(c.OutDir, "workflow")
	if path, _ := os.LookupEnv("PATH"); path != "" {
		// because we don't know what other binaries go generate will
		// end up calling (well, we do - go-bindata, but we don't know
		// where that is)
		cmd.Env = []string{"PATH=" + path}
	} else {
		cmd.Env = []string{}
	}

	return runAndLogCommand(cmd, c.Logger.With("cmd", "generate", "cwd", cmd.Dir))
}

func (c *Composer) goBuild() error {
	outBin := filepath.Join(c.OutDir, "bin", string(c.Workflow.JobId))
	cmd := exec.Command("go", "build", "-o", outBin)
	cmd.Dir = filepath.Join(c.OutDir, "workflow")
	gopath := c.OutDir
	if cur, _ := os.LookupEnv("GOPATH"); cur != "" {
		gopath = fmt.Sprintf("%s:%s", gopath, cur)
	}
	cmd.Env = []string{"GOPATH=" + gopath}

	if err := runAndLogCommand(cmd, c.Logger.With("cmd", "build", "cwd", cmd.Dir)); err != nil {
		return err
	} else {
		c.Logger.Log("compilation", "successful", "binary", outBin)
		return nil
	}
}

func runAndLogCommand(cmd *exec.Cmd, logger *Logger) error {
	if stdout, err := cmd.StdoutPipe(); err != nil {
		return err
	} else if stderr, err := cmd.StderrPipe(); err != nil {
		stdout.Close()
		return err
	} else {
		go drainToLogger(logger, stdout, "stdout")
		go drainToLogger(logger, stderr, "stderr")
		return cmd.Run()
	}
}

func drainToLogger(logger *Logger, fh io.ReadCloser, key string) {
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		logger.Log(key, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Log("error", err.Error())
	}
}

func (c *Composer) cleanOutDir() error {
	for _, leaf := range []string{"src", "workflow"} {
		if err := os.RemoveAll(filepath.Join(c.OutDir, leaf)); err != nil {
			return err
		}
	}
	return nil
}
