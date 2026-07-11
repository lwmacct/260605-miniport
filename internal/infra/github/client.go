package github

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const githubAPIURL = "https://api.github.com"

type Config struct {
	AppID          int64
	PrivateKeyFile string
}

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type Owner struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Type      string `json:"type"`
	AvatarURL string `json:"avatar_url"`
}

type Installation struct {
	ID                  int64             `json:"id"`
	Account             Owner             `json:"account"`
	RepositorySelection string            `json:"repository_selection"`
	Permissions         map[string]string `json:"permissions"`
	SuspendedAt         *time.Time        `json:"suspended_at"`
}

type Repository struct {
	ID            int64     `json:"id"`
	NodeID        string    `json:"node_id"`
	Owner         Owner     `json:"owner"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	HTMLURL       string    `json:"html_url"`
	Description   string    `json:"description"`
	DefaultBranch string    `json:"default_branch"`
	Visibility    string    `json:"visibility"`
	Private       bool      `json:"private"`
	Fork          bool      `json:"fork"`
	Archived      bool      `json:"archived"`
	Disabled      bool      `json:"disabled"`
	PushedAt      time.Time `json:"pushed_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Client struct {
	appID      int64
	apiURL     string
	httpClient HTTPClient
	now        func() time.Time
	privateKey *rsa.PrivateKey
}

func NewClient(cfg Config, client HTTPClient) (*Client, error) {
	keyData, err := os.ReadFile(cfg.PrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("read github private key: %w", err)
	}
	privateKey, err := parsePrivateKey(keyData)
	if err != nil {
		return nil, err
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		appID: cfg.AppID, apiURL: githubAPIURL, httpClient: client,
		now: time.Now, privateKey: privateKey,
	}, nil
}

func parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("github private key is not PEM encoded")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse github private key: %w", err)
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("github private key must be RSA")
	}
	return key, nil
}

func (c *Client) Installation(ctx context.Context, installationID int64) (Installation, error) {
	var out Installation
	jwt, err := c.appJWT()
	if err != nil {
		return out, err
	}
	err = c.request(ctx, http.MethodGet, fmt.Sprintf("/app/installations/%d", installationID), jwt, nil, &out)
	return out, err
}

func (c *Client) Repositories(ctx context.Context, installationID int64) ([]Repository, error) {
	token, err := c.installationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}
	all := make([]Repository, 0)
	for page := 1; ; page++ {
		var response struct {
			Repositories []Repository `json:"repositories"`
		}
		path := fmt.Sprintf("/installation/repositories?per_page=100&page=%d", page)
		if err := c.request(ctx, http.MethodGet, path, token, nil, &response); err != nil {
			return nil, err
		}
		all = append(all, response.Repositories...)
		if len(response.Repositories) < 100 {
			return all, nil
		}
	}
}

func (c *Client) installationToken(ctx context.Context, installationID int64) (string, error) {
	jwt, err := c.appJWT()
	if err != nil {
		return "", err
	}
	var response struct {
		Token string `json:"token"`
	}
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/app/installations/%d/access_tokens", installationID), jwt, map[string]any{}, &response); err != nil {
		return "", err
	}
	if response.Token == "" {
		return "", errors.New("github returned an empty installation token")
	}
	return response.Token, nil
}

func (c *Client) request(ctx context.Context, method, path, token string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	endpoint, err := url.JoinPath(c.apiURL, strings.SplitN(path, "?", 2)[0])
	if err != nil {
		return err
	}
	if parts := strings.SplitN(path, "?", 2); len(parts) == 2 {
		endpoint += "?" + parts[1]
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Github-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("github request: %w", err)
	}
	data, readErr := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	closeErr := response.Body.Close()
	if readErr != nil {
		return readErr
	}
	if closeErr != nil {
		return fmt.Errorf("close github response: %w", closeErr)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var payload struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(data, &payload)
		if payload.Message == "" {
			payload.Message = http.StatusText(response.StatusCode)
		}
		return fmt.Errorf("github API %s: %s", response.Status, payload.Message)
	}
	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode github response: %w", err)
	}
	return nil
}

func (c *Client) appJWT() (string, error) {
	now := c.now().UTC()
	claims, err := json.Marshal(struct {
		IssuedAt  int64  `json:"iat"`
		ExpiresAt int64  `json:"exp"`
		Issuer    string `json:"iss"`
	}{
		IssuedAt: now.Add(-60 * time.Second).Unix(), ExpiresAt: now.Add(9 * time.Minute).Unix(),
		Issuer: strconv.FormatInt(c.appID, 10),
	})
	if err != nil {
		return "", fmt.Errorf("encode github app JWT claims: %w", err)
	}
	encodedHeader := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	encodedClaims := base64.RawURLEncoding.EncodeToString(claims)
	unsigned := encodedHeader + "." + encodedClaims
	digest := sha256.Sum256([]byte(unsigned))
	signature, err := rsa.SignPKCS1v15(rand.Reader, c.privateKey, crypto.SHA256, digest[:])
	if err != nil {
		return "", fmt.Errorf("sign github app JWT: %w", err)
	}
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}
