package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging(t *testing.T) {
	// Create a buffer to capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	// Disable timestamps in log output for testing
	log.SetFlags(0)

	// Create a dummy handler to pass to the middleware
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the dummy handler with the logging middleware
	handler := Logging(dummyHandler)

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler with the request and response recorder
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the log output
	expectedLog := "GET /test 127.0.0.1:12345\n"
	if logBuffer.String() != expectedLog {
		t.Errorf("expected log %q but got %q", expectedLog, logBuffer.String())
	}
}
