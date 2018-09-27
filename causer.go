package errors

import (
	"fmt"
	"runtime"
	"strings"

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

func (e *causerWithCallers) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "%s\n", e.Error())
	for _, pc := range e.callers {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			fmt.Fprint(f, "unknown\n")
			continue
		}
		file, line := fn.FileLine(pc)
		pathSplit := strings.SplitN(file, "/src/", 2)
		file = pathSplit[len(pathSplit)-1]
		name := fn.Name()
		bits := strings.Split(name, ".")
		name = strings.Join(bits[1:], ".")
		fmt.Fprintf(f, "%s:%d %s\n", file, line, name)
	}
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
