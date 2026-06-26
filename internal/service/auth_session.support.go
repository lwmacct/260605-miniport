package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"
)

const AuthSessionDefaultTTL = 7 * 24 * time.Hour

type AuthSessionInput struct {
	IP                   string
	Host                 string
	UserAgent            string
	Method               string
	Path                 string
	RemoteAddr           string
	InvalidSessionCookie bool
}

type AuthSessionUser struct {
	ID        int64
	Username  string
	ExpiresAt time.Time
}

func utilNewSessionID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func utilTokenHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
