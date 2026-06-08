package inventory

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Service struct {
	repo *repository
}

func NewService(db *bun.DB) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) ListHosts(ctx context.Context) ([]Host, error) {
	return s.repo.listHosts(ctx)
}

func (s *Service) CreateHost(ctx context.Context, payload HostPayload) (*Host, error) {
	host := hostFromPayload(payload)
	if err := validateHost(host); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	host.CreatedAt = now
	host.UpdatedAt = now
	if err := s.repo.createHost(ctx, host); err != nil {
		return nil, err
	}
	return host, nil
}

func (s *Service) UpdateHost(ctx context.Context, id int64, payload HostPayload) (*Host, error) {
	host := hostFromPayload(payload)
	host.ID = id
	if err := validateHost(host); err != nil {
		return nil, err
	}
	host.UpdatedAt = time.Now().UTC()
	if err := s.repo.updateHost(ctx, id, host); err != nil {
		return nil, err
	}
	return s.repo.getHost(ctx, id)
}

func (s *Service) DeleteHost(ctx context.Context, id int64) error {
	return s.repo.deleteHost(ctx, id)
}

func (s *Service) ListPortGroups(ctx context.Context) ([]PortGroupView, error) {
	groups, err := s.repo.listPortGroups(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.buildPortGroupViews(ctx, groups)
}

func (s *Service) GetPortGroupView(ctx context.Context, id int64) (*PortGroupView, error) {
	group, err := s.repo.getPortGroupWithHost(ctx, id)
	if err != nil {
		return nil, err
	}
	views, err := s.repo.buildPortGroupViews(ctx, []PortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *Service) CreatePortGroup(ctx context.Context, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := s.repo.runInTx(ctx, func(ctx context.Context, repo *repository) error {
		group, err := portGroupFromPayload(ctx, repo, 0, payload)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		group.CreatedAt = now
		group.UpdatedAt = now
		if err := repo.createPortGroup(ctx, group); err != nil {
			return err
		}
		if err := replacePortGroupChildren(ctx, repo, group.ID, payload, now); err != nil {
			return err
		}
		view, err := repo.getPortGroupView(ctx, group.ID)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func (s *Service) UpdatePortGroup(ctx context.Context, id int64, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := s.repo.runInTx(ctx, func(ctx context.Context, repo *repository) error {
		if _, err := repo.getPortGroup(ctx, id); err != nil {
			return err
		}
		group, err := portGroupFromPayload(ctx, repo, id, payload)
		if err != nil {
			return err
		}
		group.ID = id
		group.UpdatedAt = time.Now().UTC()
		if err := repo.updatePortGroup(ctx, id, group); err != nil {
			return err
		}
		if err := replacePortGroupChildren(ctx, repo, id, payload, group.UpdatedAt); err != nil {
			return err
		}
		view, err := repo.getPortGroupView(ctx, id)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func (s *Service) DeletePortGroup(ctx context.Context, id int64) error {
	return s.repo.runInTx(ctx, func(ctx context.Context, repo *repository) error {
		return repo.deletePortGroup(ctx, id)
	})
}
