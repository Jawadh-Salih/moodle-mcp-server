package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// authClient is a dedicated HTTP client for authentication with an explicit timeout.
// This avoids using http.DefaultClient which has no timeout.
var authClient = &http.Client{Timeout: 30 * time.Second}

// GetTokenFromCredentials authenticates with Moodle and returns an API token.
// Credentials are sent in the POST body (not the URL) to avoid them appearing in server logs.
func GetTokenFromCredentials(ctx context.Context, baseURL, username, password string) (string, error) {
	endpoint := fmt.Sprintf("%s/login/token.php", baseURL)

	body := url.Values{
		"username": {username},
		"password": {password},
		"service":  {"moodle_mobile_app"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := authClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("authenticating with Moodle: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return "", fmt.Errorf("reading auth response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(data, &tokenResp); err != nil {
		return "", fmt.Errorf("parsing auth response: %w", err)
	}

	if tokenResp.Error != "" {
		return "", fmt.Errorf("%w: %s", ErrInvalidCredentials, tokenResp.Error)
	}
	if tokenResp.Token == "" {
		return "", fmt.Errorf("received empty token from Moodle")
	}

	return tokenResp.Token, nil
}
