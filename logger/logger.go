package logger

import (
	"io"
	"os"

	stdlog "log"

	kitlog "github.com/go-kit/kit/log"
)

type Logger struct {
	// To understand this a bit better, consider that implementations
	// of Logger include the keyval pairs and a reference to the
	// underlying logger...
	kitlog.Logger

	// ...whereas SwapLogger is just the underlying logger
	swappable *kitlog.SwapLogger
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

// If len(ws) == 0 then os.Stderr is used. Otherwise, the logger logs
// out via ws only. Note that NewLogger should only be called once per
// process because it grabs the stdlog and redirects that via the
// new logger.
func NewLogger(ws ...io.Writer) *Logger {
	w := io.Writer(os.Stderr)
	if len(ws) != 0 {
		w = io.MultiWriter(ws...)
	}
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(w))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	sl := &kitlog.SwapLogger{}
	sl.Swap(logger)
	l := &Logger{
		Logger:    sl,
		swappable: sl,
	}
	stdlog.SetOutput(&wrapper{l: l})
	return l
}

func (l *Logger) With(keyvals ...interface{}) *Logger {
	logger := kitlog.With(l.Logger, keyvals...)
	return &Logger{
		Logger:    logger,
		swappable: l.swappable,
	}
}

// Replace the underlying writers of not only this logger, but the
// entire tree of loggers created by any calls to With from the root
// downwards.
func (l *Logger) SwapWriters(ws ...io.Writer) {
	w := io.Writer(os.Stderr)
	if len(ws) != 0 {
		w = io.MultiWriter(ws...)
	}
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(w))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	l.swappable.Swap(logger)
}

// Fatal uses the logger to log the error and then calls
// os.Exit(1). Deferred functions are not called, and it is quite
// unlikely you want this function: you certainly should never use
// this function from within an element. It is deliberately not a
// method on Logger to make it less likely to be accidentally used.
func Fatal(l *Logger, err error) {
	l.Log("fatal", err)
	os.Exit(1)
}
