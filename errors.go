package errors

import (
	"runtime"
)

type Wrapper interface {
	Unwrap() error
}

type Causer interface {
	Cause() error
}

type Callers interface {
	Callers() []uintptr
}

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
