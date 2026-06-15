package authsession

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"

	"github.com/lwmacct/260605-miniport/internal/domain/identityuser"
)

const DefaultTTL = 7 * 24 * time.Hour

type Service struct {
	db  bun.IDB
	ttl time.Duration
}

func NewService(db bun.IDB, ttl time.Duration) *Service {
	if db == nil {
		panic("authsession.NewService: db is nil")
	}
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Service{db: db, ttl: ttl}
}

func (svc *Service) Create(ctx context.Context, userID int64, request Request) (string, time.Time, error) {
	sessionID, err := newSessionID()
	if err != nil {
		return "", time.Time{}, err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(svc.ttl)
	row := &Session{
		IDHash:        tokenHash(sessionID),
		UserID:        userID,
		LoginIP:       request.IP,
		LastIP:        request.IP,
		UserAgentHash: tokenHash(request.UserAgent),
		ExpiresAt:     expiresAt,
		CreatedAt:     now,
		LastSeenAt:    now,
	}
	_, err = svc.db.NewInsert().Model(row).Exec(ctx)
	return sessionID, expiresAt, err
}

func (svc *Service) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	_, err := svc.db.NewDelete().Model((*Session)(nil)).Where("id_hash = ?", tokenHash(sessionID)).Exec(ctx)
	return err
}

func (svc *Service) User(ctx context.Context, sessionID string, request Request, users *identityuser.Service) (*UserDTO, error) {
	row := new(Session)
	if err := svc.db.NewSelect().
		Model(row).
		Where("id_hash = ?", tokenHash(sessionID)).
		Where("revoked_at IS NULL").
		Scan(ctx); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if !row.ExpiresAt.After(now) {
		_ = svc.Delete(ctx, sessionID)
		return nil, sql.ErrNoRows
	}

	user, err := users.ByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}
	if ensureErr := users.EnsureActive(user); ensureErr != nil {
		return nil, sql.ErrNoRows
	}

	_, err = svc.db.NewUpdate().
		Model((*Session)(nil)).
		Set("last_seen_at = ?", now).
		Set("last_ip = ?", request.IP).
		Where("id_hash = ?", tokenHash(sessionID)).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		ExpiresAt: row.ExpiresAt,
	}, nil
}
