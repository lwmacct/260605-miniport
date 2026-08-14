package handler

import (
	"errors"
	"time"

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
			return utilProblem(404, "portsvc-resource-not-found", portsvcErr.Message)
		case errors.Is(portsvcErr.Kind, service.ErrPortsvcBadRequest):
			return utilProblem(400, "invalid-portsvc-request", portsvcErr.Message)
		case errors.Is(portsvcErr.Kind, service.ErrPortsvcConflict):
			return utilProblem(409, "portsvc-resource-conflict", portsvcErr.Message)
		}
	}
	if errors.Is(err, service.ErrPortsvcNotFound) {
		return utilProblem(404, "portsvc-resource-not-found", "resource not found")
	}
	if errors.Is(err, service.ErrPortsvcBadRequest) {
		return utilProblem(400, "invalid-portsvc-request", "invalid request")
	}
	if errors.Is(err, service.ErrPortsvcConflict) {
		return utilProblem(409, "portsvc-resource-conflict", "resource conflict")
	}
	return utilProblem(500, "internal-server-error", "internal server error")
}

func utilHTTPTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
