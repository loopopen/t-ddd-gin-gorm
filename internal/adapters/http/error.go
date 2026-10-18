package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	pkgerrorx "github.com/loopopen/pkg/errorx"
	"github.com/loopopen/pkg/slogx"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/adapters/http/dto"
	"github.com/loopopen/t-ddd-gin-huma-gorm/pkg/schema/errx"
)

type ErrorHandler func(error) huma.StatusError

func HandlerError(logger *slog.Logger) ErrorHandler {
	logger = slogx.Within(logger, HandlerError)

	return func(theErr error) huma.StatusError {
		var appErr *pkgerrorx.Error
		if !errors.As(theErr, &appErr) {
			appErr = errx.ErrUnspecified.WithCause(theErr)
		}
		if appErr.Code >= http.StatusInternalServerError {
			logger.Error("request failed", slog.Any("err", theErr))
		}
		return dto.Err(appErr)
	}
}

func NewHumaError(status int, message string, details ...error) huma.StatusError {
	metadata := make(map[string]string, len(details))
	for i, detail := range details {
		if detail != nil {
			metadata[fmt.Sprintf("error_%d", i+1)] = detail.Error()
		}
	}
	return dto.NewError(status, "FRAMEWORK_ERROR", message, "", metadata)
}

func ErrorMiddleware(logger *slog.Logger) gin.HandlerFunc {
	handle := HandlerError(logger)
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				writeGinError(c, handle(fmt.Errorf("panic: %v", recovered)))
			}
		}()

		c.Next()
		if !c.Writer.Written() && len(c.Errors) > 0 {
			writeGinError(c, handle(c.Errors.Last().Err))
		}
	}
}

func writeGinError(c *gin.Context, err huma.StatusError) {
	if c.Writer.Written() {
		return
	}
	c.AbortWithStatusJSON(err.GetStatus(), err)
}
