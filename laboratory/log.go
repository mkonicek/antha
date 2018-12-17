package laboratory

import (
	"os"

	stdlog "log"

	kitlog "github.com/go-kit/kit/log"
)

type Log struct {
	kitlog.Logger
}

type wrapper struct {
	l *Log
}

func (w wrapper) Write(p []byte) (n int, err error) {
	if err := w.l.Log("msg", p); err != nil {
		return 0, err
	} else {
		return len(p), nil
	}
}

func NewLog() *Log {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	l := &Log{Logger: logger}

	stdlog.SetOutput(&wrapper{l: l})

	return l
}

func (l Log) Fatal(keyvals ...interface{}) {
	l.Log(keyvals...)
	os.Exit(1)
}

func (l Log) With(keyvals ...interface{}) *Log {
	logger := kitlog.With(l.Logger, keyvals...)
	return &Log{Logger: logger}
}
