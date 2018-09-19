package errors

type Wrapper interface {
	Unwrap() error
}

type Causer interface {
	Cause() error
}

type WithCallers interface {
	Callers() []uintptr
}
