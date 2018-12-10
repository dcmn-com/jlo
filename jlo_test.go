package jlo_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dcmn-com/jlo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testTime = "2018-08-02T21:48:56.856339554Z"

type LogFunc func(string, ...interface{})

func TestMain(m *testing.M) {
	jlo.Now = func() time.Time {
		t, err := time.Parse(time.RFC3339Nano, testTime)
		if err != nil {
			panic(err)
		}
		return t
	}
	os.Exit(m.Run())
}

func Test_SetLogLevel(t *testing.T) {
	jlo.SetLogLevel(jlo.InfoLevel)
	buf := bytes.NewBuffer(nil)
	infoLogger := jlo.NewLogger(buf)

	infoLogger.Debugf("should not log")
	assert.Empty(t, buf.String())

	jlo.SetLogLevel(jlo.DebugLevel)

	debugLogger := jlo.NewLogger(buf)
	debugLogger.Debugf("should log")

	var msg map[string]interface{}
	err := json.NewDecoder(buf).Decode(&msg)
	require.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"@level":     "debug",
		"@message":   "should log",
		"@timestamp": msg["@timestamp"],
	}, msg)
}

func Test_ParseLogLevel(t *testing.T) {
	tests := map[string]jlo.LogLevel{
		"fatal":   jlo.FatalLevel,
		"FATAL":   jlo.FatalLevel,
		"error":   jlo.ErrorLevel,
		"ERROR":   jlo.ErrorLevel,
		"warn":    jlo.WarningLevel,
		"WARN":    jlo.WarningLevel,
		"warning": jlo.WarningLevel,
		"WARNING": jlo.WarningLevel,
		"info":    jlo.InfoLevel,
		"INFO":    jlo.InfoLevel,
		"debug":   jlo.DebugLevel,
		"DEBUG":   jlo.DebugLevel,
	}

	for str, level := range tests {
		parsed, err := jlo.ParseLogLevel(str)
		require.NoError(t, err)
		assert.Equal(t, level, parsed)
	}
}

func Test_ParseLogLevel_UnknownLevel(t *testing.T) {
	parsed, err := jlo.ParseLogLevel("unknown")
	assert.Error(t, err)
	assert.Equal(t, jlo.UnknownLevel, parsed)
}
func Test_Logger_Debugf(t *testing.T) {

	tests := map[string]struct {
		String  string
		Args    []interface{}
		Message string
	}{
		"simple": {
			String:  "I'm real",
			Message: "I'm real",
		},
		"with format args": {
			String:  "string: %s int: %d float: %.2f bool: %t",
			Args:    []interface{}{"I'm real", 5, 0.1, true},
			Message: "string: I'm real int: 5 float: 0.10 bool: true",
		},
		"message including double brackets": {
			String:  `I'm "real"`,
			Message: `I'm \"real\"`,
		},
		"message arg including double brackets": {
			String:  "I'm %s",
			Args:    []interface{}{`"real"`},
			Message: `I'm \"real\"`,
		},
		"message args including double brackets": {
			String:  "string: %s int: %d float: %.2f bool: %t",
			Args:    []interface{}{`"I'm real"`, 5, 0.1, true},
			Message: `string: \"I'm real\" int: 5 float: 0.10 bool: true`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			l := jlo.NewLogger(buf)
			l.SetLogLevel(jlo.DebugLevel)

			l.Debugf(test.String, test.Args...)

			assert.JSONEq(t, fmt.Sprintf(`{
				"@message":   "%s",
				"@level":     "debug",
				"@timestamp": "%s"
			}`, test.Message, testTime), buf.String())
		})
	}
}

func Test_Logger_Infof(t *testing.T) {

	tests := map[string]struct {
		String  string
		Args    []interface{}
		Message string
	}{
		"simple": {
			String:  "I'm real",
			Message: "I'm real",
		},
		"with format args": {
			String:  "string: %s int: %d float: %.2f bool: %t",
			Args:    []interface{}{"I'm real", 5, 0.1, true},
			Message: "string: I'm real int: 5 float: 0.10 bool: true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			l := jlo.NewLogger(buf)
			l.SetLogLevel(jlo.InfoLevel)

			l.Infof(test.String, test.Args...)

			assert.JSONEq(t, fmt.Sprintf(`{
				"@message":   "%s",
				"@level":     "info",
				"@timestamp": "%s"
			}`, test.Message, testTime), buf.String())
		})
	}
}

