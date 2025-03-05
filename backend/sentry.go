package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

// Default Sentry options
var (
	defaultSentryDSN         = ""
	defaultSentryEnvironment = "development"
	defaultSentryRelease     = "1.0.0"
)

// InitSentry initializes the Sentry SDK
func InitSentry() error {
	// Get configuration from environment variables or use defaults
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		dsn = defaultSentryDSN
		if dsn == "" {
			log.Println("Warning: SENTRY_DSN not set. Error reporting will be disabled.")
			return nil
		}
	}

	// Get environment and release
	environment := os.Getenv("SENTRY_ENVIRONMENT")
	if environment == "" {
		environment = defaultSentryEnvironment
	}

	release := os.Getenv("SENTRY_RELEASE")
	if release == "" {
		release = defaultSentryRelease
	}

	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		Release:          release,
		Debug:            os.Getenv("SENTRY_DEBUG") == "true",
		AttachStacktrace: true,
		// Set traces sample rate to capture some but not all requests
		TracesSampleRate: 0.2,
	})
	
	if err != nil {
		return fmt.Errorf("sentry initialization failed: %w", err)
	}

	log.Println("Sentry initialized successfully")
	log.Printf("Environment: %s, Release: %s", environment, release)
	
	// Flush buffered events before the program terminates
	defer sentry.Flush(2 * time.Second)
	
	return nil
}

// CaptureException sends an error to Sentry
func CaptureException(err error) *sentry.EventID {
	if err == nil {
		return nil
	}
	return sentry.CaptureException(err)
}

// CaptureMessage sends a message to Sentry
func CaptureMessage(message string) *sentry.EventID {
	return sentry.CaptureMessage(message)
}

// WithContext adds context data to the Sentry scope
func WithContext(f func(scope *sentry.Scope)) {
	sentry.WithScope(f)
}

// SetTag adds a tag to the current scope
func SetTag(key, value string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	})
}

// SetUser sets user information in the current scope
func SetUser(id, email, ipAddress string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:        id,
			Email:     email,
			IPAddress: ipAddress,
		})
	})
}

// SetExtra adds extra data to the current scope
func SetExtra(key string, value interface{}) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra(key, value)
	})
}

// RecoverPanic recovers from a panic and reports it to Sentry
func RecoverPanic() {
	if err := recover(); err != nil {
		eventID := sentry.CurrentHub().Recover(err)
		sentry.Flush(time.Second * 2)
		log.Printf("Recovered from panic: %v (event ID: %s)", err, *eventID)
	}
}

// SentryMiddleware creates HTTP middleware for Sentry integration
func SentryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub := sentry.CurrentHub().Clone()
		hub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetRequest(r)
			scope.SetTag("request_id", r.Header.Get("X-Request-ID"))
			scope.SetTag("endpoint", r.URL.Path)
			scope.SetUser(sentry.User{
				IPAddress: r.RemoteAddr,
			})
		})

		// Create a response writer that can track status code
		statusRecorder := &statusResponseWriter{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		// Use defer to recover from panics and report to Sentry
		defer func() {
			if err := recover(); err != nil {
				eventID := hub.RecoverWithContext(context.WithValue(r.Context(), sentry.RequestContextKey, r), err)
				log.Printf("Recovered from panic in HTTP handler: %v (event ID: %s)", err, *eventID)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Call the next handler
		next.ServeHTTP(statusRecorder, r)

		// Report non-2xx status codes to Sentry
		if statusRecorder.Status >= 500 {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelError)
				scope.SetExtra("status_code", statusRecorder.Status)
				hub.CaptureMessage(fmt.Sprintf("HTTP %d: %s", statusRecorder.Status, r.URL.Path))
			})
		}
	})
}

// statusResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type statusResponseWriter struct {
	http.ResponseWriter
	Status int
}

// WriteHeader captures the status code
func (w *statusResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

// Helper function to handle and report errors
func ReportError(err error, w http.ResponseWriter, message string, statusCode int) {
	if err != nil {
		CaptureException(err)
		log.Printf("Error: %v", err)
		sendErrorResponse(w, message, statusCode)
	}
} 