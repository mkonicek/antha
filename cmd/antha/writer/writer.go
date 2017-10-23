package writer

import (
	"bytes"
	"fmt"
	"io"
)

// A Writer is a wrapper around an io.Writer to prepend a string to each line
type Writer struct {
	prepend []byte
	writer  io.Writer
	buf     bytes.Buffer
}

func (a *Writer) println(b []byte) (int, error) {
	n1, err := a.writer.Write(a.prepend)
	if err != nil {
		return n1, err
	}
	n2, err := a.writer.Write(b)
	if err != nil {
		return n1 + n2, err
	}
	n3, err := a.writer.Write([]byte{'\n'})
	if err != nil {
		return n1 + n2 + n3, err
	}
	return n1 + n2 + n3, nil
}

func (a *Writer) output(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	prevI := 0
	i := bytes.IndexByte(b[prevI:], '\n')
	for 0 <= i && prevI+i < len(b) {
		n, err := a.println(b[prevI : prevI+i])
		if err != nil {
			return prevI + n, err
		}

		prevI = prevI + i + 1
		if prevI >= len(b) {
			break
		}
		i = bytes.IndexByte(b[prevI:], '\n')
	}
	return prevI, nil
}

// Write implements an io.Writer
func (a *Writer) Write(p []byte) (n int, err error) {
	n1, e1 := a.buf.Write(p)
	w, e2 := a.output(a.buf.Bytes())

	if e1 != nil {
		a.buf.Next(w)
		return n1, e1
	}

	if e2 != nil {
		return len(p), e2
	}

	a.buf.Next(w)
	return len(p), nil
}

// Printf is a convenience function to write formatted output
func (a *Writer) Printf(format string, args ...interface{}) error {
	s := fmt.Sprintf(format, args...)
	_, err := a.Write([]byte(s))
	return err
}

// Flush outputs any unwritten input
func (a *Writer) Flush() error {
	written, err := a.output(a.buf.Bytes())
	if err != nil {
		return err
	}

	a.buf.Next(written)
	b := a.buf.Bytes()
	if len(b) > 0 {
		_, err := a.println(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// New creates a new Writer that wraps a given writer, prepending the given
// string to each line of output
func New(w io.Writer, prepend string) *Writer {
	return &Writer{writer: w, prepend: []byte(prepend)}
}
