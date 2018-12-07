package jlo_test

import (
	"fmt"
	"os"

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

	l.Infof("I'm real")
	// Output: {"lvl":"info","msg":"I'm real","time":"2018-08-02T21:48:56.856339554Z"}
}

// Create a new logger logging to stdout
func ExampleNewLogger_customFields() {
	l := jlo.NewLogger(os.Stdout)

	l.Infof("I'm real")
	// Output: {"@level":"info","@message":"I'm real","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}

func ExampleLogger_Debugf() {
	l := jlo.NewLogger(os.Stdout)

	l.SetLogLevel(jlo.DebugLevel)
	l.Debugf("I'm real")
	// Output: {"@level":"debug","@message":"I'm real","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}

func ExampleLogger_Infof() {
	l := jlo.NewLogger(os.Stdout)

	l.Infof("I'm real")
	// Output: {"@level":"info","@message":"I'm real","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}

func ExampleLogger_Warnf() {
	l := jlo.NewLogger(os.Stdout)

	l.Warnf("I'm real")
	// Output: {"@level":"warning","@message":"I'm real","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}

func ExampleLogger_Errorf() {
	l := jlo.NewLogger(os.Stdout)

	l.Errorf("I'm real")
	// Output: {"@level":"error","@message":"I'm real","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}

func ExampleLogger_Fatalf() {
	l := jlo.NewLogger(os.Stdout)

	l.Fatalf("I'm real")
	// Output: {"@level":"fatal","@message":"I'm real","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}

func ExampleLogger_SetLogLevel() {
	l := jlo.NewLogger(os.Stdout)

	l.SetLogLevel(jlo.DebugLevel)
	l.Debugf("I'm real")
	// Output: {"@level":"debug","@message":"I'm real","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}

func ExampleLogger_WithField() {
	l := jlo.NewLogger(os.Stdout)

	l = l.WithField("@request_id", "aa33ee55")
	l.Infof("I'm real")
	// Output: {"@level":"info","@message":"I'm real","@request_id":"aa33ee55","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}

func ExampleLogger_WithField_chaining() {
	l := jlo.NewLogger(os.Stdout)

	l.WithField("@request_id", "aa33ee55").Infof("I'm real")
	// Output: {"@level":"info","@message":"I'm real","@request_id":"aa33ee55","@timestamp":"2018-08-02T21:48:56.856339554Z"}
}
