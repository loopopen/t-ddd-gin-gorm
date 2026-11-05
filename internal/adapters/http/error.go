package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/loopopen/pkg/errorx"
	"github.com/loopopen/pkg/slogx"
	"github.com/loopopen/t-ddd-gin-gorm/internal/adapters/http/dto"
	"github.com/loopopen/t-ddd-gin-gorm/pkg/schema/errx"
)

// ErrorHandler collects errors added through c.Error and writes the API error
// response after the remaining handlers have run.
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	logger = slogx.Within(logger, ErrorHandler)

	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				writeGinError(c, toErrorResponse(logger, fmt.Errorf("panic: %v", recovered)))
			}
		}()

		c.Next()
		if len(c.Errors) > 0 {
			writeGinError(c, toErrorResponse(logger, c.Errors.Last().Err))
		}
	}
}

func toErrorResponse(logger *slog.Logger, theErr error) *dto.ErrorResponse {
	var appErr *errorx.Error
	if !errors.As(theErr, &appErr) {
		appErr = errx.ErrUnspecified.WithCause(theErr)
	}
	if appErr.Code >= http.StatusInternalServerError {
		logger.Error("request failed", slog.Any("err", theErr))
	}
	return dto.Err(appErr)
}

// newHumaError configures Huma's global error factory for this API. Returned
// application errors retain their status code and API error envelope.
func newHumaError(logger *slog.Logger) func(int, string, ...error) huma.StatusError {
	logger = slogx.Within(logger, newHumaError)

	return func(status int, message string, details ...error) huma.StatusError {
		for _, detail := range details {
			var appErr *errorx.Error
			if errors.As(detail, &appErr) {
				return toErrorResponse(logger, detail)
			}
		}

		metadata := make(map[string]string, len(details))
		for i, detail := range details {
			if detail != nil {
				metadata[fmt.Sprintf("error_%d", i+1)] = detail.Error()
			}
		}
		return dto.NewError(status, "FRAMEWORK_ERROR", message, "", metadata)
	}
}

func writeGinError(c *gin.Context, err *dto.ErrorResponse) {
	if c.Writer.Written() {
		return
	}
	c.AbortWithStatusJSON(err.GetStatus(), err)
}
