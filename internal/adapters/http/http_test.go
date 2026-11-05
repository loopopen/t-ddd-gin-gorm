package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loopopen/t-ddd-gin-gorm/cmd/api/config"
	"github.com/loopopen/t-ddd-gin-gorm/pkg/schema/errx"
)

func testApp(t *testing.T) http.Handler {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewApp(
		&config.Env{Name: "test", CommitSHA: "abc123"},
		&config.Config{},
		nil,
		nil,
		nil,
		logger,
	).Engine
}

func TestOpenAPIEndpoints(t *testing.T) {
	app := testApp(t)

	for _, path := range []string{"/docs", "/openapi.json", "/openapi.yaml"} {
		t.Run(path, func(t *testing.T) {
			response := httptest.NewRecorder()
			app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
			if response.Code != http.StatusOK {
				t.Fatalf("GET %s returned %d: %s", path, response.Code, response.Body.String())
			}
		})
	}
}

func TestOpenAPIDocumentsUsersRoute(t *testing.T) {
	response := httptest.NewRecorder()
	testApp(t).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/openapi.json", nil))

	var document struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
		Paths map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	if document.Info.Version != "1.0-abc123" {
		t.Fatalf("unexpected version %q", document.Info.Version)
	}
	if _, ok := document.Paths["/api/v1/users"]; !ok {
		t.Fatal("users route is missing from OpenAPI")
	}
}

func TestMissingQueryUsesErrorEnvelope(t *testing.T) {
	response := httptest.NewRecorder()
	testApp(t).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/users", nil))

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected status %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"reason":"FRAMEWORK_ERROR"`) {
		t.Fatalf("unexpected error response: %s", response.Body.String())
	}
}

func TestErrorHandlerWritesLastGinError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	engine := gin.New()
	engine.Use(ErrorHandler(logger))
	engine.GET("/", func(c *gin.Context) {
		_ = c.Error(errx.ErrInvalidRequestFields)
	})

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"reason":"INVALID_REQUEST_FIELDS"`) {
		t.Fatalf("unexpected error response: %s", response.Body.String())
	}
}

func TestHumaErrorFactoryUsesApplicationErrorEnvelope(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	err := newHumaError(logger)(http.StatusInternalServerError, "unexpected error", errx.ErrInvalidRequestFields)

	if err.GetStatus() != http.StatusBadRequest {
		t.Fatalf("unexpected status %d", err.GetStatus())
	}
	if !strings.Contains(err.Error(), "invalid request fields") {
		t.Fatalf("unexpected error %q", err.Error())
	}
}

func TestWithTimeoutReturnsWithoutWaitingForUncooperativeHandler(t *testing.T) {
	handlerDone := make(chan struct{})
	handler := func(context.Context, *struct{}) (*struct{}, error) {
		defer close(handlerDone)
		time.Sleep(200 * time.Millisecond)
		return &struct{}{}, nil
	}

	start := time.Now()
	_, err := withTimeout(20*time.Millisecond, handler)(context.Background(), &struct{}{})
	if elapsed := time.Since(start); elapsed >= 100*time.Millisecond {
		t.Fatalf("returned after %s; want timeout to return promptly", elapsed)
	}
	if !errors.Is(err, errx.ErrServiceTimeout) {
		t.Fatalf("error = %v, want service timeout", err)
	}

	<-handlerDone
}
