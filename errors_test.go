package errors

import (
	"fmt"
	"testing"
)

func TestTrackError(t *testing.T) {
	err := func() error {
		return func() error {
			return func() error {
				return Track(fmt.Errorf("hello"))
			}()
		}()
	}()
	t.Logf("%v", err)
	t.Logf("%v", Callers(err))
}
