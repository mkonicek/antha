package spawn

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/antha-lang/antha/cmd/antharun/writer"
	"github.com/mattn/go-colorable"
	"github.com/mgutz/ansi"
)

var (
	notStarted      = errors.New("not started")
	alreadyFinished = errors.New("already finished")
	timeoutError    = errors.New("timed out waiting for server to start")
)

type fn func() error

type Spawned struct {
	Command   *exec.Cmd
	done      chan bool
	closers   []fn
	firstLine copyUntil
}

// Expecting `Server listening at : %s`
func parse(bs []byte) (string, error) {
	s := string(bs)
	ss := strings.Split(s, " ")
	if len(ss) == 0 {
		return "", fmt.Errorf("cannot parse %q", s)
	}
	return ss[len(ss)-1], nil
}

func (a *Spawned) Uri() (string, error) {
	if a.Command.Process == nil {
		return "", notStarted
	}

	select {
	case <-time.After(5 * time.Second):
		return "", timeoutError
	case <-a.done:
		return "", alreadyFinished
	case <-a.firstLine.Done:
		return parse(a.firstLine.Bytes())
	}
}

func (a *Spawned) Close() (err error) {
	for _, c := range a.closers {
		if e := c(); e != nil {
			err = e
		}
	}

	return
}

func (a *Spawned) Start() error {
	if err := a.Command.Start(); err != nil {
		return err
	} else {
		go func() {
			defer close(a.done)
			a.Command.Wait()
		}()
		a.closers = append(a.closers, func() error {
			return a.Command.Process.Kill()
		})
	}
	return nil
}

// Copy output until stop byte is reached
type copyUntil struct {
	Stop byte      // Byte to stop copying at
	Done chan bool // If not nil, signal when done
	done bool
	buf  bytes.Buffer
	lock sync.Mutex
}

func (a *copyUntil) Bytes() []byte {
	a.lock.Lock()
	defer a.lock.Unlock()
	if !a.done {
		return nil
	}
	return a.buf.Bytes()
}

func (a *copyUntil) Write(bs []byte) (n int, err error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	n = len(bs)

	if a.done {
		return
	}

	for _, b := range bs {
		if b == a.Stop {
			a.done = true
			if a.Done != nil {
				close(a.Done)
			}
			break
		}
		a.buf.WriteByte(b)
	}
	return
}

func GoPackage(gopackage string, prefix string) (*Spawned, error) {
	runCmd := func(out io.Writer, prog string, args ...string) error {
		c := exec.Command(prog, args...)
		c.Stdin = os.Stdin
		if out == nil {
			c.Stdout = os.Stdout
		} else {
			c.Stdout = out
		}
		c.Stderr = os.Stderr
		return c.Run()
	}

	var buf bytes.Buffer
	if err := runCmd(nil, "go", "get", gopackage); err != nil {
		return nil, err
	} else if err := runCmd(&buf, "go", "list", "-f", "{{.Target}}", gopackage); err != nil {
		return nil, err
	} else if t := buf.String(); filepath.Base(filepath.Dir(t)) != "bin" {
		return nil, fmt.Errorf("package does not contain executable")
	}

	dir, err := ioutil.TempDir("", "spawn")
	if err != nil {
		return nil, err
	}
	s := &Spawned{
		done: make(chan bool),
		firstLine: copyUntil{
			Stop: '\n',
			Done: make(chan bool),
		},
	}
	s.closers = append(s.closers, func() error { return os.RemoveAll(dir) })

	spath := filepath.Join(dir, "server")
	if err := runCmd(nil, "go", "build", "-o", spath, gopackage); err != nil {
		return s, err
	}

	p := ansi.Color(prefix, "red:white") + " "
	w1 := writer.New(colorable.NewColorableStdout(), p)
	w2 := writer.New(colorable.NewColorableStderr(), p)

	cmd := exec.Command(spath, "-port", "0")
	cmd.Stdout = io.MultiWriter(w1, &s.firstLine)
	cmd.Stderr = w2

	s.Command = cmd
	s.closers = append(s.closers, func() error { return w1.Flush() })
	s.closers = append(s.closers, func() error { return w2.Flush() })
	return s, nil
}
