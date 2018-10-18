package errors

import (
	"fmt"
	"runtime"
	"strings"
)

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

func (e *TrackedError) Format(f fmt.State, c rune) {
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