func Test_Logger_Warnf(t *testing.T) {

	tests := map[string]struct {
		String  string
		Args    []interface{}
		Message string
	}{
		"simple": {
			String:  "I'm real",
			Message: "I'm real",
		},
		"with format args": {
			String:  "string: %s int: %d float: %.2f bool: %t",
			Args:    []interface{}{"I'm real", 5, 0.1, true},
			Message: "string: I'm real int: 5 float: 0.10 bool: true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			l := jlo.NewLogger(buf)
			l.SetLogLevel(jlo.WarningLevel)

			l.Warnf(test.String, test.Args...)

			assert.JSONEq(t, fmt.Sprintf(`{
				"@message":   "%s",
				"@level":     "warning",
				"@timestamp": "%s"
			}`, test.Message, testTime), buf.String())
		})
	}
}

func Test_Logger_Errorf(t *testing.T) {

	tests := map[string]struct {
		String  string
		Args    []interface{}
		Message string
	}{
		"simple": {
			String:  "I'm real",
			Message: "I'm real",
		},
		"with format args": {
			String:  "string: %s int: %d float: %.2f bool: %t",
			Args:    []interface{}{"I'm real", 5, 0.1, true},
			Message: "string: I'm real int: 5 float: 0.10 bool: true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			l := jlo.NewLogger(buf)
			l.SetLogLevel(jlo.InfoLevel)

			l.Errorf(test.String, test.Args...)

			assert.JSONEq(t, fmt.Sprintf(`{
				"@message":   "%s",
				"@level":     "error",
				"@timestamp": "%s"
			}`, test.Message, testTime), buf.String())
		})
	}
}

func Test_Logger_Fatalf(t *testing.T) {

	tests := map[string]struct {
		String  string
		Args    []interface{}
		Message string
	}{
		"simple": {
			String:  "I'm real",
			Message: "I'm real",
		},
		"with format args": {
			String:  "string: %s int: %d float: %.2f bool: %t",
			Args:    []interface{}{"I'm real", 5, 0.1, true},
			Message: "string: I'm real int: 5 float: 0.10 bool: true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			l := jlo.NewLogger(buf)
			l.SetLogLevel(jlo.FatalLevel)

			l.Fatalf(test.String, test.Args...)

			assert.JSONEq(t, fmt.Sprintf(`{
				"@message":   "%s",
				"@level":     "fatal",
				"@timestamp": "%s"
			}`, test.Message, testTime), buf.String())
		})
	}
}

func Test_Logger_Infof_EnsureNewlineDelimitedJSON(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := jlo.NewLogger(buf)
	l.Infof("I'm real")
	assert.True(t, strings.HasSuffix(buf.String(), "}\n"))
}

func Test_Logger_Infof_SpecialCharacterUse(t *testing.T) {

	tests := map[string]struct {
		String  string
		Args    []interface{}
		Message string
	}{
		"double quotes": {
			String:  `"I'm" %s`,
			Args:    []interface{}{`"real"`},
			Message: `"I'm" "real"`,
		},
		"backspace": {
			String:  "I'm\b %s",
			Args:    []interface{}{"real\b"},
			Message: "I'm\b real\b",
		},
		"form feed": {
			String:  "I'm\f %s",
			Args:    []interface{}{"real\f"},
			Message: "I'm\f real\f",
		},
		"new line": {
			String:  "I'm\n %s",
			Args:    []interface{}{"real\n"},
			Message: "I'm\n real\n",
		},
		"carriage return": {
			String:  "I'm\r %s",
			Args:    []interface{}{"real\r"},
			Message: "I'm\r real\r",
		},
		"tab": {
			String:  "I'm\t %s",
			Args:    []interface{}{"real\t"},
			Message: "I'm\t real\t",
		},
		"backslash": {
			String:  "I'm\\ %s",
			Args:    []interface{}{"real\\"},
			Message: "I'm\\ real\\",
		},
		"all": {
			String:  `"` + " \b \f \n \r \t \\ %s",
			Args:    []interface{}{`"` + " \b \f \n \r \t \\"},
			Message: `"` + " \b \f \n \r \t \\ " + `"` + " \b \f \n \r \t \\",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			l := jlo.NewLogger(buf)

			l.Infof(test.String, test.Args...)

			var decoded map[string]string
			err := json.Unmarshal(buf.Bytes(), &decoded)
			if err != nil {
				t.Logf("log destination data: '%s'", buf.String())
				t.Fatal("error unmarshaling log destination data")
			}

			assert.Equal(t, map[string]string{
				"@message":   test.Message,
				"@level":     "info",
				"@timestamp": testTime,
			}, decoded)
		})
	}
}

