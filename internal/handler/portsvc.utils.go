package handler

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

func utilPortsvcAPIError(err error) error {
	if err == nil {
		return nil
	}
	var portsvcErr service.PortsvcError
	if errors.As(err, &portsvcErr) {
		switch {
		case errors.Is(portsvcErr.Kind, service.ErrPortsvcNotFound):
			return huma.Error404NotFound(portsvcErr.Message)
		case errors.Is(portsvcErr.Kind, service.ErrPortsvcBadRequest):
			return huma.Error400BadRequest(portsvcErr.Message)
		}
	}
	if errors.Is(err, service.ErrPortsvcNotFound) {
		return huma.Error404NotFound("not found")
	}
	if errors.Is(err, service.ErrPortsvcBadRequest) {
		return huma.Error400BadRequest("bad request")
	}
	return err
}
