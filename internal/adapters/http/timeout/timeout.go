package timeout

import (
	"context"
	"errors"
	"time"
)

type result[T any] struct {
	value T
	err   error
}

func Run[T any](ctx context.Context, limit time.Duration, work func(context.Context) (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, limit)
	defer cancel()

	done := make(chan result[T], 1)
	go func() {
		value, err := work(ctx)
		done <- result[T]{value: value, err: err}
	}()

	select {
	case output := <-done:
		return output.value, output.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

func IsTimeout(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}
