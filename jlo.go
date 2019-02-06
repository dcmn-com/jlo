// Light-weight json logging
package jlo

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	easyjson "github.com/mailru/easyjson"
)

//go:generate easyjson -no_std_marshalers $GOFILE

const (
	// FieldKeyLevel is the log level log field name
	FieldKeyLevel = "@level"
	// FieldKeyMsg is the log message log field name
	FieldKeyMsg = "@message"
	// FieldKeyTime is the time log field name
	FieldKeyTime = "@timestamp"
	// FieldKeyCommit is the commit log field name
	FieldKeyCommit = "@commit"
)

// LogLevel represents a log level used by Logger type
type LogLevel int

func (lvl LogLevel) MarshalEasyJSON() ([]byte, error) {
	s := lvl.String()
	out := make([]byte, len(s)+2)

	copy(out[1:len(out)-2], s[:])
	out[0], out[len(out)-1] = '"', '"'

	return out, nil
}

const (
	// UnknownLevel means the log level could not be parsed
	UnknownLevel LogLevel = iota
	// DebugLevel is the most verbose output and logs messages on all levels
	DebugLevel
	// InfoLevel logs messages on all levels except the DebugLevel
	InfoLevel
	// WarningLevel logs messages on all levels except DebugLevel and InfoLevel
	WarningLevel
	// ErrorLevel logs messages on ErrorLevel and FatalLevel
	ErrorLevel
	// FatalLevel logs messages on FatalLevel
	FatalLevel
)

// String returns a string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarningLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	default:
		return "fatal"
	}
}

// logLevel is used as initial value upon creation of a new logger
var logLevel = InfoLevel

// SetLogLevel changes the global log level
func SetLogLevel(level LogLevel) {
	logLevel = level
}

var Now = func() time.Time {
	return time.Now().UTC()
}

//easyjson:json
type Entry map[string]interface{}

// Logger logs json formatted messages to a certain output destination
type Logger struct {
	FieldKeyMsg   string
	FieldKeyLevel string
	FieldKeyTime  string
	fields        Entry
	mu            sync.RWMutex
	logLevel      LogLevel
	outMu         sync.Mutex
	out           io.Writer
}

// DefaultLogger returns a new default logger logging to stdout
func DefaultLogger() *Logger {
	return NewLogger(os.Stdout)
}

// NewLogger creates a new logger which will write to the passed in io.Writer
func NewLogger(out io.Writer) *Logger {
	return &Logger{
		FieldKeyMsg:   FieldKeyMsg,
		FieldKeyLevel: FieldKeyLevel,
		FieldKeyTime:  FieldKeyTime,
		fields:        make(Entry),
		logLevel:      logLevel,
		out:           out,
	}
}

// Fatalf logs a message on FatalLevel
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	l.log(FatalLevel, format, args...)
}

// Errorf logs a messages on ErrorLevel
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.logLevel <= ErrorLevel {
		l.log(ErrorLevel, format, args...)
	}
}

// Warnf logs a messages on WarningLevel
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.logLevel <= WarningLevel {
		l.log(WarningLevel, format, args...)
	}
}

// Infof logs a messages on InfoLevel
func (l *Logger) Infof(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.logLevel <= InfoLevel {
		l.log(InfoLevel, format, args...)
	}
}

// Debugf logs a messages on DebugLevel
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.logLevel <= DebugLevel {
		l.log(DebugLevel, format, args...)
	}
}

// SetLogLevel changes the log level
func (l *Logger) SetLogLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logLevel = level
}

// WithField returns a copy of the logger with a custom field set, which will be
// included in all subsequent logs
func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	fields := make(Entry, len(l.fields)+1)
	fields[key] = value

	for k, v := range l.fields {
		fields[k] = v
	}

	clone := NewLogger(l.out)
	clone.fields = fields
	return clone
}

// log builds the final log data by concatenating the log template array data with
// the values for log level, timestamp and log message
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	entry := l.generateLogEntry(level, format, args...)

	// wrap Write() method call in mutex to guarantee atomic writes
	l.outMu.Lock()
	defer l.outMu.Unlock()

	l.out.Write(append(entry, '\n'))
}

// generateLogEntry generates a log entry by gathering all field data and marshal
// everthing to json format
func (l *Logger) generateLogEntry(level LogLevel, format string, args ...interface{}) []byte {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}

	data := make(Entry, len(l.fields)+3)

	data[l.FieldKeyTime] = Now()
	data[l.FieldKeyLevel] = level.String()
	data[l.FieldKeyMsg] = msg

	for k, v := range l.fields {
		data[k] = v
	}

	// Error is ignored intentionally, as no errors are expected because the
	// data type to be marshaled will never change.
	entry, _ := easyjson.Marshal(data)
	return entry
}
