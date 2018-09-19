package errors

import (
	"google.golang.org/grpc/status"
)

type Wrapper interface {
	Unwrap() error
}

type Causer interface {
	Cause() error
}

type WithCallers interface {
	Callers() []uintptr
}

type WithGRPCStatus interface {
	GRPCStatus() *status.Status
}
