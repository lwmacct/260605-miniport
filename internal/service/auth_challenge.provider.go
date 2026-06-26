package service

import "context"

type AuthChallengeProvider interface {
	Name() string
	PublicConfig() AuthChallengePublicConfig
	Create(ctx context.Context, request AuthChallengeInput) (*AuthChallenge, error)
	Verify(ctx context.Context, answer AuthChallengeAnswer, request AuthChallengeInput) error
}
