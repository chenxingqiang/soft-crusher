// File: pkg/errors/errors.go

package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorCode int

const (
	ErrUnknown ErrorCode = iota
	ErrNotFound
	ErrInvalidInput
	ErrDatabaseError
	// Add more error codes as needed
)

type Error struct {
	Err        error
	Message    string
	StackTrace string
	Code       ErrorCode
}

func NewWithCode(message string, code ErrorCode) *Error {
	return &Error{
		Message:    message,
		StackTrace: getStackTrace(),
		Code:       code,
	}
}

func (e *Error) GetCode() ErrorCode {
	return e.Code
}

// New creates a new Error
func New(message string) *Error {
	return &Error{
		Message:    message,
		StackTrace: getStackTrace(),
	}
}

// Newf creates a new Error with formatted message
func Newf(format string, args ...interface{}) *Error {
	return &Error{
		Message:    fmt.Sprintf(format, args...),
		StackTrace: getStackTrace(),
	}
}

// Wrap wraps an existing error with additional message
func Wrap(err error, message string) *Error {
	return &Error{
		Err:        err,
		Message:    message,
		StackTrace: getStackTrace(),
	}
}

// Error returns the error message
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.Err
}

// getStackTrace returns the stack trace as a string
func getStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			fmt.Fprintf(&builder, "%s:%d\n", frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return builder.String()
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return fmt.Sprintf("%v", err) == fmt.Sprintf("%v", target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
func As(err error, target interface{}) bool {
	switch t := target.(type) {
	case **Error:
		if e, ok := err.(*Error); ok {
			*t = e
			return true
		}
	}
	return false
}

// GetStackTrace returns the stack trace
func (e *Error) GetStackTrace() string {
	return e.StackTrace
}
