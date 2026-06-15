package authchallengeprovider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lwmacct/260605-miniport/internal/domain/authchallenge"
	"github.com/lwmacct/260605-miniport/internal/domain/identitychallenge"
)

type ImageProvider struct {
	challenges *identitychallenge.Service
}

func NewImageProvider(maxChallenges int) *ImageProvider {
	return &ImageProvider{challenges: identitychallenge.NewService(maxChallenges)}
}

func (p *ImageProvider) Name() string {
	return authchallenge.ProviderImage
}

func (p *ImageProvider) PublicConfig() authchallenge.PublicConfigDTO {
	return authchallenge.PublicConfigDTO{Provider: authchallenge.ProviderImage}
}

func (p *ImageProvider) Create(context.Context, authchallenge.RequestDTO) (*authchallenge.ChallengeDTO, error) {
	id, image, expiresAt, err := p.challenges.Generate()
	if err != nil {
		if errors.Is(err, identitychallenge.ErrLimitExceeded) {
			return nil, authchallenge.ErrLimitExceeded
		}
		return nil, err
	}
	return &authchallenge.ChallengeDTO{
		Provider:    authchallenge.ProviderImage,
		ChallengeID: id,
		Image:       image,
		ExpiresAt:   expiresAt,
	}, nil
}

func (p *ImageProvider) Verify(_ context.Context, response authchallenge.ResponseDTO, _ authchallenge.RequestDTO) error {
	if !p.challenges.VerifyAndDelete(response.ChallengeID, response.Answer) {
		return authchallenge.ErrInvalidChallenge
	}
	return nil
}

type RemoteTokenProvider struct {
	provider  string
	siteKey   string
	secret    string
	verifyURL string
	client    *http.Client
}

func NewRemoteTokenProvider(provider string, siteKey string, secret string, verifyURL string) (*RemoteTokenProvider, error) {
	if strings.TrimSpace(provider) == "" || strings.TrimSpace(siteKey) == "" || strings.TrimSpace(secret) == "" || strings.TrimSpace(verifyURL) == "" {
		return nil, authchallenge.ErrChallengeUnsupported
	}
	return &RemoteTokenProvider{
		provider:  strings.TrimSpace(provider),
		siteKey:   strings.TrimSpace(siteKey),
		secret:    strings.TrimSpace(secret),
		verifyURL: strings.TrimSpace(verifyURL),
		client:    &http.Client{Timeout: 5 * time.Second},
	}, nil
}

func (p *RemoteTokenProvider) Name() string {
	return p.provider
}

func (p *RemoteTokenProvider) PublicConfig() authchallenge.PublicConfigDTO {
	return authchallenge.PublicConfigDTO{Provider: p.provider, SiteKey: p.siteKey}
}

func (p *RemoteTokenProvider) Create(context.Context, authchallenge.RequestDTO) (*authchallenge.ChallengeDTO, error) {
	return nil, authchallenge.ErrChallengeUnsupported
}

func (p *RemoteTokenProvider) Verify(ctx context.Context, response authchallenge.ResponseDTO, request authchallenge.RequestDTO) error {
	if strings.TrimSpace(response.Token) == "" {
		return authchallenge.ErrInvalidChallenge
	}
	form := url.Values{}
	form.Set("secret", p.secret)
	form.Set("response", response.Token)
	if request.IP != "" {
		form.Set("remoteip", request.IP)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return authchallenge.ErrInvalidChallenge
	}
	var body struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}
	if !body.Success {
		return authchallenge.ErrInvalidChallenge
	}
	return nil
}
