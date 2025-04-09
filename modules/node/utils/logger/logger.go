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

// Debug logs a message at the debug level with a bug emoji
func Debug(message string) {
	logMessage(DebugLevel, Blue, "🐞  "+message, true)
}

// Info logs a message at the info level with an information emoji
func Info(message string) {
	logMessage(InfoLevel, Green, "ℹ️  "+message, true)
}

// Warn logs a message at the warning level with a warning emoji
func Warn(message string) {
	logMessage(WarnLevel, Yellow, "⚠️  "+message, true)
}

// Error logs an error message with file and line number and a cross mark emoji
func Error(message string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		message = message + " (at " + file + ":" + strconv.Itoa(line) + ")"
	}
	logMessage(ErrorLevel, Red, "❌  "+message, true)
}

// Fatal logs a critical error message, adds a skull emoji, and exits the application
func Fatal(message string) {
	_, file, line, ok := runtime.Caller(1)
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()

	// // Simulate a long-running task
	// select {
	// case <-time.After(20 * time.Second):
	// 	fmt.Println("Task completed")
	// case <-ctx.Done():
	// 	fmt.Println("Timeout:", ctx.Err())
	// }
	if ok {
		message = message + " (at " + file + ":" + strconv.Itoa(line) + ")"
	}
	logMessage(ErrorLevel, Purple, "💀 "+message, true)
	log.Fatal(message)
}
