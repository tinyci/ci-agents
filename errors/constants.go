package errors

var errorMapping = map[string]*Error{
	"record not found": ErrNotFound,
}

// MapError finds an error by string and returns an appropriate *Error for it.
// The stack will NOT be preserved in the error and you will want to Wrap() it.
// If there is no potential mapping, a new *Error is returned.
func MapError(err interface{}) *Error {
	if err == nil {
		return nil
	}

	if e, ok := errorMapping[err.(error).Error()]; ok {
		return e
	}

	return New(err)
}

var (
	// ErrNotFound is a generic error for things that don't exist.
	ErrNotFound = New("not found")

	// ErrInvalidAuth is for authentication events that do not succeed.
	ErrInvalidAuth = New("invalid authentication")

	// ErrRunCanceled is thrown to github when the run is canceled.
	ErrRunCanceled = New("run canceled by user intervention")
)
