package handlers

import "context"

type Handler[I, O any] = func(context.Context, *I) (*O, error)
