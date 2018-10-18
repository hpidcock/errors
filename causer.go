package errors

type TrackedError struct {
	cause   error
	callers []uintptr
}

func (e *TrackedError) Error() string {
	return e.cause.Error()
}

func (e *TrackedError) Cause() error {
	return e.cause
}

func (e *TrackedError) Unwrap() error {
	return e.cause
}

func (e *TrackedError) Callers() []uintptr {
	return e.callers
}
