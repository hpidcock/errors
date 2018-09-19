package errors

import (
	"runtime"
)

type causerWithCallers struct {
	cause   error
	callers []uintptr
}

func (e *causerWithCallers) Error() string {
	return e.cause.Error()
}

func (e *causerWithCallers) Cause() error {
	return e.cause
}

func (e *causerWithCallers) Unwrap() error {
	return e.cause
}

func (e *causerWithCallers) Callers() []uintptr {
	return e.callers
}

func _ReturnTraceMarker() uintptr {
	var programCounter [1]uintptr
	runtime.Callers(1, programCounter[:])
	return programCounter[0]
}

var returnTraceMarker = _ReturnTraceMarker()

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	var programCounter [32]uintptr
	depth := runtime.Callers(2, programCounter[:])
	return &causerWithCallers{
		cause:   err,
		callers: programCounter[:depth],
	}
}

func Callers(err error) []uintptr {
	c, ok := err.(WithCallers)
	if ok == false {
		return nil
	}
	return c.Callers()
}

func Track(err error) error {
	if err == nil {
		return nil
	}

	c, ok := err.(*causerWithCallers)
	if ok == false {
		var programCounter [32]uintptr
		depth := runtime.Callers(2, programCounter[:])
		return &causerWithCallers{
			cause:   err,
			callers: programCounter[:depth],
		}
	}

	if c.callers[len(c.callers)-1] != returnTraceMarker {
		c.callers = append(c.callers, returnTraceMarker)
	}

	var programCounter [1]uintptr
	depth := runtime.Callers(2, programCounter[:])
	if depth == 1 {
		c.callers = append(c.callers, programCounter[0])
	}
	return c
}

func Is(err, target error) bool {
	for {
		if err == target {
			return true
		}
		switch w := err.(type) {
		case Wrapper:
			err = w.Unwrap()
			continue
		case Causer:
			err = w.Cause()
			continue
		}
		return false
	}
}

func Root(err error) error {
	switch w := err.(type) {
	case Wrapper:
		return Root(w.Unwrap())
	case Causer:
		return Root(w.Cause())
	}
	return err
}
