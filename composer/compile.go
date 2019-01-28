package composer

import (
	"bufio"
	"io"
	"os/exec"
	"path/filepath"
)

func (c *Composer) CompileWorkflow() error {
	c.Logger.Log("progress", "compiling workflow")
	if err := c.goGenerate(); err != nil {
		return err
	} else {
		return c.goBuild()
	}
}

func (c *Composer) goGenerate() error {
	cmd := exec.Command("go", "generate", "-x")
	cmd.Dir = filepath.Join(c.OutDir, "workflow")
	return runAndLogCommand(cmd, c.Logger.With("cmd", "generate", "cwd", cmd.Dir))
}

func (c *Composer) goBuild() error {
	cmd := exec.Command("go", "build")
	cmd.Dir = filepath.Join(c.OutDir, "workflow")
	if err := runAndLogCommand(cmd, c.Logger.With("cmd", "build", "cwd", cmd.Dir)); err != nil {
		return err
	} else {
		c.Logger.Log("compilation", "successful", "binary", filepath.Join(cmd.Dir, "workflow"))
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
		defer stderr.Close()
		defer stdout.Close()
		go drainToLogger(logger, stdout, "stdout")
		go drainToLogger(logger, stderr, "stderr")
		return cmd.Run()
	}
}

func drainToLogger(logger *Logger, fh io.ReadCloser, key string) {
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		logger.Log(key, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Log("error", err.Error())
	}
}
