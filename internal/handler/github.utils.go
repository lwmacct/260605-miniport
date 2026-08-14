package handler

import (
	"errors"

	"github.com/lwmacct/260605-miniport/internal/service"
)

func utilGithubAPIError(err error) error {
	switch {
	case errors.Is(err, service.ErrGithubDisabled):
		return utilProblem(503, "github-disabled", err.Error())
	case errors.Is(err, service.ErrGithubInvalidSignature):
		return utilProblem(403, "invalid-github-signature", err.Error())
	case errors.Is(err, service.ErrGithubInvalidState):
		return utilProblem(400, "invalid-github-state", err.Error())
	case errors.Is(err, service.ErrGithubNotFound):
		return utilProblem(404, "github-resource-not-found", err.Error())
	default:
		return utilProblem(500, "internal-server-error", "internal server error")
	}
}
