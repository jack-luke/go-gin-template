package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

// TestLivenessProbe checks that Kubernetes healthchecks succeed
func TestLivenessProbe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, err := setupRouter()
	if err != nil {
		t.Fatalf("Error setting up router: %v", err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	// verify response is 200 OK
	assert.Equal(t, 200, w.Code)
}

// TestNoRoute verifies that unhandled routes are correctly reported
func TestNoRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, err := setupRouter()
	if err != nil {
		t.Fatalf("Error setting up router: %v", err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/this-is-not-a-route", nil)
	router.ServeHTTP(w, req)

	// Verify not found status
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestPanic forces a panic and verifies server recovery
func TestPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, err := setupRouter()
	if err != nil {
		t.Fatalf("Error setting up router: %v", err)
		return
	}

	// Attach an endpoint that panics
	router.GET("/panic", func(c *gin.Context) {
		panic(errors.New("this is a test panic"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	// Verify internal server error status
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Test that an error is caught, logged, and relevant info returned to user
func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture JSON logs in a buffer
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))

	router, err := setupRouter()
	if err != nil {
		t.Fatalf("Error setting up router: %v", err)
		return
	}

	// Attach an endpoint that halts the request with an error
	errMsg := "this is a test error"
	router.GET("/error", func(c *gin.Context) {
		c.AbortWithError(http.StatusInternalServerError, errors.New(errMsg))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	router.ServeHTTP(w, req)

	// Test that internal server error is returned
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Test that the error is logged
	if !strings.Contains(buf.String(), errMsg) {
		t.Fatalf("Log: \n %s \n did not contain error: %s", buf.String(), errMsg)
	}

	// Test that `the response contains the error message
	if !strings.Contains(w.Body.String(), errMsg) {
		t.Fatalf("Response: \n %s \n did not contain error: %s", w.Body.String(), errMsg)
	}
}

type RequestLog struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Error    string `json:"error"`
	ClientIP string `json:"client_ip"`
	Msg      string `json:"msg"`
	Status   int    `json:"status"`
	Level    string `json:"level"`
}

// TestLoggingFields that request logs have the fields and values expected by
// downstream consumers.
func TestLoggingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// capture JSON logs in a buffer
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))

	router, err := setupRouter()
	if err != nil {
		t.Fatalf("Error setting up router: %v", err)
		return
	}

	// Make request to an endpoint to produce a request log
	req, _ := http.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify structure of JSON log
	var logOutput RequestLog
	err = json.Unmarshal(buf.Bytes(), &logOutput)
	if err != nil {
		t.Fatalf("Failed to unmarshal request log into expected structure: %v", err)
	}

	// Verify all expected log values
	if logOutput.Level != "INFO" {
		t.Errorf("Request log field 'level', want %v, got %v", "INFO", logOutput.Level)
	}
	if logOutput.Method != "GET" {
		t.Errorf("Request log field 'method', want %v, got %v", "GET", logOutput.Method)
	}
	if logOutput.Path != "/healthz" {
		t.Errorf("Request log field 'path', want %v, got %v", "/healthz", logOutput.Path)
	}
	if logOutput.Status != 200 {
		t.Errorf("Request log field 'status', want %v, got %v", 200, logOutput.Status)
	}
	if logOutput.Error != "" {
		t.Errorf("Request log field 'error', want %v, got %v", "", logOutput.Error)
	}
	if logOutput.Msg != "HTTP 200 (OK)" {
		t.Errorf("Request log field 'msg', want %v, got %v", "HTTP 200 (OK)", logOutput.Msg)
	}
}
