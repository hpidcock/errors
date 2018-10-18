package extgrpc

import (
	"google.golang.org/grpc/status"

	"github.com/hpidcock/errors"
)

type WithGRPCStatus interface {
	GRPCStatus() *status.Status
}

func init() {
	errors.RegisterMiddleware(middleware)
}

func middleware(err error) errors.Synthesizer {
	if _, ok := err.(WithGRPCStatus); ok {
		return synthesizer
	}
	return nil
}

func synthesizer(te errors.TrackedError) errors.CauserWithCallers {
	return &grpcTrackedError{
		TrackedError: te,
	}
}

type grpcTrackedError struct {
	errors.TrackedError
}

func (e *grpcTrackedError) GRPCStatus() *status.Status {
	g, ok := e.Cause().(WithGRPCStatus)
	if ok == false {
		return nil
	}
	return g.GRPCStatus()
}
