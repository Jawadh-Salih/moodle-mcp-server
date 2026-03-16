package tools

import (
	"context"
	"fmt"

	"github.com/jawadh/moodle-mcp-server/internal/api"
)

// --- Get Grades Tool ---

// GradeItem mirrors a single grade entry from gradereport_user_get_grade_items.
type GradeItem struct {
	ID             int     `json:"id"`
	ItemName       string  `json:"itemname"`
	ItemType       string  `json:"itemtype"`
	ItemModule     string  `json:"itemmodule"`
	GradeFormatted string  `json:"gradeformatted"`
	GradeMin       float64 `json:"grademin"`
	GradeMax       float64 `json:"grademax"`
	Feedback       string  `json:"feedback,omitempty"`
	Percentage     string  `json:"percentageformatted,omitempty"`
}

type userGrade struct {
	CourseID   int         `json:"courseid"`
	UserID     int         `json:"userid"`
	FullName   string      `json:"userfullname"`
	GradeItems []GradeItem `json:"gradeitems"`
}

type gradeReport struct {
	UserGrades []userGrade `json:"usergrades"`
}

type GetGradesInput struct {
	CourseID int `json:"course_id"`
}

func HandleGetGrades(ctx context.Context, client *api.Client, input GetGradesInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}
	if input.CourseID == 0 {
		return "", fmt.Errorf("course_id is required")
	}

	userID := client.GetUserID()
	data, err := client.Call(ctx, "gradereport_user_get_grade_items", map[string]string{
		"courseid": fmt.Sprintf("%d", input.CourseID),
		"userid":   fmt.Sprintf("%d", userID),
	})
	if err != nil {
		return "", err
	}

	var report gradeReport
	if err := unmarshal(data, &report); err != nil {
		return "", fmt.Errorf("parsing grades: %w", err)
	}

	return marshalResult(report)
}

// --- Get Grades Overview Tool ---

// CourseGradeSummary is a condensed grade summary for a single course.
type CourseGradeSummary struct {
	CourseID   int    `json:"course_id"`
	CourseName string `json:"course_name"`
	TotalGrade string `json:"total_grade"`
	GradeMax   string `json:"grade_max"`
	ItemCount  int    `json:"graded_items"`
}

type GetGradesOverviewInput struct{}

func HandleGetGradesOverview(ctx context.Context, client *api.Client, _ GetGradesOverviewInput) (string, error) {
	if !client.IsAuthenticated() {
		return "", api.ErrNotAuthenticated
	}

	// getEnrolledCourses returns names alongside IDs, eliminating N+1 lookups.
	courses, err := getEnrolledCourses(ctx, client)
	if err != nil {
		return "", fmt.Errorf("getting enrolled courses: %w", err)
	}

	userID := client.GetUserID()
	summaries := make([]CourseGradeSummary, 0, len(courses))

	for _, course := range courses {
		data, err := client.Call(ctx, "gradereport_user_get_grade_items", map[string]string{
			"courseid": fmt.Sprintf("%d", course.ID),
			"userid":   fmt.Sprintf("%d", userID),
		})
		if err != nil {
			continue
		}

		var report gradeReport
		if err := unmarshal(data, &report); err != nil || len(report.UserGrades) == 0 {
			continue
		}

		ug := report.UserGrades[0]
		var totalGrade, gradeMax string
		gradedCount := 0
		for _, item := range ug.GradeItems {
			if item.ItemType == "course" {
				totalGrade = item.GradeFormatted
				gradeMax = fmt.Sprintf("%.0f", item.GradeMax)
			}
			if item.ItemType == "mod" && item.GradeFormatted != "-" && item.GradeFormatted != "" {
				gradedCount++
			}
		}

		summaries = append(summaries, CourseGradeSummary{
			CourseID:   course.ID,
			CourseName: course.Name,
			TotalGrade: totalGrade,
			GradeMax:   gradeMax,
			ItemCount:  gradedCount,
		})
	}

	return marshalResult(map[string]any{
		"total_courses": len(summaries),
		"grades":        summaries,
	})
}
