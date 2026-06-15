package identitychallenge

import (
	"crypto/sha256"
	"strings"
)

func utilAnswerHash(answer string) [32]byte {
	return sha256.Sum256([]byte(strings.ToUpper(strings.TrimSpace(answer))))
}
