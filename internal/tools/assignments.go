package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/jawadh/moodle-mcp-server/internal/api"
)

// assignment mirrors the API response shape from mod_assign_get_assignments.
type assignment struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	DueDate    int64  `json:"duedate"`
	CutoffDate int64  `json:"cutoffdate"`
	Grade      int    `json:"grade"`
	Intro      string `json:"intro"`
}

// assignmentDisplay is the public-facing shape for a single assignment.
type assignmentDisplay struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	DueDate    string `json:"due_date"`
	CutoffDate string `json:"cutoff_date,omitempty"`
	MaxGrade   int    `json:"max_grade"`
	Overdue    bool   `json:"overdue"`
	DaysLeft   int    `json:"days_left,omitempty"`
	Intro      string `json:"description,omitempty"`
}

type assignmentCourse struct {
	ID          int          `json:"id"`
	FullName    string       `json:"fullname"`
	Assignments []assignment `json:"assignments"`
}

type assignmentResponse struct {
	Courses []assignmentCourse `json:"courses"`
}

// --- Get Assignments Tool ---

type GetAssignmentsInput struct {
	CourseID int `json:"course_id"`
}

func HandleGetAssignments(ctx context.Context, client *api.Client, input GetAssignmentsInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}
	if input.CourseID == 0 {
		return "", fmt.Errorf("course_id is required")
	}

	data, err := client.Call(ctx, "mod_assign_get_assignments", map[string]string{
		"courseids[0]": fmt.Sprintf("%d", input.CourseID),
	})
	if err != nil {
		return "", err
	}

	var resp assignmentResponse
	if err := unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("parsing assignments: %w", err)
	}

	now := time.Now()
	display := make([]assignmentDisplay, 0)
	for _, c := range resp.Courses {
		for _, a := range c.Assignments {
			display = append(display, buildAssignmentDisplay(a, now))
		}
	}

	return marshalResult(map[string]any{
		"course_id":         input.CourseID,
		"total_assignments": len(display),
		"assignments":       display,
	})
}

// --- Get Upcoming Assignments Tool ---

// upcomingAssignment enriches assignmentDisplay with its course context.
type upcomingAssignment struct {
	CourseID   int    `json:"course_id"`
	CourseName string `json:"course_name"`
	assignmentDisplay
}

type GetUpcomingAssignmentsInput struct {
	DaysAhead int `json:"days_ahead,omitempty"`
}

func HandleGetUpcomingAssignments(ctx context.Context, client *api.Client, input GetUpcomingAssignmentsInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}

	daysAhead := input.DaysAhead
	if daysAhead <= 0 {
		daysAhead = 30
	}

	courseIDs, err := getEnrolledCourseIDs(ctx, client)
	if err != nil {
		return "", err
	}

	now := time.Now()
	cutoff := now.Add(time.Duration(daysAhead) * 24 * time.Hour)
	upcoming := make([]upcomingAssignment, 0)

	for _, cid := range courseIDs {
		data, err := client.Call(ctx, "mod_assign_get_assignments", map[string]string{
			"courseids[0]": fmt.Sprintf("%d", cid),
		})
		if err != nil {
			continue
		}

		var resp assignmentResponse
		if err := unmarshal(data, &resp); err != nil {
			continue
		}

		for _, c := range resp.Courses {
			for _, a := range c.Assignments {
				if a.DueDate == 0 {
					continue
				}
				due := time.Unix(a.DueDate, 0)
				if due.After(now) && due.Before(cutoff) {
					upcoming = append(upcoming, upcomingAssignment{
						CourseID:          cid,
						CourseName:        c.FullName,
						assignmentDisplay: buildAssignmentDisplay(a, now),
					})
				}
			}
		}
	}

	return marshalResult(map[string]any{
		"days_ahead":           daysAhead,
		"upcoming_count":       len(upcoming),
		"upcoming_assignments": upcoming,
	})
}

// --- Submit Assignment Tool ---

type SubmitAssignmentInput struct {
	AssignmentID int    `json:"assignment_id"`
	Text         string `json:"text"`
}

func HandleSubmitAssignment(ctx context.Context, client *api.Client, input SubmitAssignmentInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}
	if input.AssignmentID == 0 {
		return "", fmt.Errorf("assignment_id is required")
	}
	if input.Text == "" {
		return "", fmt.Errorf("text content is required for submission")
	}

	_, err := client.Call(ctx, "mod_assign_save_submission", map[string]string{
		"assignmentid":                         fmt.Sprintf("%d", input.AssignmentID),
		"plugindata[onlinetext_editor][text]":   input.Text,
		"plugindata[onlinetext_editor][format]": "1",
		"plugindata[onlinetext_editor][itemid]": "0",
	})
	if err != nil {
		return "", fmt.Errorf("submitting assignment: %w", err)
	}

	return fmt.Sprintf("Assignment %d submitted successfully at %s.",
		input.AssignmentID, time.Now().Format("2006-01-02 15:04:05")), nil
}

// buildAssignmentDisplay converts a raw assignment into its display form.
func buildAssignmentDisplay(a assignment, now time.Time) assignmentDisplay {
	d := assignmentDisplay{
		ID:       a.ID,
		Name:     a.Name,
		MaxGrade: a.Grade,
		Intro:    truncate(a.Intro, 200),
	}
	if a.DueDate > 0 {
		due := time.Unix(a.DueDate, 0)
		d.DueDate = due.Format("2006-01-02 15:04")
		d.Overdue = now.After(due)
		if !d.Overdue {
			d.DaysLeft = int(time.Until(due).Hours() / 24)
		}
	}
	if a.CutoffDate > 0 {
		d.CutoffDate = time.Unix(a.CutoffDate, 0).Format("2006-01-02 15:04")
	}
	return d
}
