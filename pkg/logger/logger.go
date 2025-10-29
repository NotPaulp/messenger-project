package logger

import (
	"log"
	"os"
)

// Logger wraps the standard logger with levels
type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
}

// New creates a new logger
func New(debugMode bool) *Logger {
	flags := log.Ldate | log.Ltime

	logger := &Logger{
		infoLog:  log.New(os.Stdout, "INFO: ", flags),
		errorLog: log.New(os.Stderr, "ERROR: ", flags),
	}

	if debugMode {
		logger.debugLog = log.New(os.Stdout, "DEBUG: ", flags)
	}

	return logger
}

// Info logs informational messages
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}

// Error logs error messages
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}

// Debug logs debug messages (only in debug mode)
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debugLog != nil {
		l.debugLog.Printf(format, v...)
	}
}
