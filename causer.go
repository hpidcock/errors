package errors

import (
	"google.golang.org/grpc/status"
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

type grpcCauserWithCallers struct {
	causerWithCallers
}

func (e *grpcCauserWithCallers) GRPCStatus() *status.Status {
	g, ok := e.cause.(WithGRPCStatus)
	if ok == false {
		return nil
	}
	return g.GRPCStatus()
}
