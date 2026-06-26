package service

import "context"

type AuthChallengeProvider interface {
	Name() string
	PublicConfig() AuthChallengePublicConfig
	Create(context.Context, AuthChallengeInput) (*AuthChallenge, error)
	Verify(context.Context, AuthChallengeAnswer, AuthChallengeInput) error
}
