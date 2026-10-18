package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/loopopen/t-ddd-gin-huma-gorm/pkg/schema/errx"
)

type From int32

const (
	FromBody From = iota
	FromQuery
	FromParams
)

var validate = validator.New()

type Failure struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value any    `json:"value"`
}

func validateStruct(value any) []*Failure {
	var failures []*Failure
	if err := validate.Struct(value); err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) {
			failures = append(failures, &Failure{
				Field: validationErr.StructNamespace(),
				Tag:   validationErr.Tag(),
				Value: validationErr.Value(),
			})
		}
	}
	return failures
}

// Validate is a generic Gin middleware for endpoints that are not registered through Huma.
// Huma endpoints use Huma's built-in OpenAPI-driven parsing and validation.
func Validate[T any](source From) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := new(T)
		var err error
		switch source {
		case FromQuery:
			err = c.ShouldBindQuery(payload)
		case FromParams:
			err = c.ShouldBindUri(payload)
		case FromBody:
			err = c.ShouldBindJSON(payload)
		}
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		if failures := validateStruct(payload); len(failures) > 0 {
			metadata := make(map[string]string, len(failures))
			for _, failure := range failures {
				metadata[failure.Field] = fmt.Sprintf("failed on the '%s' tag, value: '%v'", failure.Tag, failure.Value)
			}
			_ = c.Error(errx.ErrInvalidRequestFields.WithMetadata(metadata))
			c.Abort()
			return
		}

		c.Set("payload", payload)
		c.Next()
	}
}
