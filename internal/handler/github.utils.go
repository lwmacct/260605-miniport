package handler

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260605-miniport/internal/service"
)

func utilGithubAPIError(err error) error {
	switch {
	case errors.Is(err, service.ErrGithubDisabled):
		return huma.Error503ServiceUnavailable(err.Error())
	case errors.Is(err, service.ErrGithubInvalidSignature):
		return huma.Error403Forbidden(err.Error())
	case errors.Is(err, service.ErrGithubInvalidState):
		return huma.Error400BadRequest(err.Error())
	case errors.Is(err, service.ErrGithubNotFound):
		return huma.Error404NotFound(err.Error())
	default:
		return err
	}
}
