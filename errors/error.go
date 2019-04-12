package errors

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

// Error encapsulates a series of errors.
type Error struct {
	Errs   []string        `json:"errors"`
	Frames []xerrors.Frame `json:"-"`
	Log    bool            `json:"log"`
}

// New constructs a new error
func New(res interface{}) *Error {
	if res == nil {
		return nil
	}

	if _, ok := res.(*Error); ok {
		return res.(*Error)
	}

	switch res := res.(type) {
	case error:
		return doinit(res.Error(), true)
	case string:
		return doinit(res, true)
	}

	return nil
}

// Errorf is the formatted version of New().
func Errorf(f string, args ...interface{}) *Error {
	return New(fmt.Sprintf(f, args...))
}

func doinit(str string, log bool) *Error {
	return &Error{
		Errs:   []string{str},
		Frames: []xerrors.Frame{xerrors.Caller(2)},
		Log:    log,
	}
}

// NewNoLog returns an error which will not be logged.
func NewNoLog(str string) *Error {
	return doinit(str, false)
}

// SetLog sets the state for whether or not logging is permitted.
func (e *Error) SetLog(log bool) {
	e.Log = log
}

// GetLog returns the state of the log argument; for usage by error reporting utilities.
func (e *Error) GetLog() bool {
	return e.Log
}

// Copy makes a copy of the error: for great justice
func (e *Error) Copy() *Error {
	e2 := *e
	e2.Errs = []string{}
	e2.Frames = []xerrors.Frame{}
	e2.Errs = append(e2.Errs, e.Errs...)
	e2.Frames = append(e2.Frames, e.Frames...)

	return &e2
}

func (e *Error) Error() string {
	if e == nil {
		return "<invalid nil error!>"
	}

	return strings.Join(e.Errs, ": ")
}

func (e *Error) String() string {
	return e.Error()
}

// Format formats the errors implementing fmt.Formatter
func (e *Error) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
}

// FormatError implements xerrors.Formatter and allows it to provide extended detail.
func (e *Error) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())

	if p.Detail() {
		for i := len(e.Frames) - 1; i >= 0; i-- {
			p.Print(e.Errs[i], ": ")
			e.Frames[i].Format(p)
		}
	}
	return nil
}

// Wrap wraps an error string. If the original error is nil, no wrapping occurs and nil is returned.
func (e *Error) Wrap(err interface{}) *Error {
	if e == nil {
		return nil
	}

	e2 := e.Copy()

	switch err := err.(type) {
	case *Error:
		if e.Contains(err) {
			return e
		}

		e2.Errs = append(err.Errs, e2.Errs...)
		e2.Frames = append(err.Frames, e2.Frames...)
		return e2
	case error, string:
		if e.Contains(err) {
			return e
		}
		e2.Errs = append(New(err).Errs, e2.Errs...)
	default:
		panic("invalid value passed to errors.Wrap: must be string or error")
	}

	e2.Frames = append(e2.Frames, xerrors.Caller(1))
	return e2
}

// Wrapf wraps an error with formatting
func (e *Error) Wrapf(str string, args ...interface{}) *Error {
	e2 := e.Copy()
	wrapped := New(fmt.Sprintf(str, args...))
	e2.Errs = append(wrapped.Errs, e2.Errs...)
	e2.Frames = append(wrapped.Frames, e2.Frames...)
	return e2
}

// Contains checks if e contains err. Strings and errors are supported; otherwise it will return false.
func (e *Error) Contains(err interface{}) bool {
	for _, e2 := range e.Errs {
		if s, ok := err.(string); ok {
			if e2 == s {
				return true
			}
		} else if _, ok := err.(error); !ok {
			return false
		} else {
			if strings.Contains(e2, err.(error).Error()) {
				return true
			}
		}
	}

	return false
}
