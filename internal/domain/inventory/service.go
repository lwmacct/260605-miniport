package inventory

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Service struct {
	store *store
}

func NewService(db *bun.DB) *Service {
	return &Service{store: newStore(db)}
}

func (svc *Service) ListHosts(ctx context.Context) ([]Host, error) {
	return svc.store.listHosts(ctx)
}

func (svc *Service) CreateHost(ctx context.Context, payload HostPayload) (*Host, error) {
	host := hostFromPayload(payload)
	if err := validateHost(host); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	host.CreatedAt = now
	host.UpdatedAt = now
	if err := svc.store.createHost(ctx, host); err != nil {
		return nil, err
	}
	return host, nil
}

func (svc *Service) UpdateHost(ctx context.Context, id int64, payload HostPayload) (*Host, error) {
	host := hostFromPayload(payload)
	host.ID = id
	if err := validateHost(host); err != nil {
		return nil, err
	}

	host.UpdatedAt = time.Now().UTC()
	if err := svc.store.updateHost(ctx, id, host); err != nil {
		return nil, err
	}
	return svc.store.getHost(ctx, id)
}

func (svc *Service) DeleteHost(ctx context.Context, id int64) error {
	return svc.store.deleteHost(ctx, id)
}

func (svc *Service) ListPortGroups(ctx context.Context) ([]PortGroupView, error) {
	groups, err := svc.store.listPortGroups(ctx)
	if err != nil {
		return nil, err
	}
	return svc.store.buildPortGroupViews(ctx, groups)
}

func (svc *Service) GetPortGroup(ctx context.Context, id int64) (*PortGroupView, error) {
	group, err := svc.store.getPortGroupWithHost(ctx, id)
	if err != nil {
		return nil, err
	}
	views, err := svc.store.buildPortGroupViews(ctx, []PortGroup{*group})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (svc *Service) CreatePortGroup(ctx context.Context, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := svc.store.runInTx(ctx, func(ctx context.Context, tx *store) error {
		group, err := portGroupFromPayload(ctx, tx, 0, payload)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		group.CreatedAt = now
		group.UpdatedAt = now
		if err := tx.createPortGroup(ctx, group); err != nil {
			return err
		}
		if err := replacePortGroupChildren(ctx, tx, group.ID, payload, now); err != nil {
			return err
		}

		view, err := tx.getPortGroupView(ctx, group.ID)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func (svc *Service) UpdatePortGroup(ctx context.Context, id int64, payload PortGroupPayload) (*PortGroupView, error) {
	var out *PortGroupView
	err := svc.store.runInTx(ctx, func(ctx context.Context, tx *store) error {
		if _, err := tx.getPortGroup(ctx, id); err != nil {
			return err
		}

		group, err := portGroupFromPayload(ctx, tx, id, payload)
		if err != nil {
			return err
		}

		group.ID = id
		group.UpdatedAt = time.Now().UTC()
		if err := tx.updatePortGroup(ctx, id, group); err != nil {
			return err
		}
		if err := replacePortGroupChildren(ctx, tx, id, payload, group.UpdatedAt); err != nil {
			return err
		}

		view, err := tx.getPortGroupView(ctx, id)
		if err != nil {
			return err
		}
		out = view
		return nil
	})
	return out, err
}

func (svc *Service) DeletePortGroup(ctx context.Context, id int64) error {
	return svc.store.runInTx(ctx, func(ctx context.Context, tx *store) error {
		return tx.deletePortGroup(ctx, id)
	})
}
