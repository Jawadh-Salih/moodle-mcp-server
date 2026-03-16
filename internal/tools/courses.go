package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/Jawadh-Salih/moodle-mcp-server/internal/api"
)

// CourseRef is a lightweight course identifier used to avoid N+1 API lookups.
// It is populated once from core_enrol_get_users_courses and reused across tools.
type CourseRef struct {
	ID   int
	Name string
}

// course is the raw API shape returned by core_enrol_get_users_courses.
type course struct {
	ID        int      `json:"id"`
	ShortName string   `json:"shortname"`
	FullName  string   `json:"fullname"`
	StartDate int64    `json:"startdate"`
	EndDate   int64    `json:"enddate"`
	Progress  *float64 `json:"progress,omitempty"` // float64: Moodle may return e.g. 3.2
}

// getEnrolledCourses fetches all courses the user is enrolled in.
// It is the single source of truth for course enumeration across all tools.
func getEnrolledCourses(ctx context.Context, client *api.Client) ([]CourseRef, error) {
	userID := client.GetUserID()
	if userID == 0 {
		return nil, fmt.Errorf("user ID not set — please run get_site_info first")
	}
	params := map[string]string{
		"userid": fmt.Sprintf("%d", userID),
	}
	data, err := client.Call(ctx, "core_enrol_get_users_courses", params)
	if err != nil {
		return nil, err
	}
	var courses []course
	if err := unmarshal(data, &courses); err != nil {
		return nil, fmt.Errorf("parsing courses: %w", err)
	}
	refs := make([]CourseRef, len(courses))
	for i, c := range courses {
		refs[i] = CourseRef{ID: c.ID, Name: c.FullName}
	}
	return refs, nil
}

// getEnrolledCourseIDs returns only the IDs of enrolled courses.
func getEnrolledCourseIDs(ctx context.Context, client *api.Client) ([]int, error) {
	courses, err := getEnrolledCourses(ctx, client)
	if err != nil {
		return nil, err
	}
	ids := make([]int, len(courses))
	for i, c := range courses {
		ids[i] = c.ID
	}
	return ids, nil
}

// --- List Courses Tool ---

// courseDisplay is the public-facing shape for a listed course.
type courseDisplay struct {
	ID        int    `json:"id"`
	ShortName string `json:"shortname"`
	FullName  string `json:"fullname"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date,omitempty"`
	Active    bool   `json:"active"`
}

type ListCoursesInput struct{}

func HandleListCourses(ctx context.Context, client *api.Client, _ ListCoursesInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}

	userID := client.GetUserID()
	if userID == 0 {
		return "", fmt.Errorf("user ID not set — please run get_site_info first")
	}

	params := map[string]string{
		"userid": fmt.Sprintf("%d", userID),
	}
	data, err := client.Call(ctx, "core_enrol_get_users_courses", params)
	if err != nil {
		return "", err
	}

	var courses []course
	if err := unmarshal(data, &courses); err != nil {
		return "", fmt.Errorf("parsing courses: %w", err)
	}

	now := time.Now().Unix()
	display := make([]courseDisplay, 0, len(courses))
	for _, c := range courses {
		d := courseDisplay{
			ID:        c.ID,
			ShortName: c.ShortName,
			FullName:  c.FullName,
			Active:    c.StartDate <= now && (c.EndDate == 0 || c.EndDate >= now),
		}
		if c.StartDate > 0 {
			d.StartDate = time.Unix(c.StartDate, 0).Format("2006-01-02")
		}
		if c.EndDate > 0 {
			d.EndDate = time.Unix(c.EndDate, 0).Format("2006-01-02")
		}
		display = append(display, d)
	}

	return marshalResult(map[string]any{
		"total_courses": len(display),
		"courses":       display,
	})
}

// --- Get Course Contents Tool ---

// courseSection mirrors the core_course_get_contents API response.
type courseSection struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Summary string         `json:"summary"`
	Visible int            `json:"visible"`
	Modules []courseModule `json:"modules"`
}

type courseModule struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	ModName string `json:"modname"`
	URL     string `json:"url,omitempty"`
	Visible int    `json:"visible"`
}

type GetCourseContentsInput struct {
	CourseID int `json:"course_id"`
}

func HandleGetCourseContents(ctx context.Context, client *api.Client, input GetCourseContentsInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}
	if input.CourseID == 0 {
		return "", fmt.Errorf("course_id is required")
	}

	data, err := client.Call(ctx, "core_course_get_contents", map[string]string{
		"courseid": fmt.Sprintf("%d", input.CourseID),
	})
	if err != nil {
		return "", err
	}

	var sections []courseSection
	if err := unmarshal(data, &sections); err != nil {
		return "", fmt.Errorf("parsing course contents: %w", err)
	}

	return marshalResult(map[string]any{
		"course_id":      input.CourseID,
		"total_sections": len(sections),
		"sections":       sections,
	})
}

// --- Get Course Details Tool ---

// courseDetails is the public-facing shape for a single course.
type courseDetails struct {
	ID        int    `json:"id"`
	ShortName string `json:"shortname"`
	FullName  string `json:"fullname"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date,omitempty"`
	Active    bool   `json:"active"`
}

type GetCourseDetailsInput struct {
	CourseID int `json:"course_id"`
}

// HandleGetCourseDetails returns details for a specific enrolled course.
// It uses core_enrol_get_users_courses (student-accessible) rather than
// core_course_get_courses which requires admin-level permissions.
func HandleGetCourseDetails(ctx context.Context, client *api.Client, input GetCourseDetailsInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}
	if input.CourseID == 0 {
		return "", fmt.Errorf("course_id is required")
	}

	userID := client.GetUserID()
	if userID == 0 {
		return "", fmt.Errorf("user ID not set — please run get_site_info first")
	}

	data, err := client.Call(ctx, "core_enrol_get_users_courses", map[string]string{
		"userid": fmt.Sprintf("%d", userID),
	})
	if err != nil {
		return "", err
	}

	var courses []course
	if err := unmarshal(data, &courses); err != nil {
		return "", fmt.Errorf("parsing courses: %w", err)
	}

	now := time.Now().Unix()
	for _, c := range courses {
		if c.ID != input.CourseID {
			continue
		}
		d := courseDetails{
			ID:        c.ID,
			ShortName: c.ShortName,
			FullName:  c.FullName,
			Active:    c.StartDate <= now && (c.EndDate == 0 || c.EndDate >= now),
		}
		if c.StartDate > 0 {
			d.StartDate = time.Unix(c.StartDate, 0).Format("2006-01-02")
		}
		if c.EndDate > 0 {
			d.EndDate = time.Unix(c.EndDate, 0).Format("2006-01-02")
		}
		return marshalResult(d)
	}

	return "", fmt.Errorf("course %d not found in your enrolled courses", input.CourseID)
}
