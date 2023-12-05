package scheduler

// Logger defines a go-scheduler logger
type Logger interface {
	// Error logs a message on error level
	Error(message string)

	// Errorf formats a message according to a format specifier and logs the message on error level
	Errorf(format string, a ...any)
}

// emptyLogger represents an empty logger that logs nothing
type emptyLogger struct{}

func (el *emptyLogger) Error(message string) {}

func (el *emptyLogger) Errorf(format string, a ...any) {}
