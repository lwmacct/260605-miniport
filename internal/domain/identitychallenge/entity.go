package identitychallenge

import "time"

type item struct {
	answerHash [32]byte
	expiresAt  time.Time
}
