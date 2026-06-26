package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

func utilRequest(ctx context.Context, config Config) (service.AuthSessionInput, error) {
	if config.Request == nil {
		return service.AuthSessionInput{}, huma.Error400BadRequest("invalid request source")
	}
	request, ok := config.Request(ctx)
	if !ok {
		return service.AuthSessionInput{}, huma.Error400BadRequest("invalid request source")
	}
	return request, nil
}

func utilRegisterErrorMessage(err error) string {
	switch {
	case errors.Is(err, service.ErrAuthPasswordWeakPassword):
		return "weak password"
	case errors.Is(err, service.ErrUserUsernameTaken):
		return "username taken"
	case errors.Is(err, service.ErrUserInvalidCredentials):
		return "invalid credentials"
	default:
		return "register failed"
	}
}

func utilSessionCookieValue(value string, expiresAt time.Time, secure bool) string {
	//nolint:gosec // Secure is configuration-driven so local HTTP development can work.
	return (&http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/api",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	}).String()
}

func utilClearSessionCookie(secure bool) string {
	//nolint:gosec // Secure is configuration-driven so local HTTP development can work.
	return (&http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/api",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
	}).String()
}

func utilIsRuntimeAdmin(user *service.User, admins []string) bool {
	if user == nil {
		return false
	}
	return utilIsRuntimeAdminUsername(user.Username, admins)
}

func utilIsRuntimeAdminUsername(username string, admins []string) bool {
	for _, admin := range admins {
		if strings.EqualFold(strings.TrimSpace(admin), username) {
			return true
		}
	}
	return false
}

func utilHTTPTime(value time.Time) string {
	return value.UTC().Format(http.TimeFormat)
}
