package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lwmacct/260630-go-hsr-shared/pkg/idgen"
	"github.com/uptrace/bun"
)

func (s *Store) UpsertGithubInstallation(ctx context.Context, item GithubInstallationRecord) (*GithubInstallationRecord, error) {
	row := new(GithubInstallationsModel)
	err := s.db.NewSelect().Model(row).Where("github_installation_id = ?", item.GithubInstallationID).Scan(ctx)
	if err != nil && !errors.Is(WrapNotFound(err), ErrNotFound) {
		return nil, err
	}
	if err != nil {
		row.ID = idgen.NewUUID7()
		row.CreatedAt = item.CreatedAt
		if row.CreatedAt.IsZero() {
			row.CreatedAt = item.UpdatedAt
		}
	}
	row.GithubInstallationID = item.GithubInstallationID
	row.AccountID = item.AccountID
	row.AccountLogin = item.AccountLogin
	row.AccountType = item.AccountType
	row.AvatarURL = item.AvatarURL
	row.RepositorySelection = item.RepositorySelection
	row.Permissions = item.Permissions
	row.Status = item.Status
	row.SuspendedAt = item.SuspendedAt
	row.UpdatedAt = item.UpdatedAt
	if err != nil {
		if _, insertErr := s.db.NewInsert().Model(row).Exec(ctx); insertErr != nil {
			return nil, insertErr
		}
	} else if _, updateErr := s.db.NewUpdate().Model(row).WherePK().Exec(ctx); updateErr != nil {
		return nil, updateErr
	}
	return utilGithubInstallationRecord(row), nil
}

func (s *Store) FetchGithubInstallationByID(ctx context.Context, id string) (*GithubInstallationRecord, error) {
	row := new(GithubInstallationsModel)
	if err := s.db.NewSelect().Model(row).Where("github_installation.id = ?", id).Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	return utilGithubInstallationRecord(row), nil
}

func (s *Store) FetchGithubInstallationByExternalID(ctx context.Context, id int64) (*GithubInstallationRecord, error) {
	row := new(GithubInstallationsModel)
	if err := s.db.NewSelect().Model(row).Where("github_installation_id = ?", id).Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	return utilGithubInstallationRecord(row), nil
}

func (s *Store) ListGithubInstallations(ctx context.Context) ([]GithubInstallationRecord, error) {
	var rows []GithubInstallationsModel
	err := s.db.NewSelect().Model(&rows).
		Order("github_installation.account_login ASC", "github_installation.id ASC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]GithubInstallationRecord, 0, len(rows))
	for idx := range rows {
		out = append(out, *utilGithubInstallationRecord(&rows[idx]))
	}
	return out, nil
}

func (s *Store) ListGithubInstallationsByActiveStatus(ctx context.Context) ([]GithubInstallationRecord, error) {
	var rows []GithubInstallationsModel
	if err := s.db.NewSelect().Model(&rows).Where("status = 'active'").Order("id ASC").Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]GithubInstallationRecord, 0, len(rows))
	for idx := range rows {
		out = append(out, *utilGithubInstallationRecord(&rows[idx]))
	}
	return out, nil
}

func (s *Store) UpdateGithubInstallationSync(ctx context.Context, id string, syncedAt time.Time, syncError string) (*GithubInstallationRecord, error) {
	_, err := s.db.NewUpdate().Model((*GithubInstallationsModel)(nil)).
		Set("last_synced_at = ?", syncedAt).Set("last_sync_error = ?", syncError).Set("updated_at = ?", syncedAt).
		Where("id = ?", id).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.FetchGithubInstallationByID(ctx, id)
}

