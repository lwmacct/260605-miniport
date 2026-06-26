package handler

import (
	"context"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type requestKey struct{}

func ContextWithRequest(ctx context.Context, request service.AuthSessionInput) context.Context {
	return context.WithValue(ctx, requestKey{}, request)
}

func RequestFromContext(ctx context.Context) (service.AuthSessionInput, bool) {
	request, ok := ctx.Value(requestKey{}).(service.AuthSessionInput)
	return request, ok
}
