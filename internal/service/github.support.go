package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/lwmacct/260605-miniport/internal/repository"
)

var (
	ErrGithubDisabled     = errors.New("github integration is disabled")
	ErrGithubUnauthorized = errors.New("github resource is not accessible")
	ErrGithubInvalidState = errors.New("github connection state is invalid or expired")
	ErrGithubNotFound     = errors.New("github resource not found")
)

const githubConnectionStateTTL = 10 * time.Minute

type GithubConfig struct {
	Enabled           bool
	AppID             int64
	AppSlug           string
	PrivateKeyFile    string
	WebhookSecret     string
	APIURL            string
	WebURL            string
	SetupReturnURL    string
	ReconcileInterval time.Duration
}

type GithubInstallation = repository.GithubInstallationRecord
type GithubRepository = repository.GithubRepositoryRecord

type GithubStatus struct {
	Enabled    bool
	AppSlug    string
	InstallURL string
}

type GithubConnectionStart struct {
	URL string
}

type githubRemoteInstallation struct {
	ID                  int64
	Account             githubRemoteOwner
	RepositorySelection string
	Permissions         map[string]string
	SuspendedAt         *time.Time
}

type githubRemoteOwner struct {
	ID        int64
	Login     string
	Type      string
	AvatarURL string
}

type githubRemoteRepository struct {
	ID            int64
	NodeID        string
	Owner         githubRemoteOwner
	Name          string
	FullName      string
	HTMLURL       string
	Description   string
	DefaultBranch string
	Visibility    string
	Private       bool
	Fork          bool
	Archived      bool
	Disabled      bool
	PushedAt      time.Time
	UpdatedAt     time.Time
}

type githubAPI interface {
	Installation(ctx context.Context, installationID int64) (githubRemoteInstallation, error)
	Repositories(ctx context.Context, installationID int64) ([]githubRemoteRepository, error)
}

func utilHashGithubState(state string) string {
	digest := sha256.Sum256([]byte(state))
	return hex.EncodeToString(digest[:])
}

func validateGithubWebhookSignature(secret, signature string, body []byte) bool {
	if secret == "" || !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	provided, err := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return hmac.Equal(provided, mac.Sum(nil))
}

func utilGithubInstallURL(webURL, slug string) string {
	return strings.TrimRight(webURL, "/") + "/apps/" + url.PathEscape(slug) + "/installations/new"
}

func utilGithubConnectionURL(installURL, state string) string {
	parsed, err := url.Parse(installURL)
	if err != nil {
		return installURL
	}
	query := parsed.Query()
	query.Set("state", state)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func utilGithubServiceError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return ErrGithubNotFound
	}
	return err
}
