package service

import (
	"context"
	"time"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

type AuthSessionService struct {
	store *repository.Store
	ttl   time.Duration
}

func NewAuthSessionService(store *repository.Store, ttl time.Duration) *AuthSessionService {
	if store == nil {
		panic("NewAuthSessionService: store is nil")
	}
	if ttl <= 0 {
		ttl = AuthSessionDefaultTTL
	}
	return &AuthSessionService{store: store, ttl: ttl}
}

func (s *AuthSessionService) Create(ctx context.Context, userID int64, request AuthSessionInput) (string, time.Time, error) {
	sessionID, err := utilNewSessionID()
	if err != nil {
		return "", time.Time{}, err
	}
	now := time.Now().UTC()
	expiresAt := now.Add(s.ttl)
	_, err = s.store.CreateAuthSessionFromInput(ctx, repository.AuthSessionCreateInput{
		IDHash:        utilTokenHash(sessionID),
		UserID:        userID,
		LoginIP:       request.IP,
		LastIP:        request.IP,
		UserAgentHash: utilTokenHash(request.UserAgent),
		ExpiresAt:     expiresAt,
		CreatedAt:     now,
		LastSeenAt:    now,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return sessionID, expiresAt, nil
}

func (s *AuthSessionService) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	return s.store.DeleteAuthSessionByHash(ctx, utilTokenHash(sessionID))
}

func (s *AuthSessionService) User(ctx context.Context, sessionID string, request AuthSessionInput, users *UserService) (*AuthSessionUser, error) {
	row, err := s.store.FetchAuthSessionByHash(ctx, utilTokenHash(sessionID))
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	if !row.ExpiresAt.After(now) {
		_ = s.Delete(ctx, sessionID)
		return nil, repository.ErrNotFound
	}
	user, err := users.ByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}
	if err := users.EnsureActive(user); err != nil {
		return nil, repository.ErrNotFound
	}
	if _, err := s.store.UpdateAuthSessionTouch(ctx, row.IDHash, request.IP, now); err != nil {
		return nil, err
	}
	return &AuthSessionUser{ID: user.ID, Username: user.Username, ExpiresAt: row.ExpiresAt}, nil
}
