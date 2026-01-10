package error

import "errors"

var (
	ErrSubmitFailed               = errors.New("Submit error")
	ErrExecuteFailed              = errors.New("Execute error")
	ErrCompileFailed              = errors.New("Compilation error")
	ErrLanguageNotSupported       = errors.New("This language is not supported")
	ErrDebuggerIsClosed           = errors.New("debug is closed")
	ErrProgramIsRunningOptionFail = errors.New("The program is running")
)
