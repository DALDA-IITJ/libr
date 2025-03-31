package logger

import (
	"log"
	"runtime"
	"strconv"
	"time"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
)

// LogLevel represents the severity of the log message
type LogLevel string

const (
	DebugLevel LogLevel = "DEBUG"
	InfoLevel  LogLevel = "INFO"
	WarnLevel  LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
)

// logMessage formats and logs a message with the given level and color
func logMessage(level LogLevel, color string, message string, includeTimestamp bool) {
	timestamp := ""
	if includeTimestamp {
		timestamp = time.Now().Format("2006-01-02 15:04:05") + " "
	}
	log.Printf("\n%s%s[%s]%s %s%s\n", color, timestamp, level, Reset, message, Reset)
}

// Debug logs a message at the debug level
func Debug(message string) {
	logMessage(DebugLevel, Blue, message, true)
}

// Info logs a message at the info level
func Info(message string) {
	logMessage(InfoLevel, Green, message, true)
}

// Warn logs a message at the warning level
func Warn(message string) {
	logMessage(WarnLevel, Yellow, message, true)
}

// Error logs an error message with file and line number
func Error(message string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		message = message + " (at " + file + ":" + strconv.Itoa(line) + ")"
	}
	logMessage(ErrorLevel, Red, message, true)
}

// Fatal logs a critical error message and exits the application
func Fatal(message string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		message = message + " (at " + file + ":" + strconv.Itoa(line) + ")"
	}
	logMessage(ErrorLevel, Purple, message, true)
	log.Fatal(message)
}
