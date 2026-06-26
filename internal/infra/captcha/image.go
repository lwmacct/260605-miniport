package captcha

import (
	"context"
	"crypto/sha256"
	"image/color"
	"strings"
	"sync"
	"time"

	"github.com/golang-module/base64Captcha/driver"
	"github.com/lwmacct/260605-miniport/internal/infra/token"
	"github.com/lwmacct/260605-miniport/internal/service"
)

type ImageProvider struct {
	challenges *imageChallengeStore
}

func NewImageProvider(maxChallenges int) *ImageProvider {
	return &ImageProvider{challenges: newImageChallengeStore(maxChallenges)}
}

func (p *ImageProvider) Name() string {
	return service.AuthChallengeProviderImage
}

func (p *ImageProvider) PublicConfig() service.AuthChallengePublicConfig {
	return service.AuthChallengePublicConfig{Provider: service.AuthChallengeProviderImage}
}

func (p *ImageProvider) Create(context.Context, service.AuthChallengeInput) (*service.AuthChallenge, error) {
	id, image, expiresAt, err := p.challenges.Generate()
	if err != nil {
		return nil, err
	}
	return &service.AuthChallenge{
		Provider:    service.AuthChallengeProviderImage,
		ChallengeID: id,
		Image:       image,
		ExpiresAt:   expiresAt,
	}, nil
}

func (p *ImageProvider) Verify(_ context.Context, answer service.AuthChallengeAnswer, _ service.AuthChallengeInput) error {
	if !p.challenges.VerifyAndDelete(answer.ChallengeID, answer.Answer) {
		return service.ErrAuthChallengeInvalid
	}
	return nil
}

type imageChallengeStore struct {
	mu         sync.Mutex
	challenges map[string]imageChallengeItem
	driver     driver.Driver
	ttl        time.Duration
	maxItems   int
}

type imageChallengeItem struct {
	answerHash [32]byte
	expiresAt  time.Time
}

func newImageChallengeStore(maxItems int) *imageChallengeStore {
	return &imageChallengeStore{
		challenges: make(map[string]imageChallengeItem),
		driver: driver.NewDriverString(driver.DriverString{
			Width:           180,
			Height:          56,
			Length:          4,
			NoiseCount:      12,
			ShowLineOptions: driver.OptionShowHollowLine | driver.OptionShowSlimeLine,
			Source:          "23456789ABCDEFGHJKLMNPQRSTUVWXYZ",
			BgColor:         &color.RGBA{R: 248, G: 250, B: 252, A: 255},
		}),
		ttl:      5 * time.Minute,
		maxItems: maxItems,
	}
}

func (s *imageChallengeStore) Generate() (string, string, time.Time, error) {
	_, content, answer := s.driver.GenerateCaptcha()
	image, err := s.driver.DrawCaptcha(content)
	if err != nil {
		return "", "", time.Time{}, err
	}
	id, expiresAt, err := s.Put(answer)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return id, image.Encoder(), expiresAt, nil
}

func (s *imageChallengeStore) Put(answer string) (string, time.Time, error) {
	id, err := token.NewWithPrefix("cap")
	if err != nil {
		return "", time.Time{}, err
	}
	now := time.Now().UTC()
	expiresAt := now.Add(s.ttl)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupLocked(now)
	if s.maxItems > 0 && len(s.challenges) >= s.maxItems {
		return "", time.Time{}, service.ErrAuthChallengeLimitExceeded
	}
	s.challenges[id] = imageChallengeItem{answerHash: imageAnswerHash(answer), expiresAt: expiresAt}
	return id, expiresAt, nil
}

func (s *imageChallengeStore) VerifyAndDelete(id string, answer string) bool {
	id = strings.TrimSpace(id)
	if id == "" || strings.TrimSpace(answer) == "" {
		return false
	}
	now := time.Now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	challenge, ok := s.challenges[id]
	if !ok {
		return false
	}
	delete(s.challenges, id)
	if !challenge.expiresAt.After(now) {
		return false
	}
	return challenge.answerHash == imageAnswerHash(answer)
}

func (s *imageChallengeStore) cleanupLocked(now time.Time) {
	for id, challenge := range s.challenges {
		if !challenge.expiresAt.After(now) {
			delete(s.challenges, id)
		}
	}
}

func imageAnswerHash(answer string) [32]byte {
	return sha256.Sum256([]byte(strings.ToUpper(strings.TrimSpace(answer))))
}
