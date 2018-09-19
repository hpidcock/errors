package errors

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