func Test_Logger_WithField(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := jlo.NewLogger(buf)
	l.WithField("@request_id", "e44c2a9").
		WithField("@version", 2.1).
		WithField("@revision", 6).
		Infof("I'm real")

	// check that field is in logs
	assert.JSONEq(t, fmt.Sprintf(`{
		"@level": "info",
		"@message": "I'm real",
		"@timestamp": "%s",
		"@request_id": "e44c2a9",
		"@revision": 6,
		"@version": 2.1
	}`, testTime), buf.String())

	// check that original logger is unaffected
	buf.Reset()
	l.Infof("I'm real")
	assert.JSONEq(t, fmt.Sprintf(`{
		"@message":    "I'm real",
		"@level":      "info",
		"@timestamp":  "%s"
	}`, testTime), buf.String())
}

func Test_Logger_WithField_EnsureNoChangeWhenIgnoringReturnValue(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := jlo.NewLogger(buf)
	l.WithField("WithField", "will only be set in returned logger")
	l.Infof("I'm real")

	assert.JSONEq(t, fmt.Sprintf(`{
		"@message":   "I'm real",
		"@level":     "info",
		"@timestamp": "%s"
	}`, testTime), buf.String())
}

func Test_Logger_SetLogLevel(t *testing.T) {

	tests := map[string]struct {
		GetLogFunc  func(l *jlo.Logger) LogFunc
		LogLevel    jlo.LogLevel
		SetLogLevel jlo.LogLevel
		Message     string
	}{
		"Debug log with DebugLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Debugf
			},
			LogLevel:    jlo.DebugLevel,
			SetLogLevel: jlo.DebugLevel,
			Message:     "I'm real",
		},
		"Debug log with InfoLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Debugf
			},
			LogLevel:    jlo.DebugLevel,
			SetLogLevel: jlo.InfoLevel,
			Message:     "",
		},
		"Debug log with WarningLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Debugf
			},
			LogLevel:    jlo.DebugLevel,
			SetLogLevel: jlo.WarningLevel,
			Message:     "",
		},
		"Debug log with ErrorLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Debugf
			},
			LogLevel:    jlo.DebugLevel,
			SetLogLevel: jlo.ErrorLevel,
			Message:     "",
		},
		"Info log with DebugLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Infof
			},
			LogLevel:    jlo.InfoLevel,
			SetLogLevel: jlo.DebugLevel,
			Message:     "I'm real",
		},
		"Info log with InfoLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Infof
			},
			LogLevel:    jlo.InfoLevel,
			SetLogLevel: jlo.InfoLevel,
			Message:     "I'm real",
		},
		"Info log with WarningLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Infof
			},
			SetLogLevel: jlo.WarningLevel,
			Message:     "",
		},
		"Info log with ErrorLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Infof
			},
			LogLevel:    jlo.InfoLevel,
			SetLogLevel: jlo.ErrorLevel,
			Message:     "",
		},
		"Warning log with DebugLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Warnf
			},
			LogLevel:    jlo.WarningLevel,
			SetLogLevel: jlo.DebugLevel,
			Message:     "I'm real",
		},
		"Warning log with InfoLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Warnf
			},
			LogLevel:    jlo.WarningLevel,
			SetLogLevel: jlo.InfoLevel,
			Message:     "I'm real",
		},
		"Warning log with WarningLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Warnf
			},
			LogLevel:    jlo.WarningLevel,
			SetLogLevel: jlo.WarningLevel,
			Message:     "I'm real",
		},
		"Warning log with ErrorLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Warnf
			},
			LogLevel:    jlo.WarningLevel,
			SetLogLevel: jlo.ErrorLevel,
			Message:     "",
		},
		"Error log with DebugLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Errorf
			},
			LogLevel:    jlo.ErrorLevel,
			SetLogLevel: jlo.DebugLevel,
			Message:     "I'm real",
		},
		"Error log with InfoLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Errorf
			},
			LogLevel:    jlo.ErrorLevel,
			SetLogLevel: jlo.InfoLevel,
			Message:     "I'm real",
		},
		"Error log with WarningLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Errorf
			},
			LogLevel:    jlo.ErrorLevel,
			SetLogLevel: jlo.WarningLevel,
			Message:     "I'm real",
		},
		"Error log with ErrorLevel": {
			GetLogFunc: func(l *jlo.Logger) LogFunc {
				return l.Errorf
			},
			LogLevel:    jlo.ErrorLevel,
			SetLogLevel: jlo.ErrorLevel,
			Message:     "I'm real",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			l := jlo.NewLogger(buf)
			l.SetLogLevel(test.SetLogLevel)

			test.GetLogFunc(l)("I'm real")

			if test.Message == "" {
				if buf.Bytes() != nil {
					t.Error("Expected log destination data to be empty")
				}
				return
			}

			assert.JSONEq(t, fmt.Sprintf(`{
				"@message":   "%s",
				"@level":     "%s",
				"@timestamp": "%s"
			}`, test.Message, test.LogLevel.String(), testTime), buf.String())
		})
	}
}

