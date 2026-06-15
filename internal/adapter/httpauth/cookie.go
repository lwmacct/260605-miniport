package httpauth

import (
	"net/http"
	"time"
)

func SessionCookieValue(sessionID string, expiresAt time.Time, secure bool) *http.Cookie {
	// #nosec G124 -- Secure is controlled by the runtime TLS setting so local HTTP development can still log in.
	return &http.Cookie{
		Name:     SessionCookie,
		Value:    sessionID,
		Path:     SessionCookiePath,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	}
}

func ClearSessionCookie(secure bool) *http.Cookie {
	// #nosec G124 -- Secure must match the session cookie setting so browsers accept the clearing cookie.
	return &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     SessionCookiePath,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
	}
}
