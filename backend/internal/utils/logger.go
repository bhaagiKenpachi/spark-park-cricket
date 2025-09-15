package utils

import (
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance
var Logger *logrus.Logger
var loggerOnce sync.Once

// InitLogger initializes the structured logger
func InitLogger() {
	loggerOnce.Do(func() {
		Logger = logrus.New()

		// Set log level based on environment
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info"
		}

		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			level = logrus.InfoLevel
		}
		Logger.SetLevel(level)

		// Set JSON formatter for structured logging
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
			},
		})

		// Set output to stdout
		Logger.SetOutput(os.Stdout)
	})
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	InitLogger() // This is safe to call multiple times due to sync.Once
	return Logger
}

// LogWithContext creates a logger with context fields
func LogWithContext(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// LogError logs an error with context
func LogError(err error, message string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	fields["error"] = err.Error()
	GetLogger().WithFields(fields).Error(message)
}

// LogInfo logs an info message with context
func LogInfo(message string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	GetLogger().WithFields(fields).Info(message)
}

// LogWarn logs a warning message with context
func LogWarn(message string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	GetLogger().WithFields(fields).Warn(message)
}

// LogDebug logs a debug message with context
func LogDebug(message string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	GetLogger().WithFields(fields).Debug(message)
}

// LogHTTPRequest logs HTTP request details
func LogHTTPRequest(method, path, userAgent string, statusCode int, duration time.Duration, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}

	fields["method"] = method
	fields["path"] = path
	fields["user_agent"] = userAgent
	fields["status_code"] = statusCode
	fields["duration_ms"] = duration.Milliseconds()

	GetLogger().WithFields(fields).Info("HTTP request")
}

// LogWebSocketEvent logs WebSocket events
func LogWebSocketEvent(eventType, roomID, clientID string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}

	fields["event_type"] = eventType
	fields["room_id"] = roomID
	fields["client_id"] = clientID

	GetLogger().WithFields(fields).Info("WebSocket event")
}

// LogDatabaseOperation logs database operations
func LogDatabaseOperation(operation, table string, duration time.Duration, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}

	fields["operation"] = operation
	fields["table"] = table
	fields["duration_ms"] = duration.Milliseconds()

	GetLogger().WithFields(fields).Debug("Database operation")
}

// LogCricketEvent logs cricket-specific events
func LogCricketEvent(eventType, matchID string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}

	fields["cricket_event"] = eventType
	fields["match_id"] = matchID

	GetLogger().WithFields(fields).Info("Cricket event")
}
