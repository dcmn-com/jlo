// Light-weight json logging
package jlo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel represents a log level used by Logger type
type LogLevel int

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

const (
	// DebugLevel is the most verbose output and logs messages on all levels
	DebugLevel LogLevel = iota
	// InfoLevel logs messages on all levels except the DebugLevel
	InfoLevel
	// WarningLevel logs messages on all levels except DebugLevel and InfoLevel
	WarningLevel
	// ErrorLevel logs messages on ErrorLevel and FatalLevel
	ErrorLevel
	// FatalLevel logs messages on FatalLevel
	FatalLevel

	// FieldKeyLevel is the log level log field name
	FieldKeyLevel = "@level"
	// FieldKeyMsg is the log message log field name
	FieldKeyMsg = "@message"
	// FieldKeyTime is the time log field name
	FieldKeyTime = "@timestamp"

	// default log level used for initialization of global logger and all additional ones
	defaultLogLevel = InfoLevel
	// RFC3339 time format "2006-01-02T15:04:05.999999999Z07:00"
	timeFormat = time.RFC3339Nano
)

// Logger logs json formatted messages to a certain output destination
type Logger struct {
	FieldKeyMsg   string
	FieldKeyLevel string
	FieldKeyTime  string
	fields        map[string]string
	mu            sync.RWMutex
	logLevel      LogLevel
	now           func() time.Time
	outMu         sync.Mutex
	out           io.Writer
}

// DefaultLogger is a default logger logging to stdout
var DefaultLogger = NewLogger(os.Stdout)

// NewLogger creates a new logger which will write to the passed in io.Writer
func NewLogger(out io.Writer) *Logger {
	return &Logger{
		FieldKeyMsg:   FieldKeyMsg,
		FieldKeyLevel: FieldKeyLevel,
		FieldKeyTime:  FieldKeyTime,
		fields:        make(map[string]string),
		logLevel:      defaultLogLevel,
		now:           func() time.Time { return time.Now().UTC() },
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
func (l *Logger) WithField(key, value string) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	fields := map[string]string{
		key: value,
	}

	for k, v := range l.fields {
		fields[k] = v
	}

	clone := NewLogger(l.out)
	clone.now = l.now
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

	l.out.Write(entry)
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

	data := map[string]string{
		l.FieldKeyTime:  l.now().Format(timeFormat),
		l.FieldKeyLevel: level.String(),
		l.FieldKeyMsg:   msg,
	}

	for k, v := range l.fields {
		data[k] = v
	}

	// Error is ignored intentionally, as no errors are expected because the
	// data type to be marshaled will never change.
	entry, _ := json.Marshal(data)
	entry = append(entry, '\n')
	return entry
}
