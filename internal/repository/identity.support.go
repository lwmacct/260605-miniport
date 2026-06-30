package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type IdentityUserRecord struct {
	ID         int64
	Username   string
	Status     string
	DisabledAt *time.Time
}

func (s *Store) FetchIdentityUser(ctx context.Context, id int64) (*IdentityUserRecord, error) {
	var row IdentityUserRecord
	err := s.db.NewSelect().
		TableExpr("users").
		Column("id", "username", "status", "disabled_at").
		Where("id = ?", id).
		Scan(ctx, &row)
	return &row, WrapNotFound(err)
}

func (s *Store) identityUsernames(ctx context.Context, userIDs []int64) (map[int64]string, error) {
	ids := utilUniqueInt64s(userIDs)
	if len(ids) == 0 {
		return map[int64]string{}, nil
	}
	var rows []struct {
		ID       int64  `bun:"id"`
		Username string `bun:"username"`
	}
	if err := s.db.NewSelect().
		TableExpr("users").
		Column("id", "username").
		Where("id IN (?)", bun.List(ids)).
		Scan(ctx, &rows); err != nil {
		return nil, err
	}
	out := make(map[int64]string, len(rows))
	for _, row := range rows {
		out[row.ID] = row.Username
	}
	return out, nil
}

func utilUniqueInt64s(values []int64) []int64 {
	seen := make(map[int64]struct{}, len(values))
	out := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
