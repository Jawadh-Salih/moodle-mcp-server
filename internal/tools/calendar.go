package tools

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Jawadh-Salih/moodle-mcp-server/internal/api"
)

// calendarEvent mirrors the API response shape from core_calendar_get_calendar_events.
type calendarEvent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CourseID    int    `json:"courseid"`
	ModuleName  string `json:"modulename"`
	EventType   string `json:"eventtype"`
	TimeStart   int64  `json:"timestart"`
}

// calendarEventDisplay is the public-facing shape for a calendar event.
type calendarEventDisplay struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Course string `json:"course,omitempty"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Module string `json:"module,omitempty"`
}

type calendarResponse struct {
	Events []calendarEvent `json:"events"`
}

// --- Get Calendar Events Tool ---

type GetCalendarEventsInput struct {
	DaysAhead int `json:"days_ahead,omitempty"`
}

func HandleGetCalendarEvents(ctx context.Context, client *api.Client, input GetCalendarEventsInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}

	daysAhead := input.DaysAhead
	if daysAhead <= 0 {
		daysAhead = 30
	}

	// Fetch courses once — provides both IDs and names without extra API calls.
	courses, err := getEnrolledCourses(ctx, client)
	if err != nil {
		return "", err
	}

	courseNames := make(map[int]string, len(courses))
	for _, c := range courses {
		courseNames[c.ID] = c.Name
	}

	now := time.Now()
	params := map[string]string{
		"options[timestart]": fmt.Sprintf("%d", now.Unix()),
		"options[timeend]":   fmt.Sprintf("%d", now.Add(time.Duration(daysAhead)*24*time.Hour).Unix()),
	}
	for i, c := range courses {
		params[fmt.Sprintf("events[courseids][%d]", i)] = fmt.Sprintf("%d", c.ID)
	}

	data, err := client.Call(ctx, "core_calendar_get_calendar_events", params)
	if err != nil {
		return "", err
	}

	var resp calendarResponse
	if err := unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("parsing calendar events: %w", err)
	}

	display := make([]calendarEventDisplay, 0, len(resp.Events))
	for _, e := range resp.Events {
		t := time.Unix(e.TimeStart, 0)
		display = append(display, calendarEventDisplay{
			ID:     e.ID,
			Name:   e.Name,
			Type:   e.EventType,
			Course: courseNames[e.CourseID],
			Date:   t.Format("2006-01-02"),
			Time:   t.Format("15:04"),
			Module: e.ModuleName,
		})
	}

	sort.Slice(display, func(i, j int) bool {
		return display[i].Date < display[j].Date
	})

	return marshalResult(map[string]any{
		"days_ahead":   daysAhead,
		"total_events": len(display),
		"events":       display,
	})
}

// --- Get Upcoming Deadlines Tool ---

// deadline represents a single upcoming deadline from any source.
type deadline struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	CourseID   int    `json:"course_id"`
	CourseName string `json:"course"`
	DueDate    string `json:"due_date"`
	DueTime    string `json:"due_time"`
	DaysLeft   int    `json:"days_left"`
	Urgent     bool   `json:"urgent"`
}

type GetUpcomingDeadlinesInput struct {
	DaysAhead int `json:"days_ahead,omitempty"`
}

func HandleGetUpcomingDeadlines(ctx context.Context, client *api.Client, input GetUpcomingDeadlinesInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}

	daysAhead := input.DaysAhead
	if daysAhead <= 0 {
		daysAhead = 14
	}

	// getEnrolledCourses returns names, eliminating N+1 lookups for course names.
	courses, err := getEnrolledCourses(ctx, client)
	if err != nil {
		return "", err
	}

	courseNames := make(map[int]string, len(courses))
	for _, c := range courses {
		courseNames[c.ID] = c.Name
	}

	now := time.Now()
	cutoff := now.Add(time.Duration(daysAhead) * 24 * time.Hour)
	var deadlines []deadline

	// Collect assignment deadlines.
	for _, c := range courses {
		data, err := client.Call(ctx, "mod_assign_get_assignments", map[string]string{
			"courseids[0]": fmt.Sprintf("%d", c.ID),
		})
		if err != nil {
			continue
		}

		var resp assignmentResponse
		if err := unmarshal(data, &resp); err != nil {
			continue
		}

		for _, course := range resp.Courses {
			for _, a := range course.Assignments {
				if a.DueDate == 0 {
					continue
				}
				due := time.Unix(a.DueDate, 0)
				if due.After(now) && due.Before(cutoff) {
					daysLeft := int(time.Until(due).Hours() / 24)
					deadlines = append(deadlines, deadline{
						Name:       a.Name,
						Type:       "assignment",
						CourseID:   c.ID,
						CourseName: c.Name,
						DueDate:    due.Format("2006-01-02"),
						DueTime:    due.Format("15:04"),
						DaysLeft:   daysLeft,
						Urgent:     daysLeft <= 3,
					})
				}
			}
		}
	}

	// Collect calendar event deadlines.
	calParams := map[string]string{
		"options[timestart]": fmt.Sprintf("%d", now.Unix()),
		"options[timeend]":   fmt.Sprintf("%d", cutoff.Unix()),
	}
	for i, c := range courses {
		calParams[fmt.Sprintf("events[courseids][%d]", i)] = fmt.Sprintf("%d", c.ID)
	}
	if calData, err := client.Call(ctx, "core_calendar_get_calendar_events", calParams); err == nil {
		var calResp calendarResponse
		if unmarshal(calData, &calResp) == nil {
			for _, e := range calResp.Events {
				if e.EventType == "due" || e.EventType == "close" {
					due := time.Unix(e.TimeStart, 0)
					daysLeft := int(time.Until(due).Hours() / 24)
					deadlines = append(deadlines, deadline{
						Name:       e.Name,
						Type:       e.EventType,
						CourseID:   e.CourseID,
						CourseName: courseNames[e.CourseID],
						DueDate:    due.Format("2006-01-02"),
						DueTime:    due.Format("15:04"),
						DaysLeft:   daysLeft,
						Urgent:     daysLeft <= 3,
					})
				}
			}
		}
	}

	// Sort by urgency (fewest days left first).
	sort.Slice(deadlines, func(i, j int) bool {
		return deadlines[i].DaysLeft < deadlines[j].DaysLeft
	})

	// Deduplicate: key includes course ID to avoid dropping same-named assignments
	// from different courses that happen to share a due date.
	seen := make(map[string]bool, len(deadlines))
	unique := make([]deadline, 0, len(deadlines))
	for _, d := range deadlines {
		key := fmt.Sprintf("%d|%s|%s", d.CourseID, d.Name, d.DueDate)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, d)
		}
	}

	return marshalResult(map[string]any{
		"days_ahead":      daysAhead,
		"total_deadlines": len(unique),
		"deadlines":       unique,
	})
}
