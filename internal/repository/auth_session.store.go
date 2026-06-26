package repository

import (
	"context"
	"time"
)

func (s *Store) CreateAuthSession(ctx context.Context, row *AuthSessionModel) (*AuthSessionModel, error) {
	if _, err := s.db.NewInsert().Model(row).Exec(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Store) CreateAuthSessionFromInput(ctx context.Context, input AuthSessionCreateInput) (*AuthSessionModel, error) {
	return s.CreateAuthSession(ctx, &AuthSessionModel{
		IDHash:        input.IDHash,
		UserID:        input.UserID,
		LoginIP:       input.LoginIP,
		LastIP:        input.LastIP,
		UserAgentHash: input.UserAgentHash,
		ExpiresAt:     input.ExpiresAt,
		CreatedAt:     input.CreatedAt,
		LastSeenAt:    input.LastSeenAt,
	})
}

func (s *Store) FetchAuthSessionByHash(ctx context.Context, idHash string) (*AuthSessionModel, error) {
	row := new(AuthSessionModel)
	err := s.db.NewSelect().
		Model(row).
		Where("id_hash = ?", idHash).
		Where("revoked_at IS NULL").
		Scan(ctx)
	return row, WrapNotFound(err)
}

func (s *Store) DeleteAuthSessionByHash(ctx context.Context, idHash string) error {
	_, err := s.db.NewDelete().Model((*AuthSessionModel)(nil)).Where("id_hash = ?", idHash).Exec(ctx)
	return err
}

func (s *Store) UpdateAuthSessionTouch(ctx context.Context, idHash string, lastIP string, lastSeenAt time.Time) (*AuthSessionModel, error) {
	_, err := s.db.NewUpdate().
		Model((*AuthSessionModel)(nil)).
		Set("last_seen_at = ?", lastSeenAt).
		Set("last_ip = ?", lastIP).
		Where("id_hash = ?", idHash).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.FetchAuthSessionByHash(ctx, idHash)
}
