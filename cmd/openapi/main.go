package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/loopopen/t-ddd-gin-huma-gorm/cmd/api/config"
	apihttp "github.com/loopopen/t-ddd-gin-huma-gorm/internal/adapters/http"
)

func main() {
	env := &config.Env{Name: "dev"}
	c := &config.Config{
		OpenAPI: config.OpenAPI{
			ServerURL: "http://localhost:8080",
			DocsPath:  "/docs",
			SpecPath:  "/openapi",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	spec := apihttp.OpenAPISpec(env, c, nil, logger)

	jsonDoc, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		panic(err)
	}
	yamlDoc, err := spec.YAML()
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll("docs", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("docs/openapi.json", append(jsonDoc, '\n'), 0o644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("docs/openapi.yaml", yamlDoc, 0o644); err != nil {
		panic(err)
	}
}
