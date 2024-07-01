package handler

import (
	"context"

	"github.com/oktopUSP/backend/services/bulkdata/internal/bridge"
)

type Handler struct {
	ctx context.Context
	b   *bridge.Bridge
}

func NewHandler(ctx context.Context, b *bridge.Bridge) Handler {
	return Handler{
		ctx: ctx,
		b:   b,
	}
}
