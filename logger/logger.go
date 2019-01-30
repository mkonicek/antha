package logger

import (
	"os"

	stdlog "log"

	kitlog "github.com/go-kit/kit/log"
)

type Logger struct {
	kitlog.Logger
}

type wrapper struct {
	l *Logger
}

func (w wrapper) Write(p []byte) (n int, err error) {
	if err := w.l.Log("msg", p); err != nil {
		return 0, err
	} else {
		return len(p), nil
	}
}

func NewLogger() *Logger {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	l := &Logger{Logger: logger}
	stdlog.SetOutput(&wrapper{l: l})
	return l
}

func (l Logger) With(keyvals ...interface{}) *Logger {
	logger := kitlog.With(l.Logger, keyvals...)
	return &Logger{Logger: logger}
}

func (l Logger) Fatal(err error) {
	l.Log("fatal", err.Error())
	os.Exit(1)
}
