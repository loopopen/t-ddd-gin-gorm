package errx

import (
	"github.com/loopopen/pkg/errorx"
	"github.com/loopopen/t-ddd-gin-gorm/pkg/schema"
)

func FrameworkErr(status int, message string) *errorx.Error {
	return errorx.New(status, schema.ErrorReason_FRAMEWORK_ERROR.String(), message)
}

func ArgumentIsNil(argName string) *errorx.Error {
	return ErrNilArgument.WithMetadata(map[string]string{
		"argumentName": argName,
	})
}

var ErrUnspecified = errorx.InternalServer(schema.ErrorReason_UNSPECIFIED.String(), "")

var ErrInvalidIdempotencyKey = errorx.BadRequest(schema.ErrorReason_INVALID_IDEMPOTENCY_KEY.String(), "invalid idempotency key")
var ErrInvalidRequestFields = errorx.BadRequest(schema.ErrorReason_INVALID_REQUEST_FIELDS.String(), "invalid request fields")

var ErrNilArgument = errorx.InternalServer(schema.ErrorReason_NIL_ARGUMENT.String(), "nil argument")
var ErrServiceTimeout = errorx.GatewayTimeout(schema.ErrorReason_SERVICE_TIMEOUT.String(), "service timeout")
