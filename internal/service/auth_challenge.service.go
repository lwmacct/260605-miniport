package service

import "context"

type AuthChallengeService struct {
	provider AuthChallengeProvider
}

func NewAuthChallengeService(provider AuthChallengeProvider) *AuthChallengeService {
	return &AuthChallengeService{provider: provider}
}

func (s *AuthChallengeService) PublicConfig() AuthChallengePublicConfig {
	if s == nil || s.provider == nil {
		return AuthChallengePublicConfig{}
	}
	return s.provider.PublicConfig()
}

func (s *AuthChallengeService) Create(ctx context.Context, request AuthChallengeInput) (*AuthChallenge, error) {
	if s == nil || s.provider == nil {
		return nil, ErrAuthChallengeUnsupported
	}
	return s.provider.Create(ctx, request)
}

func (s *AuthChallengeService) Verify(ctx context.Context, answer AuthChallengeAnswer, request AuthChallengeInput) error {
	if s == nil || s.provider == nil || answer.Provider != s.provider.Name() {
		return ErrAuthChallengeInvalid
	}
	if err := s.provider.Verify(ctx, answer, request); err != nil {
		return ErrAuthChallengeInvalid
	}
	return nil
}