func (s *Store) UpdateGithubInstallationStatus(ctx context.Context, externalID int64, status string, suspendedAt, now time.Time) (*GithubInstallationRecord, error) {
	_, err := s.db.NewUpdate().Model((*GithubInstallationsModel)(nil)).
		Set("status = ?", status).Set("suspended_at = ?", suspendedAt).Set("updated_at = ?", now).
		Where("github_installation_id = ?", externalID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.FetchGithubInstallationByExternalID(ctx, externalID)
}

func (s *Store) ReplaceGithubRepositoriesStateWithUnavailable(ctx context.Context, installationID string, now time.Time) error {
	_, err := s.db.NewUpdate().Model((*GithubRepositoriesModel)(nil)).
		Set("state = 'unavailable'").Set("updated_at = ?", now).
		Where("installation_id = ?", installationID).Exec(ctx)
	return err
}

func (s *Store) AddGithubConnectionState(ctx context.Context, stateHash string, expiresAt, now time.Time) error {
	_, err := s.db.NewDelete().Model((*GithubConnectionStatesModel)(nil)).Where("expires_at <= ?", now).Exec(ctx)
	if err != nil {
		return err
	}
	row := &GithubConnectionStatesModel{StateHash: stateHash, ExpiresAt: expiresAt, CreatedAt: now}
	_, err = s.db.NewInsert().Model(row).Exec(ctx)
	return err
}

func (s *Store) FetchGithubConnectionStateByHash(ctx context.Context, stateHash string, now time.Time) (*GithubConnectionStateRecord, error) {
	row := new(GithubConnectionStatesModel)
	if err := s.db.NewSelect().Model(row).Where("state_hash = ?", stateHash).
		Where("expires_at > ?", now).Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	return utilGithubConnectionStateRecord(row), nil
}

func (s *Store) DeleteGithubConnectionState(ctx context.Context, stateHash string, now time.Time) (bool, error) {
	result, err := s.db.NewDelete().Model((*GithubConnectionStatesModel)(nil)).
		Where("state_hash = ?", stateHash).Where("expires_at > ?", now).Exec(ctx)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}

func (s *Store) ReplaceGithubRepositories(ctx context.Context, installationID string, items []GithubRepositoryRecord, now time.Time) error {
	return s.RunInTx(ctx, func(ctx context.Context, tx *Store) error {
		ids := make([]int64, 0, len(items))
		for _, item := range items {
			ids = append(ids, item.GithubRepositoryID)
			row := &GithubRepositoriesModel{
				ID: idgen.NewUUID7(), InstallationID: installationID, GithubRepositoryID: item.GithubRepositoryID,
				NodeID: item.NodeID, OwnerLogin: item.OwnerLogin, Name: item.Name, FullName: item.FullName,
				HTMLURL: item.HTMLURL, Description: item.Description, DefaultBranch: item.DefaultBranch,
				Visibility: item.Visibility, Private: item.Private, Fork: item.Fork, Archived: item.Archived,
				Disabled: item.Disabled, State: "active", PushedAt: item.PushedAt, RemoteUpdatedAt: item.RemoteUpdatedAt,
				LastSeenAt: now, CreatedAt: now, UpdatedAt: now,
			}
			_, err := tx.db.NewInsert().Model(row).
				On("CONFLICT (github_repository_id) DO UPDATE").
				Set("installation_id = EXCLUDED.installation_id").Set("node_id = EXCLUDED.node_id").
				Set("owner_login = EXCLUDED.owner_login").Set("name = EXCLUDED.name").Set("full_name = EXCLUDED.full_name").
				Set("html_url = EXCLUDED.html_url").Set("description = EXCLUDED.description").
				Set("default_branch = EXCLUDED.default_branch").Set("visibility = EXCLUDED.visibility").
				Set("private = EXCLUDED.private").Set("fork = EXCLUDED.fork").Set("archived = EXCLUDED.archived").
				Set("disabled = EXCLUDED.disabled").Set("state = 'active'").Set("pushed_at = EXCLUDED.pushed_at").
				Set("remote_updated_at = EXCLUDED.remote_updated_at").Set("last_seen_at = EXCLUDED.last_seen_at").
				Set("updated_at = EXCLUDED.updated_at").Exec(ctx)
			if err != nil {
				return err
			}
		}

		query := tx.db.NewUpdate().Model((*GithubRepositoriesModel)(nil)).
			Set("state = 'unavailable'").Set("updated_at = ?", now).Where("installation_id = ?", installationID)
		if len(ids) > 0 {
			query = query.Where("github_repository_id NOT IN (?)", bun.List(ids))
		}
		_, err := query.Exec(ctx)
		return err
	})
}

func (s *Store) ListGithubRepositories(ctx context.Context, queryText, state string) ([]GithubRepositoryRecord, error) {
	var rows []GithubRepositoriesModel
	query := s.db.NewSelect().Model(&rows)
	if state != "" {
		query = query.Where("github_repository.state = ?", state)
	}
	if keyword := utilCompactString(queryText); keyword != "" {
		pattern := utilSearchPattern(keyword)
		query = query.Where(utilJoinSearchClauses([]string{"github_repository.full_name", "github_repository.description"}), pattern, pattern)
	}
	if err := query.Order("github_repository.full_name ASC").Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]GithubRepositoryRecord, 0, len(rows))
	for idx := range rows {
		out = append(out, *utilGithubRepositoryRecord(&rows[idx]))
	}
	return out, nil
}

func (s *Store) FetchGithubRepositoryByID(ctx context.Context, id string) (*GithubRepositoryRecord, error) {
	row := new(GithubRepositoriesModel)
	if err := s.db.NewSelect().Model(row).Where("github_repository.id = ?", id).Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	return utilGithubRepositoryRecord(row), nil
}

func (s *Store) ExistsGithubRepository(ctx context.Context, repositoryID string) (bool, error) {
	count, err := s.db.NewSelect().Model((*GithubRepositoriesModel)(nil)).
		Where("github_repository.id = ?", repositoryID).Count(ctx)
	return count > 0, err
}

func (s *Store) AddGithubWebhookDelivery(ctx context.Context, item GithubWebhookDeliveryRecord) error {
	existing := new(GithubWebhookDeliveriesModel)
	err := s.db.NewSelect().Model(existing).Where("delivery_id = ?", item.DeliveryID).Scan(ctx)
	if err == nil {
		if existing.Status != "failed" {
			return ErrGithubWebhookDeliveryExists
		}
		_, err = s.db.NewUpdate().Model(existing).Set("status = 'processing'").Set("error = ''").Set("processed_at = NULL").WherePK().Exec(ctx)
		return err
	}
	if !errors.Is(WrapNotFound(err), ErrNotFound) {
		return err
	}
	row := &GithubWebhookDeliveriesModel{
		DeliveryID: item.DeliveryID, Event: item.Event, Action: item.Action,
		InstallationID: item.InstallationID, Status: item.Status, ReceivedAt: item.ReceivedAt,
	}
	result, err := s.db.NewInsert().Model(row).On("CONFLICT (delivery_id) DO NOTHING").Exec(ctx)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrGithubWebhookDeliveryExists
	}
	return nil
}

func (s *Store) ReplaceGithubWebhookDeliveryResult(ctx context.Context, deliveryID, status, deliveryError string, now time.Time) error {
	_, err := s.db.NewUpdate().Model((*GithubWebhookDeliveriesModel)(nil)).
		Set("status = ?", status).Set("error = ?", deliveryError).Set("processed_at = ?", now).
		Where("delivery_id = ?", deliveryID).Exec(ctx)
	return err
}
