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
	"time"

	"github.com/antha-lang/antha/workflow"
)

func (mc *mainComposer) compileWorkflow() error {
	mc.Logger.Log("progress", "compiling workflow")
	genStart := time.Now()
	if err := mc.goGenerate(); err != nil {
		return err
	}
	buildStart := time.Now()
	if err := mc.goBuild(); err != nil {
		return err
	}
	mc.Logger.Log("go_generate", buildStart.Sub(genStart), "go_build", time.Now().Sub(buildStart))
	return mc.cleanOutDir()
}

func (cb *ComposerBase) goGenerate() error {
	cmd := exec.Command("go", "generate", "-x")
	cmd.Dir = filepath.Join(cb.OutDir, "workflow")

	return RunAndLogCommand(cmd, cb.Logger.With("cmd", "generate").Log)
}

func (mc *mainComposer) goBuild() error {
	outBin := filepath.Join(mc.OutDir, "bin", "workflow")
	cmd := exec.Command("go", "build", "-o", outBin)
	if mc.LinkedDrivers {
		cmd.Args = append(cmd.Args, "-tags", "linkedDrivers protobuf")
	}
	cmd.Dir = filepath.Join(mc.OutDir, "workflow")
	cmd.Env = SetEnvGoPath(os.Environ(), mc.OutDir)

	if err := RunAndLogCommand(cmd, mc.Logger.With("cmd", "build").Log); err != nil {
		return err
	} else {
		mc.Logger.Log("compilation", "successful", "binary", outBin)
		return nil
	}
}

func (tc *testComposer) goTest() error {
	cmd := exec.Command("go", "test", "-v")
	if tc.LinkedDrivers {
		cmd.Args = append(cmd.Args, "-tags", "linkedDrivers")
	}
	cmd.Dir = filepath.Join(tc.OutDir, "workflow")
	cmd.Env = SetEnvGoPath(os.Environ(), tc.OutDir)

	if err := RunAndLogCommand(cmd, tc.Logger.With("cmd", "test").Log); err != nil {
		return err
	} else {
		tc.Logger.Log("testing", "successful")
		return nil
	}
}

// NB: this may well need revising when we move to go mod
func SetEnvGoPath(env []string, outDir string) []string {
	for idx, s := range env {
		if len(s) >= 7 && "GOPATH=" == s[:7] {
			env[idx] = fmt.Sprintf("GOPATH=%s:%s", outDir, s[7:])
			break
		}
	}
	return env
}

func (mc *mainComposer) cleanOutDir() error {
	if mc.Keep {
		mc.Logger.Log("msg", "-keep set; not cleaning up", "path", mc.OutDir)
		return nil
	}
	for _, leaf := range []string{"src", "workflow"} {
		if err := os.RemoveAll(filepath.Join(mc.OutDir, leaf)); err != nil {
			return err
		}
	}
	return nil
}

func (mc *mainComposer) runWorkflow() error {
	if !mc.Run {
		mc.Logger.Log("msg", "running workflow disabled by flags")
		return nil
	}

	runOutDir, err := ioutil.TempDir(mc.OutDir, "antha-run-outputs")
	if err != nil {
		return err
	}
	mc.Logger.Log("progress", "running compiled workflow", "outdir", runOutDir, "indir", mc.InDir)
	outBin := filepath.Join(mc.OutDir, "bin", "workflow")
	cmd := exec.Command(outBin, "-outdir", runOutDir, "-indir", mc.InDir)
	cmd.Env = []string{}

	// the workflow uses a proper logger these days so we don't need to do any wrapping
	logFunc := func(vals ...interface{}) error {
		// we are guaranteed len(vals) == 2, and that at [0] we have the key, which we ignore here
		fmt.Println(vals[1])
		return nil
	}
	return RunAndLogCommand(cmd, logFunc)
}

func (cb *ComposerBase) prepareDrivers(cfg *workflow.Config) error {
	// Here, if we're meant to compile something, we attempt that, on
	// the basis that when we come to run the workflow itself, we may
	// not have the necessary sources around or build environment.  If
	// we're just meant to be running some command, we make sure we
	// take a copy of that binary into a sensible place within the
	// outdir, again so that we should be able to guarantee it exists
	// when we come to workflow execution.
	conns := make(map[workflow.DeviceInstanceID]*workflow.ParsedConnection)

	for id, cfg := range cfg.GilsonPipetMax.Devices {
		conns[id] = &cfg.ParsedConnection
	}
	for id, cfg := range cfg.Tecan.Devices {
		conns[id] = &cfg.ParsedConnection
	}
	for id, cfg := range cfg.CyBio.Devices {
		conns[id] = &cfg.ParsedConnection
	}
	for id, cfg := range cfg.Labcyte.Devices {
		conns[id] = &cfg.ParsedConnection
	}
	for id, cfg := range cfg.Hamilton.Devices {
		conns[id] = &cfg.ParsedConnection
	}

	for id, cfg := range conns {
		outBin := filepath.Join(cb.OutDir, "bin", "drivers", string(id))
		if err := os.MkdirAll(filepath.Dir(outBin), 0700); err != nil {
			return err

		} else if cfg.CompileAndRun != "" {
			cb.Logger.Log("instructionPlugin", string(id), "building", cfg.CompileAndRun)
			cmd := exec.Command("go", "build", "-o", outBin, cfg.CompileAndRun)
			cmd.Dir = filepath.Dir(outBin)
			if err := RunAndLogCommand(cmd, cb.Logger.With("cmd", "build", "instructionPlugin", string(id)).Log); err != nil {
				return err
			}
			cfg.ExecFile = outBin
			cfg.CompileAndRun = ""

		} else if cfg.ExecFile != "" {
			src, err := os.Open(cfg.ExecFile)
			if err != nil {
				return err
			}
			defer src.Close()
			dst, err := os.OpenFile(outBin, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
			if err != nil {
				return err
			}
			defer dst.Close()
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}
			cfg.ExecFile = outBin
		}
	}
	return nil
}

// Only starts the command, does not Wait for it.
func StartAndLogCommand(cmd *exec.Cmd, logger func(...interface{}) error) error {
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
		return cmd.Start()
	}
}

// Starts and Waits for the command (unless error)
func RunAndLogCommand(cmd *exec.Cmd, logger func(...interface{}) error) error {
	if err := StartAndLogCommand(cmd, logger); err != nil {
		return err
	} else {
		return cmd.Wait()
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
