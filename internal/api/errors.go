package api

import (
	"fmt"
)

// APIError represents an error returned by the Moodle REST API.
type APIError struct {
	ErrorCode string `json:"errorcode"`
	Message   string `json:"message"`
	Exception string `json:"exception"`
}

func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("moodle api error [%s]: %s", e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("moodle api error: %s", e.Message)
}

// ErrNotAuthenticated is returned when no active session exists.
var ErrNotAuthenticated = fmt.Errorf("not authenticated — please use the 'login' tool first with your Moodle URL, username, and password")

// ErrInvalidCredentials is returned when login credentials are wrong.
var ErrInvalidCredentials = fmt.Errorf("invalid username or password")

// TokenResponse is the response from Moodle's login/token.php endpoint.
type TokenResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}