const (
	testStringShort = "this is a short log string"
	testStringLong  = `Lorem ipsum dolor sit amet, consectetuer adipiscing elit.
		Proin in tellus sit amet nibh dignissim sagittis. In convallis. Etiam ligula
		pede, sagittis quis, interdum ultricies, scelerisque eu. Donec iaculis gravida
		nulla. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum
		dolore eu fugiat nulla pariatur. Nunc tincidunt ante vitae massa. Ut enim ad
		minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam,
		nisi ut aliquid ex ea commodi consequatur? Vestibulum erat nulla, ullamcorper
		nec, rutrum non, nonummy ac, erat. Curabitur bibendum justo non orci. Etiam
		sapien elit, consequat eget, tristique non, venenatis quis, ante. Duis pulvinar.
		Aliquam erat volutpat. In convallis. Vestibulum fermentum tortor id mi.`
)

func Benchmark_Logger_ShortString(b *testing.B) {
	benchmarkLogger(b, testStringShort)
}

func Benchmark_Logger_ShortString_WithArgs(b *testing.B) {
	benchmarkLogger(b, testStringShort+"%s", "I'm real")
}

func Benchmark_Logger_ShortString_WithField(b *testing.B) {
	benchmarkLoggerWithField(b, testStringShort)
}

func Benchmark_Logger_ShortString_WithFieldAndArgs(b *testing.B) {
	benchmarkLoggerWithField(b, testStringShort+"%s", "I'm real")
}

func Benchmark_Logger_VeryLongString(b *testing.B) {
	benchmarkLogger(b, testStringLong)
}

func Benchmark_Logger_VeryLongString_WithArgs(b *testing.B) {
	benchmarkLogger(b, testStringLong+"%s", "I'm real")
}

func Benchmark_Logger_VeryLongString_WithFields(b *testing.B) {
	benchmarkLoggerWithField(b, testStringLong)
}

func Benchmark_Logger_VeryLongString_WithFieldsAndArgs(b *testing.B) {
	benchmarkLoggerWithField(b, testStringLong+"%s", "I'm real")
}

func benchmarkLogger(b *testing.B, format string, args ...interface{}) {
	l := jlo.NewLogger(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		l.Infof(format, args...)
	}
}

func benchmarkLoggerWithField(b *testing.B, format string, args ...interface{}) {
	l := jlo.NewLogger(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		l.WithField("I'm", "real").Infof(format, args...)
	}
}
