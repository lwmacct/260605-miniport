package authchallenge

import "errors"

var (
	ErrInvalidChallenge     = errors.New("invalid challenge")
	ErrChallengeUnsupported = errors.New("challenge provider unsupported")
	ErrLimitExceeded        = errors.New("challenge limit exceeded")
)
