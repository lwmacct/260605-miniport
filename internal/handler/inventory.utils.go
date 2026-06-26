package handler

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

func utilInventoryAPIError(err error) error {
	if err == nil {
		return nil
	}
	var invErr service.InventoryError
	if errors.As(err, &invErr) {
		switch {
		case errors.Is(invErr.Kind, service.ErrInventoryNotFound):
			return huma.Error404NotFound(invErr.Message)
		case errors.Is(invErr.Kind, service.ErrInventoryBadRequest):
			return huma.Error400BadRequest(invErr.Message)
		}
	}
	if errors.Is(err, service.ErrInventoryNotFound) {
		return huma.Error404NotFound("not found")
	}
	if errors.Is(err, service.ErrInventoryBadRequest) {
		return huma.Error400BadRequest("bad request")
	}
	return err
}
