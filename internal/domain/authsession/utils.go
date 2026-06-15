package authsession

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func newSessionID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "sess_" + base64.RawURLEncoding.EncodeToString(buf), nil
}

func tokenHash(value string) []byte {
	sum := sha256.Sum256([]byte(value))
	return sum[:]
}
