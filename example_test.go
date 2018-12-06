package jlo_test

import (
	"fmt"
	"os"
	"time"

	"github.com/dcmn-com/jlo"
)

func ExampleLogLevel_String() {
	level := jlo.DebugLevel.String()
	fmt.Println(level)
	// Output: debug
}

// Create a new logger logging to stdout
func ExampleNewLogger() {
	l := jlo.NewLogger(os.Stdout)
	l.FieldKeyLevel = "lvl"
	l.FieldKeyMsg = "msg"
	l.FieldKeyTime = "time"

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.Infof("I'm real")
	// Output: {"lvl":"info","msg":"I'm real","time":"0001-01-01T00:00:00Z"}
}

// Create a new logger logging to stdout
func ExampleNewLogger_customFields() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.Infof("I'm real")
	// Output: {"@level":"info","@message":"I'm real","@timestamp":"0001-01-01T00:00:00Z"}
}

func ExampleLogger_Debugf() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.SetLogLevel(jlo.DebugLevel)
	l.Debugf("I'm real")
	// Output: {"@level":"debug","@message":"I'm real","@timestamp":"0001-01-01T00:00:00Z"}
}

func ExampleLogger_Infof() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.Infof("I'm real")
	// Output: {"@level":"info","@message":"I'm real","@timestamp":"0001-01-01T00:00:00Z"}
}

func ExampleLogger_Warnf() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.Warnf("I'm real")
	// Output: {"@level":"warning","@message":"I'm real","@timestamp":"0001-01-01T00:00:00Z"}
}

func ExampleLogger_Errorf() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.Errorf("I'm real")
	// Output: {"@level":"error","@message":"I'm real","@timestamp":"0001-01-01T00:00:00Z"}
}

func ExampleLogger_Fatalf() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.Fatalf("I'm real")
	// Output: {"@level":"fatal","@message":"I'm real","@timestamp":"0001-01-01T00:00:00Z"}
}

func ExampleLogger_SetLogLevel() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.SetLogLevel(jlo.DebugLevel)
	l.Debugf("I'm real")
	// Output: {"@level":"debug","@message":"I'm real","@timestamp":"0001-01-01T00:00:00Z"}
}

func ExampleLogger_WithField() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l = l.WithField("@request_id", "aa33ee55")
	l.Infof("I'm real")
	// Output: {"@level":"info","@message":"I'm real","@request_id":"aa33ee55","@timestamp":"0001-01-01T00:00:00Z"}
}

func ExampleLogger_WithField_chaining() {
	l := jlo.NewLogger(os.Stdout)

	// mocking time, ignore this line!
	l.Now = func() time.Time { return time.Time{} }

	l.WithField("@request_id", "aa33ee55").Infof("I'm real")
	// Output: {"@level":"info","@message":"I'm real","@request_id":"aa33ee55","@timestamp":"0001-01-01T00:00:00Z"}
}
