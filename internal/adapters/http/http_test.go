package http

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/loopopen/t-ddd-gin-huma-gorm/cmd/api/config"
)

func testApp(t *testing.T) http.Handler {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewApp(
		&config.Env{Name: "test", CommitSHA: "abc123"},
		&config.Config{},
		nil,
		nil,
		logger,
	)
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
