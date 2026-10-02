package backup

// diagnosticError contains only developer-controlled text. Never construct it
// from subprocess output, provider bodies, filesystem paths, or credentials.
type diagnosticError struct {
	code, message string
}

func (e *diagnosticError) Error() string { return e.message }

func preflightFailure(code, message string) error {
	return &diagnosticError{code: "preflight_" + code, message: message}
}
