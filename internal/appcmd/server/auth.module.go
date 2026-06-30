package server

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/260630-go-hsr-auth/pkg/auth"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/appmodule"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"

	"github.com/lwmacct/260605-miniport/internal/config"
)

type AuthModule struct {
	cfg   *config.Config
	value *auth.Module
}

var _ appmodule.Module = (*AuthModule)(nil)
var _ identity.SessionResolver = (*AuthModule)(nil)
var _ identity.Directory = (*AuthModule)(nil)

func NewAuthSpec(cfg *config.Config) appmodule.Spec {
	module := &AuthModule{cfg: cfg}
	return appmodule.Spec{
		Name:   module.Name(),
		Schema: auth.ApplySchema,
		Build: func(ctx *appmodule.Context) (appmodule.Module, error) {
			module := &AuthModule{cfg: cfg}
			authModule, err := auth.New(auth.Options{
				DB:         ctx.DB(),
				Config:     module.config(),
				SessionTTL: cfg.Server.Auth.Session.TTL,
			})
			if err != nil {
				return nil, err
			}
			module.value = authModule
			return module, nil
		},
	}
}

func (m *AuthModule) Name() string {
	return "auth"
}

func (m *AuthModule) Register(api huma.API) {
	m.value.Register(api)
}

func (m *AuthModule) CurrentUser(ctx context.Context, sessionID string, request auth.SessionRequest) (*auth.User, error) {
	return m.value.CurrentUser(ctx, sessionID, request)
}

func (m *AuthModule) CurrentPrincipal(ctx context.Context, sessionID string, request requestctx.Request) (*identity.Principal, error) {
	return m.value.CurrentPrincipal(ctx, sessionID, request)
}

func (m *AuthModule) CreateExternalUser(ctx context.Context, input auth.ExternalUserInput) (*auth.User, error) {
	return m.value.CreateExternalUser(ctx, input)
}

func (m *AuthModule) CreateSession(ctx context.Context, userID string, request auth.SessionRequest) (*auth.Session, error) {
	return m.value.CreateSession(ctx, userID, request)
}

func (m *AuthModule) Principal(ctx context.Context, subject string) (*identity.Principal, error) {
	return m.value.Principal(ctx, subject)
}

func (m *AuthModule) Principals(ctx context.Context, subjects []string) (map[string]*identity.Principal, error) {
	return m.value.Principals(ctx, subjects)
}

func (m *AuthModule) config() auth.Config {
	return auth.Config{
		Local: auth.LocalConfig{
			LoginEnabled:        m.cfg.Server.Auth.Local.LoginEnabled,
			RegistrationEnabled: m.cfg.Server.Auth.Local.RegistrationEnabled,
		},
		Session: auth.SessionConfig{
			Cookie: auth.SessionCookieConfig{
				Name:     m.cfg.Server.Auth.Session.Cookie.Name,
				Path:     m.cfg.Server.Auth.Session.Cookie.Path,
				Secure:   m.cfg.Server.Auth.Session.Cookie.Secure,
				SameSite: http.SameSiteStrictMode,
			},
		},
		RuntimeAdmins:     m.cfg.Server.Auth.Admins,
		ChallengeProvider: m.challengeProvider(),
		Request:           auth.RequestFromContext,
	}
}

func (m *AuthModule) challengeProvider() challenge.Provider {
	cfg := m.cfg.Server.Auth.Challenge
	switch cfg.Provider {
	case challenge.ProviderHCaptcha:
		provider, err := challenge.NewRemoteTokenProvider(
			challenge.ProviderHCaptcha,
			cfg.HCaptcha.SiteKey,
			cfg.HCaptcha.Secret,
			cfg.HCaptcha.VerifyURL,
		)
		if err == nil {
			return provider
		}
	case challenge.ProviderTurnstile:
		provider, err := challenge.NewRemoteTokenProvider(
			challenge.ProviderTurnstile,
			cfg.Turnstile.SiteKey,
			cfg.Turnstile.Secret,
			cfg.Turnstile.VerifyURL,
		)
		if err == nil {
			return provider
		}
	}
	return challenge.NewImageProvider(cfg.Image.MaxChallenges)
}
