package tools

import (
	"context"
	"fmt"

	"github.com/jawadh/moodle-mcp-server/internal/api"
	"github.com/jawadh/moodle-mcp-server/internal/config"
)

// --- Login Tool ---

type LoginInput struct {
	MoodleURL string `json:"moodle_url"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func HandleLogin(ctx context.Context, client *api.Client, input LoginInput) (string, error) {
	if input.MoodleURL == "" || input.Username == "" || input.Password == "" {
		return "", fmt.Errorf("moodle_url, username, and password are all required")
	}

	baseURL := config.NormalizeURL(input.MoodleURL)

	token, err := api.GetTokenFromCredentials(ctx, baseURL, input.Username, input.Password)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	client.SetSession(baseURL, token)

	info, err := getSiteInfoRaw(ctx, client)
	if err != nil {
		return "", fmt.Errorf("login succeeded but could not fetch site info: %w", err)
	}

	client.SetUserID(info.UserID)

	return fmt.Sprintf("Successfully logged in to %s as %s (%s). You now have access to all Moodle tools.",
		info.SiteName, info.FullName, info.Username), nil
}

// --- Get Site Info Tool ---

// SiteInfo mirrors the fields returned by core_webservice_get_site_info.
type SiteInfo struct {
	SiteName  string `json:"sitename"`
	Username  string `json:"username"`
	FullName  string `json:"fullname"`
	UserID    int    `json:"userid"`
	SiteURL   string `json:"siteurl"`
	Lang      string `json:"lang"`
	UserEmail string `json:"useremail"`
}

// getSiteInfoRaw calls core_webservice_get_site_info and returns a typed result.
// It is used internally by both HandleLogin and HandleGetSiteInfo.
func getSiteInfoRaw(ctx context.Context, client *api.Client) (*SiteInfo, error) {
	data, err := client.Call(ctx, "core_webservice_get_site_info", nil)
	if err != nil {
		return nil, err
	}
	var info SiteInfo
	if err := unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("parsing site info: %w", err)
	}
	return &info, nil
}

type GetSiteInfoInput struct{}

func HandleGetSiteInfo(ctx context.Context, client *api.Client, _ GetSiteInfoInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}
	info, err := getSiteInfoRaw(ctx, client)
	if err != nil {
		return "", err
	}
	return marshalResult(info)
}

// --- Get User Profile Tool ---

// UserProfile mirrors the fields returned by core_user_get_users_by_field.
type UserProfile struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	FullName    string `json:"fullname"`
	Email       string `json:"email"`
	Department  string `json:"department"`
	Institution string `json:"institution"`
	Description string `json:"description"`
	City        string `json:"city"`
	Country     string `json:"country"`
	ProfileURL  string `json:"profileurl"`
	ImageURL    string `json:"profileimageurl"`
}

type GetUserProfileInput struct{}

func HandleGetUserProfile(ctx context.Context, client *api.Client, _ GetUserProfileInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}

	userID := client.GetUserID()
	if userID == 0 {
		info, err := getSiteInfoRaw(ctx, client)
		if err != nil {
			return "", err
		}
		userID = info.UserID
		client.SetUserID(userID)
	}

	params := map[string]string{
		"field":     "id",
		"values[0]": fmt.Sprintf("%d", userID),
	}

	data, err := client.Call(ctx, "core_user_get_users_by_field", params)
	if err != nil {
		return "", err
	}

	var users []UserProfile
	if err := unmarshal(data, &users); err != nil {
		return "", fmt.Errorf("parsing user profile: %w", err)
	}
	if len(users) == 0 {
		return "", fmt.Errorf("user not found")
	}

	return marshalResult(users[0])
}
