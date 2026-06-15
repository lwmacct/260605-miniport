package identitychallenge

import (
	"image/color"
	"strings"
	"sync"
	"time"

	"github.com/golang-module/base64Captcha/driver"

	"github.com/lwmacct/260605-miniport/internal/infra/token"
)

type Service struct {
	mu         sync.Mutex
	challenges map[string]item
	driver     driver.Driver
	ttl        time.Duration
	maxItems   int
}

func NewService(maxItems int) *Service {
	return &Service{
		challenges: make(map[string]item),
		driver: driver.NewDriverString(driver.DriverString{
			Width:           180,
			Height:          56,
			Length:          4,
			NoiseCount:      12,
			ShowLineOptions: driver.OptionShowHollowLine | driver.OptionShowSlimeLine,
			Source:          "23456789ABCDEFGHJKLMNPQRSTUVWXYZ",
			BgColor:         &color.RGBA{R: 248, G: 250, B: 252, A: 255},
		}),
		ttl:      defaultTTL,
		maxItems: maxItems,
	}
}

func (s *Service) Generate() (string, string, time.Time, error) {
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

func (s *Service) Put(answer string) (string, time.Time, error) {
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
		return "", time.Time{}, ErrLimitExceeded
	}
	s.challenges[id] = item{answerHash: utilAnswerHash(answer), expiresAt: expiresAt}
	return id, expiresAt, nil
}

func (s *Service) VerifyAndDelete(id string, answer string) bool {
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
	return challenge.answerHash == utilAnswerHash(answer)
}

func (s *Service) cleanupLocked(now time.Time) {
	for id, challenge := range s.challenges {
		if !challenge.expiresAt.After(now) {
			delete(s.challenges, id)
		}
	}
}
