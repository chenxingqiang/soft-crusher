package logging

import (
	"io/ioutil"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	// Test debug mode
	InitLogger(true, "")
	if logger == nil {
		t.Error("Logger should not be nil")
	}

	// Test production mode
	InitLogger(false, "")
	if logger == nil {
		t.Error("Logger should not be nil")
	}

	// Test with log file
	tmpfile, err := ioutil.TempFile("", "test_log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	InitLogger(false, tmpfile.Name())
	if logger == nil {
		t.Error("Logger should not be nil")
	}

	// Test logging
	Info("Test info message")
	Error("Test error message")
	Debug("Test debug message")
	Warn("Test warning message")

	// Check if the log file contains the messages
	content, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) == "" {
		t.Error("Log file should not be empty")
	}
}

func TestSetLogger(t *testing.T) {
	customLogger, _ := zap.NewProduction()
	SetLogger(customLogger)
	if logger != customLogger {
		t.Error("SetLogger did not set the custom logger")
	}
}

func TestWith(t *testing.T) {
	InitLogger(true, "")
	childLogger := With(zap.String("key", "value"))
	if childLogger == nil {
		t.Error("Child logger should not be nil")
	}
}

func TestGetLogger(t *testing.T) {
	InitLogger(true, "")
	if GetLogger() != logger {
		t.Error("GetLogger did not return the correct logger")
	}
}
