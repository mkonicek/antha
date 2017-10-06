package writer

import (
	"bytes"
	"fmt"
	"io"
)

type Writer struct {
	prepend []byte
	writer  io.Writer
	buf     bytes.Buffer
}

func (a *Writer) println(b []byte) error {
	if _, err := a.writer.Write(a.prepend); err != nil {
		return err
	} else if _, err := a.writer.Write(b); err != nil {
		return err
	} else if _, err := a.writer.Write([]byte{'\n'}); err != nil {
		return err
	}
	return nil
}

func (a *Writer) output(b []byte) int {
	if len(b) == 0 {
		return 0
	}
	prevI := 0
	i := bytes.IndexByte(b[prevI:], '\n')
	for 0 <= i && prevI+i < len(b) {
		a.println(b[prevI : prevI+i])
		prevI = prevI + i + 1
		if prevI >= len(b) {
			break
		}
		i = bytes.IndexByte(b[prevI:], '\n')
	}
	return prevI
}

func (a *Writer) Write(p []byte) (n int, err error) {
	a.buf.Write(p)
	written := a.output(a.buf.Bytes())
	a.buf.Next(written)
	return len(p), nil
}

func (a *Writer) Printf(format string, args ...interface{}) error {
	s := fmt.Sprintf(format, args...)
	_, err := a.Write([]byte(s))
	return err
}

func (a *Writer) Flush() error {
	written := a.output(a.buf.Bytes())
	a.buf.Next(written)
	b := a.buf.Bytes()
	if len(b) > 0 {
		a.println(b)
	}
	return nil
}

func New(w io.Writer, prepend string) *Writer {
	return &Writer{writer: w, prepend: []byte(prepend)}
}
