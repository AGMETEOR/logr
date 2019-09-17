package logr

// Target represents a destination for log records such as file,
// database, TCP socket, etc.
type Target interface {
	// Start initializes the target and should start a new
	// goroutine to accept incoming log records.
	// This is a good place for a target to create a file,
	// connect to a database, or other data sink.
	Start() error

	// IsLevelEnabled returns true if this target should emit
	// logs for the specified level. Also determines if
	// a stack trace is required.
	IsLevelEnabled(Level) (enabled bool, stacktrace bool)

	// Formatter returns the Formatter associated with this Target.
	Formatter() Formatter

	// Log outputs the log record to this targets destination.
	Log(rec *LogRec)

	// Shutdown makes best effort to flush target queue and
	// frees/closes all resources.
	Shutdown() error
}