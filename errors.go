package errors

import (
	"runtime"
)

func _BeginReturnTraceMarker() uintptr {
	var programCounter [1]uintptr
	runtime.Callers(1, programCounter[:])
	return programCounter[0]
}

func _EndReturnTraceMarker() uintptr {
	var programCounter [1]uintptr
	runtime.Callers(1, programCounter[:])
	return programCounter[0]
}

var beginReturnTraceMarker = _BeginReturnTraceMarker()
var endReturnTraceMarker = _EndReturnTraceMarker()

func createCauserWithCallers(err error) error {
	var programCounter [32]uintptr
	depth := runtime.Callers(3, programCounter[:])
	return synthesizer(err)(TrackedError{
		cause:   err,
		callers: programCounter[:depth],
	})
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return createCauserWithCallers(err)
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

	c, ok := err.(*TrackedError)
	if ok == false {
		return createCauserWithCallers(err)
	}

	if c.callers[len(c.callers)-1] != endReturnTraceMarker {
		c.callers = append(c.callers, beginReturnTraceMarker)
		c.callers = append(c.callers, endReturnTraceMarker)
	}

	runtime.Callers(2, c.callers[len(c.callers)-1:])
	c.callers = append(c.callers, endReturnTraceMarker)
	return err
}
