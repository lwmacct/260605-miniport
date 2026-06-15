package authchallenge

import "context"

type Provider interface {
	Name() string
	PublicConfig() PublicConfigDTO
	Create(context.Context, RequestDTO) (*ChallengeDTO, error)
	Verify(context.Context, ResponseDTO, RequestDTO) error
}

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
	return &Service{provider: provider}
}

func (s *Service) PublicConfig() PublicConfigDTO {
	if s == nil || s.provider == nil {
		return PublicConfigDTO{}
	}
	return s.provider.PublicConfig()
}

func (s *Service) Create(ctx context.Context, request RequestDTO) (*ChallengeDTO, error) {
	if s == nil || s.provider == nil {
		return nil, ErrChallengeUnsupported
	}
	return s.provider.Create(ctx, request)
}

func (s *Service) Verify(ctx context.Context, response ResponseDTO, request RequestDTO) error {
	if s == nil || s.provider == nil || response.Provider != s.provider.Name() {
		return ErrInvalidChallenge
	}
	if err := s.provider.Verify(ctx, response, request); err != nil {
		return ErrInvalidChallenge
	}
	return nil
}
