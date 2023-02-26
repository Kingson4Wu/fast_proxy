// Package logger
// copy from github.com/robfig/cron/v3@v3.0.1/logger.go
package logger

import (
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// DefaultLogger is used by Cron if none is specified.
var DefaultLogger = PrintfLogger(log.New(os.Stdout, "fastProxy: ", log.LstdFlags), os.Stdout)

// DiscardLogger can be used by callers to discard all log messages.
var DiscardLogger = PrintfLogger(log.New(io.Discard, "", 0), io.Discard)

var VerboseLogger = VerbosePrintfLogger(log.New(os.Stdout, "fastProxy: ", log.LstdFlags), os.Stdout)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger is the interface used in this package for logging, so that any backend
// can be plugged in. It is a subset of the github.com/go-logr/logr interface.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Debugf(template string, args ...interface{})

	Info(msg string, keysAndValues ...interface{})
	Infof(template string, args ...interface{})

	Warn(msg string, keysAndValues ...interface{})
	Warnf(template string, args ...interface{})

	Error(msg string, keysAndValues ...interface{})
	Errorf(template string, args ...interface{})

	GetWriter() io.Writer

	Flush()
}

// PrintfLogger wraps a Printf-based logger (such as the standard library "log")
// into an implementation of the Logger interface which logs errors only.
func PrintfLogger(l interface{ Printf(string, ...interface{}) }, w io.Writer) Logger {
	return printfLogger{l, ERROR, w}
}

// VerbosePrintfLogger wraps a Printf-based logger (such as the standard library
// "log") into an implementation of the Logger interface which logs everything.
func VerbosePrintfLogger(l interface{ Printf(string, ...interface{}) }, w io.Writer) Logger {
	return printfLogger{l, DEBUG, w}
}

type printfLogger struct {
	logger   interface{ Printf(string, ...interface{}) }
	logLevel Level
	w        io.Writer
}

func (pl printfLogger) Debug(msg string, keysAndValues ...interface{}) {
	if pl.logLevel <= DEBUG {
		keysAndValues = formatTimes(keysAndValues)
		pl.logger.Printf(
			formatString(len(keysAndValues)),
			append([]interface{}{"DEBUG", msg}, keysAndValues...)...)
	}
}

func (pl printfLogger) Debugf(template string, args ...interface{}) {
	if pl.logLevel <= DEBUG {
		pl.logger.Printf("DEBUG "+
			template, args...)
	}
}

func (pl printfLogger) Info(msg string, keysAndValues ...interface{}) {
	if pl.logLevel <= INFO {
		keysAndValues = formatTimes(keysAndValues)
		pl.logger.Printf(
			formatString(len(keysAndValues)),
			append([]interface{}{"INFO", msg}, keysAndValues...)...)
	}
}

func (pl printfLogger) Infof(template string, args ...interface{}) {
	if pl.logLevel <= INFO {
		pl.logger.Printf("INFO "+
			template, args...)
	}
}

func (pl printfLogger) Warn(msg string, keysAndValues ...interface{}) {
	if pl.logLevel <= WARN {
		keysAndValues = formatTimes(keysAndValues)
		pl.logger.Printf(
			formatString(len(keysAndValues)),
			append([]interface{}{"WARN", msg}, keysAndValues...)...)
	}
}

func (pl printfLogger) Warnf(template string, args ...interface{}) {
	if pl.logLevel <= WARN {
		pl.logger.Printf("WARN "+
			template, args...)
	}
}

func (pl printfLogger) Error(msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	if pl.logLevel <= ERROR {
		pl.logger.Printf(
			formatString(len(keysAndValues)),
			append([]interface{}{"ERROR", msg}, keysAndValues...)...)
	}
}

func (pl printfLogger) Errorf(template string, args ...interface{}) {
	if pl.logLevel <= ERROR {
		pl.logger.Printf("ERROR "+
			template, args...)
	}
}

func (pl printfLogger) GetWriter() io.Writer {
	return pl.w
}

func (pl printfLogger) Flush() {
}

// formatString returns a logfmt-like format string for the number of
// key/values.
func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("%v=%v")
	}
	return sb.String()
}

// formatTimes formats any time.Time values as RFC3339.
func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}
