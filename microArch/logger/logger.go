// microArch/logger/logger.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"strings"
)

//Logging Functions
func Track(message string, extra ...interface{}) {
	for _, h := range getMiddlewareList() {
		h.Log(TRACK, time.Now().Unix(), getSource(), message, extra...)
	}
}

func Info(message string, extra ...interface{}) {
	for _, h := range getMiddlewareList() {
		h.Log(INFO, time.Now().Unix(), getSource(), message, extra...)
	}
}

func Debug(message string, extra ...interface{}) {
	for _, h := range getMiddlewareList() {
		h.Log(DEBUG, time.Now().Unix(), getSource(), message, extra...)
	}
}

func Warning(message string, extra ...interface{}) {
	for _, h := range getMiddlewareList() {
		h.Log(WARNING, time.Now().Unix(), getSource(), message, extra...)
	}
}

func Error(message string, extra ...interface{}) {
	for _, h := range getMiddlewareList() {
		h.Log(ERROR, time.Now().Unix(), getSource(), message, extra...)
	}
}

func Fatal(message string, extra ...interface{}) {
	buf := make([]byte, 8192)
	runtime.Stack(buf, true)
	stack := make(map[string]string)
	stack["Stack"] = strings.TrimRight(string(buf), "\u0000")
	extra = append(extra, stack)
	for _, h := range getMiddlewareList() {
		h.Log(FATAL, time.Now().Unix(), getSource(), message, extra...)
	}

	pargs := []interface{}{"logger.Fatal", message}
	pargs = append(pargs, extra...)
	// Die if none of our middlewares dies first
	panic(fmt.Sprint(pargs...))
}

//telemetry and sensors functions
func Measure(message string, extra ...interface{}) {
	for _, h := range getMiddlewareList() {
		h.Measure(time.Now().Unix(), getSource(), message, extra...)
	}
}

func Sensor(message string, extra ...interface{}) {
	for _, h := range getMiddlewareList() {
		h.Sensor(time.Now().Unix(), getSource(), message, extra...)
	}
}

func Data(data interface{}, extra ...interface{}) {
	for _, h := range getMiddlewareList() {
		h.Data(time.Now().Unix(), data, append(extra, getSource()))
	}
}

var (
	middlewares []LoggerMiddleware
	_defaultmw  LoggerMiddleware //only used if none other available
)

// Register all middleware before making any logger calls
func RegisterMiddleware(m LoggerMiddleware) {
	middlewares = append(middlewares, m)
}

func getMiddlewareList() []LoggerMiddleware {
	//return the list or list with default inside if empty
	if len(middlewares) == 0 {
		if _defaultmw == nil {
			_defaultmw = &LogMiddleware{log.New(os.Stdout, "", log.LstdFlags)}
		}
		return []LoggerMiddleware{_defaultmw}
	}
	return middlewares
}

//LoggerMiddleware a means to react to specific log events
type LoggerMiddleware interface {
	Log(level LogLevel, ts int64, source, msg string, extra ...interface{})
	//Measure react to specific telemetry messages
	Measure(ts int64, source, msg string, extra ...interface{})
	//Sensor react to specific sensor readouts
	Sensor(ts int64, source, msg string, extra ...interface{})
	//Data saves data to database witha  timestamp and a set of extra fields in order to localise
	// it is an upper level log call that allows to dump unstructured data into our backend
	Data(ts int64, data interface{}, extra ...interface{})
}

//getSource returns a string representing the line of code that the preceeding call was generated.
// It makes 2 jumps on the call heap. getSource, and Log call are ignored
func getSource() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "NOSRC"
}
