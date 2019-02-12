package composer

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	return RunAndLogCommand(cmd, c.Logger.With("cmd", "generate").Log)
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

	if err := RunAndLogCommand(cmd, c.Logger.With("cmd", "build").Log); err != nil {
		return err
	} else {
		c.Logger.Log("compilation", "successful", "binary", outBin)
		return nil
	}
}

func (c *Composer) cleanOutDir() error {
	if c.Keep {
		c.Logger.Log("msg", "-keep set; not cleaning up", "path", c.OutDir)
		return nil
	}
	for _, leaf := range []string{"src", "workflow"} {
		if err := os.RemoveAll(filepath.Join(c.OutDir, leaf)); err != nil {
			return err
		}
	}
	return nil
}

func (c *Composer) RunWorkflow() error {
	if !c.Run {
		c.Logger.Log("msg", "running workflow disabled by flags")
		return nil
	}

	runOutDir, err := ioutil.TempDir(c.OutDir, fmt.Sprintf("antha-run-%s", c.Workflow.JobId))
	if err != nil {
		return err
	}
	c.Logger.Log("progress", "running compiled workflow", "outdir", runOutDir)
	outBin := filepath.Join(c.OutDir, "bin", string(c.Workflow.JobId))
	cmd := exec.Command(outBin, "-outdir", runOutDir)
	cmd.Env = []string{}

	// the workflow uses a proper logger these days so we don't need to do any wrapping
	logFunc := func(vals ...interface{}) error {
		// we are guaranteed len(vals) == 2, and that at [0] we have the key, which we ignore here
		fmt.Println(vals[1])
		return nil
	}
	return RunAndLogCommand(cmd, logFunc)
}

func RunAndLogCommand(cmd *exec.Cmd, logger func(...interface{}) error) error {
	if stdout, err := cmd.StdoutPipe(); err != nil {
		return err
	} else if stderr, err := cmd.StderrPipe(); err != nil {
		stdout.Close()
		return err
	} else {
		go drainToLogger(logger, stdout, "stdout")
		go drainToLogger(logger, stderr, "stderr")
		// lock to the current thread to ensure that thread state is
		// predictably inherited (eg namespaces etc) (see docs in
		// cmd.Run())
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		return cmd.Run()
	}
}

func drainToLogger(logger func(...interface{}) error, fh io.ReadCloser, key string) {
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		logger(key, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger("error", err.Error())
	}
}
