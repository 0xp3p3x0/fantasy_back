package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	// LogFile is the file where all logs will be written
	LogFile *os.File
	// Logger is the logger instance that writes to both file and stdout
	Logger *log.Logger
)

// Init initializes the logger
func Init() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join("logs", fmt.Sprintf("app_%s.log", timestamp))

	// Open log file
	var err error
	LogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Create multi-writer to write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, LogFile)

	flags := log.Ldate | log.Ltime | log.Lshortfile

	// Route standard library log (log.Printf in middleware, db, etc.) to the same file + stdout
	log.SetOutput(multiWriter)
	log.SetFlags(flags)

	// Create logger with timestamp and file:line information
	Logger = log.New(multiWriter, "", flags)

	// Log initialization
	Logger.Printf("Logger initialized. Log file: %s", logPath)

	return nil
}

// Close closes the log file
func Close() error {
	if LogFile != nil {
		return LogFile.Close()
	}
	return nil
}

// LogAPIRequest logs an API request with method and path
func LogAPIRequest(r *http.Request, statusCode int, duration time.Duration) {
	if Logger != nil {
		Logger.Printf("[API] %s %s %d %v", r.Method, r.URL.Path, statusCode, duration)
	}
}

// LogAPIError logs an API error with method and path
func LogAPIError(r *http.Request, err error) {
	if Logger != nil {
		Logger.Printf("[API-ERROR] %s %s: %v", r.Method, r.URL.Path, err)
	}
}

// Printf logs a formatted message
func Printf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf(format, v...)
	}
}

// Print logs a message
func Print(v ...interface{}) {
	if Logger != nil {
		Logger.Print(v...)
	}
}

// Fatal logs a message and exits
func Fatal(v ...interface{}) {
	if Logger != nil {
		Logger.Fatal(v...)
	}
}

// Fatalf logs a formatted message and exits
func Fatalf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Fatalf(format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[ERROR] "+format, v...)
	}
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[INFO] "+format, v...)
	}
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[DEBUG] "+format, v...)
	}
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[WARN] "+format, v...)
	}
}
