package httpauth

import (
	"context"

	"github.com/lwmacct/260605-miniport/internal/domain/authsession"
)

type requestKey struct{}
type sessionIDKey struct{}
type userKey struct{}

func ContextWithRequest(ctx context.Context, request authsession.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, request)
}

func RequestFromContext(ctx context.Context) (authsession.Request, bool) {
	request, ok := ctx.Value(requestKey{}).(authsession.Request)
	return request, ok
}

func ContextWithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDKey{}, sessionID)
}

func SessionIDFromContext(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(sessionIDKey{}).(string)
	return sessionID, ok && sessionID != ""
}

func ContextWithAuthActor(ctx context.Context, user AuthActor) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

func AuthActorFromContext(ctx context.Context) (AuthActor, bool) {
	user, ok := ctx.Value(userKey{}).(AuthActor)
	return user, ok
}
