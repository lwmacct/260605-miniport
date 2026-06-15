package httpauth

import (
	"context"

	"github.com/lwmacct/260605-miniport/internal/domain/authsession"
)

type requestKey struct{}
type userKey struct{}

func ContextWithRequest(ctx context.Context, request authsession.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, request)
}

func RequestFromContext(ctx context.Context) (authsession.Request, bool) {
	request, ok := ctx.Value(requestKey{}).(authsession.Request)
	return request, ok
}

func ContextWithAuthActor(ctx context.Context, user AuthActor) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

func AuthActorFromContext(ctx context.Context) (AuthActor, bool) {
	user, ok := ctx.Value(userKey{}).(AuthActor)
	return user, ok
}
