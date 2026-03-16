package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Jawadh-Salih/moodle-mcp-server/internal/api"
	"github.com/Jawadh-Salih/moodle-mcp-server/internal/config"
	"github.com/Jawadh-Salih/moodle-mcp-server/internal/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load config from environment (all optional — login tool can provide at runtime)
	cfg := config.LoadFromEnv()
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	// Create the Moodle API client
	client := api.NewClient()

	// If credentials are provided via env vars, authenticate now
	ctx := context.Background()
	if cfg.HasAuth() {
		token := cfg.Token
		if token == "" {
			var err error
			token, err = api.GetTokenFromCredentials(ctx, cfg.MoodleURL, cfg.Username, cfg.Password)
			if err != nil {
				fmt.Fprintf(os.Stderr, "authentication failed: %v\n", err)
				os.Exit(1)
			}
		}
		client.SetSession(cfg.MoodleURL, token)

		// Get user ID from site info
		if data, err := client.Call(ctx, "core_webservice_get_site_info", nil); err == nil {
			var info struct {
				UserID int `json:"userid"`
			}
			if json.Unmarshal(data, &info) == nil {
				client.SetUserID(info.UserID)
			}
		}
	}

	// Create MCP server
	s := server.NewMCPServer(
		"moodle-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register all tools
	registerTools(s, client)

	// Run the server over stdio
	log.Println("Moodle MCP Server starting...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func registerTools(s *server.MCPServer, client *api.Client) {

	// ── Login ──────────────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("login",
			mcp.WithDescription("Log in to your Moodle account. This is the first step — provide your Moodle site URL, username, and password."),
			mcp.WithString("moodle_url", mcp.Required(), mcp.Description("Your Moodle site URL (e.g. https://online.uom.lk)")),
			mcp.WithString("username", mcp.Required(), mcp.Description("Your Moodle username or email")),
			mcp.WithString("password", mcp.Required(), mcp.Description("Your Moodle password")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			url := mcp.ParseString(req, "moodle_url", "")
			user := mcp.ParseString(req, "username", "")
			pass := mcp.ParseString(req, "password", "")
			result, err := tools.HandleLogin(ctx, client, tools.LoginInput{
				MoodleURL: url, Username: user, Password: pass,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Site Info ──────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_site_info",
			mcp.WithDescription("Get information about the connected Moodle site and the logged-in user."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := tools.HandleGetSiteInfo(ctx, client, tools.GetSiteInfoInput{})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get User Profile ──────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_user_profile",
			mcp.WithDescription("Get the detailed profile of the currently logged-in user."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := tools.HandleGetUserProfile(ctx, client, tools.GetUserProfileInput{})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── List Courses ──────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("list_courses",
			mcp.WithDescription("List all courses the student is enrolled in."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := tools.HandleListCourses(ctx, client, tools.ListCoursesInput{})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Course Contents ───────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_course_contents",
			mcp.WithDescription("Get sections, resources, and activities of a specific course."),
			mcp.WithNumber("course_id", mcp.Required(), mcp.Description("The Moodle course ID (get it from list_courses)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := mcp.ParseInt(req, "course_id", 0)
			result, err := tools.HandleGetCourseContents(ctx, client, tools.GetCourseContentsInput{CourseID: id})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Course Details ────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_course_details",
			mcp.WithDescription("Get detailed metadata about a specific course (description, dates, format)."),
			mcp.WithNumber("course_id", mcp.Required(), mcp.Description("The Moodle course ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := mcp.ParseInt(req, "course_id", 0)
			result, err := tools.HandleGetCourseDetails(ctx, client, tools.GetCourseDetailsInput{CourseID: id})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Grades ────────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_grades",
			mcp.WithDescription("Get all grade items and scores for a specific course."),
			mcp.WithNumber("course_id", mcp.Required(), mcp.Description("The Moodle course ID to get grades for")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := mcp.ParseInt(req, "course_id", 0)
			result, err := tools.HandleGetGrades(ctx, client, tools.GetGradesInput{CourseID: id})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Grades Overview ───────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_grades_overview",
			mcp.WithDescription("Get a summary of grades across all enrolled courses."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := tools.HandleGetGradesOverview(ctx, client, tools.GetGradesOverviewInput{})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Assignments ───────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_assignments",
			mcp.WithDescription("Get all assignments for a specific course with due dates and status."),
			mcp.WithNumber("course_id", mcp.Required(), mcp.Description("The Moodle course ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := mcp.ParseInt(req, "course_id", 0)
			result, err := tools.HandleGetAssignments(ctx, client, tools.GetAssignmentsInput{CourseID: id})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Upcoming Assignments ──────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_upcoming_assignments",
			mcp.WithDescription("Get all upcoming assignments across all enrolled courses, sorted by due date."),
			mcp.WithNumber("days_ahead", mcp.Description("Number of days to look ahead (default: 30)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			days := mcp.ParseInt(req, "days_ahead", 0)
			result, err := tools.HandleGetUpcomingAssignments(ctx, client, tools.GetUpcomingAssignmentsInput{DaysAhead: days})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Submit Assignment ─────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("submit_assignment",
			mcp.WithDescription("Submit text content for an online text assignment."),
			mcp.WithNumber("assignment_id", mcp.Required(), mcp.Description("The assignment ID to submit to")),
			mcp.WithString("text", mcp.Required(), mcp.Description("The text content to submit")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := mcp.ParseInt(req, "assignment_id", 0)
			text := mcp.ParseString(req, "text", "")
			result, err := tools.HandleSubmitAssignment(ctx, client, tools.SubmitAssignmentInput{
				AssignmentID: id, Text: text,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Calendar Events ───────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_calendar_events",
			mcp.WithDescription("Get upcoming calendar events from all enrolled courses."),
			mcp.WithNumber("days_ahead", mcp.Description("Number of days to look ahead (default: 30)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			days := mcp.ParseInt(req, "days_ahead", 0)
			result, err := tools.HandleGetCalendarEvents(ctx, client, tools.GetCalendarEventsInput{DaysAhead: days})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Upcoming Deadlines ────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_upcoming_deadlines",
			mcp.WithDescription("Get a consolidated view of all upcoming deadlines (assignments, quizzes, etc.), sorted by urgency."),
			mcp.WithNumber("days_ahead", mcp.Description("Number of days to look ahead (default: 14)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			days := mcp.ParseInt(req, "days_ahead", 0)
			result, err := tools.HandleGetUpcomingDeadlines(ctx, client, tools.GetUpcomingDeadlinesInput{DaysAhead: days})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	// ── Get Notifications ─────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_notifications",
			mcp.WithDescription("Get recent messages and notifications from Moodle."),
			mcp.WithNumber("limit", mcp.Description("Maximum number of notifications (default: 20)")),
			mcp.WithBoolean("unread_only", mcp.Description("Only show unread notifications (default: true)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			limit := mcp.ParseInt(req, "limit", 0)
			unreadOnly := mcp.ParseBoolean(req, "unread_only", true)
			result, err := tools.HandleGetNotifications(ctx, client, tools.GetNotificationsInput{
				Limit: limit, UnreadOnly: unreadOnly,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(result), nil
		},
	)
}
