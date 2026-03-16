package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/Jawadh-Salih/moodle-mcp-server/internal/api"
)

// message mirrors the API response shape from core_message_get_messages.
type message struct {
	ID               int    `json:"id"`
	UserFromFullName string `json:"userfromfullname"`
	Subject          string `json:"subject"`
	Text             string `json:"text"`
	FullMessage      string `json:"fullmessage"`
	TimeCreated      int64  `json:"timecreated"`
	TimeRead         int64  `json:"timeread"`
}

// messageDisplay is the public-facing shape for a notification.
type messageDisplay struct {
	ID      int    `json:"id"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Preview string `json:"preview"`
	Date    string `json:"date"`
	Read    bool   `json:"read"`
}

type messagesResponse struct {
	Messages []message `json:"messages"`
}

// --- Get Notifications Tool ---

type GetNotificationsInput struct {
	Limit      int  `json:"limit,omitempty"`
	UnreadOnly bool `json:"unread_only,omitempty"`
}

func HandleGetNotifications(ctx context.Context, client *api.Client, input GetNotificationsInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	readType := "unread"
	if !input.UnreadOnly {
		readType = "both"
	}

	userID := client.GetUserID()
	data, err := client.Call(ctx, "core_message_get_messages", map[string]string{
		"useridto":    fmt.Sprintf("%d", userID),
		"type":        readType,
		"limitnum":    fmt.Sprintf("%d", limit),
		"newestfirst": "1",
	})
	if err != nil {
		return "", err
	}

	var resp messagesResponse
	if err := unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("parsing messages: %w", err)
	}

	display := make([]messageDisplay, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		preview := m.Text
		if preview == "" {
			preview = m.FullMessage
		}
		preview = truncate(stripHTML(preview), 150)

		display = append(display, messageDisplay{
			ID:      m.ID,
			From:    m.UserFromFullName,
			Subject: m.Subject,
			Preview: preview,
			Date:    time.Unix(m.TimeCreated, 0).Format("2006-01-02 15:04"),
			Read:    m.TimeRead > 0,
		})
	}

	return marshalResult(map[string]any{
		"total_messages": len(display),
		"messages":       display,
	})
}
