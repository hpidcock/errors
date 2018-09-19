package errors

import (
	"fmt"
	"runtime"
	"testing"
)

func TestTrackError(t *testing.T) {
	err := func() error {
		return Track(func() error {
			return Track(func() error {
				return Track(fmt.Errorf("hello"))
			}())
		}())
	}()
	t.Logf("%v", err)
	c := Callers(err)
	for _, v := range c {
		f := runtime.FuncForPC(v)
		file, line := f.FileLine(v)
		t.Logf("%s:%d %s", file, line, f.Name())
	}
}
