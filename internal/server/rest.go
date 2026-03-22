package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jawadh/moodle-mcp-server/internal/api"
	"github.com/jawadh/moodle-mcp-server/internal/tools"
)

// RESTServer wraps the Moodle API client and exposes it via HTTP REST API.
type RESTServer struct {
	client *api.Client
	mux    *http.ServeMux
	port   int
}

// NewRESTServer creates a new REST server.
func NewRESTServer(client *api.Client, port int) *RESTServer {
	srv := &RESTServer{
		client: client,
		mux:    http.NewServeMux(),
		port:   port,
	}
	srv.registerRoutes()
	return srv
}

// Run starts the REST server.
func (s *RESTServer) Run(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.mux,
	}
	log.Printf("REST API listening on http://localhost:%d", s.port)
	log.Println("OpenAPI docs: http://localhost:%d/api/docs", s.port)
	return httpServer.ListenAndServe()
}

func (s *RESTServer) registerRoutes() {
	// Health check
	s.mux.HandleFunc("/health", s.handleHealth)

	// OpenAPI spec
	s.mux.HandleFunc("/api/docs", s.handleOpenAPISpec)

	// Authentication
	s.mux.HandleFunc("/api/login", s.handleLogin)
	s.mux.HandleFunc("/api/site-info", s.handleSiteInfo)
	s.mux.HandleFunc("/api/user-profile", s.handleUserProfile)

	// Courses
	s.mux.HandleFunc("/api/courses", s.handleListCourses)
	s.mux.HandleFunc("/api/courses/details", s.handleCourseDetails)
	s.mux.HandleFunc("/api/courses/contents", s.handleCourseContents)

	// Grades
	s.mux.HandleFunc("/api/grades", s.handleGetGrades)
	s.mux.HandleFunc("/api/grades/overview", s.handleGradesOverview)

	// Assignments
	s.mux.HandleFunc("/api/assignments", s.handleGetAssignments)
	s.mux.HandleFunc("/api/assignments/upcoming", s.handleUpcomingAssignments)
	s.mux.HandleFunc("/api/assignments/submit", s.handleSubmitAssignment)
	s.mux.HandleFunc("/api/assignments/update", s.handleUpdateAssignment)
	s.mux.HandleFunc("/api/journal/entry", s.handleGetJournalEntry)
	s.mux.HandleFunc("/api/journal/submit", s.handleSubmitJournal)

	// Calendar
	s.mux.HandleFunc("/api/calendar/events", s.handleCalendarEvents)
	s.mux.HandleFunc("/api/calendar/deadlines", s.handleUpcomingDeadlines)

	// Messages
	s.mux.HandleFunc("/api/notifications", s.handleNotifications)
}

// --- Handlers ---

func (s *RESTServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"authenticated": s.client.IsAuthenticated(),
	})
}

func (s *RESTServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		MoodleURL string `json:"moodle_url"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	input := tools.LoginInput{
		MoodleURL: req.MoodleURL,
		Username:  req.Username,
		Password:  req.Password,
	}

	result, err := tools.HandleLogin(ctx, s.client, input)
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleSiteInfo(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleGetSiteInfo(ctx, s.client, tools.GetSiteInfoInput{})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleUserProfile(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleGetUserProfile(ctx, s.client, tools.GetUserProfileInput{})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleListCourses(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleListCourses(ctx, s.client, tools.ListCoursesInput{})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleCourseDetails(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	courseID := s.getIntQuery(r, "course_id")
	if courseID == 0 {
		http.Error(w, "Missing course_id parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleGetCourseDetails(ctx, s.client, tools.GetCourseDetailsInput{CourseID: courseID})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleCourseContents(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	courseID := s.getIntQuery(r, "course_id")
	if courseID == 0 {
		http.Error(w, "Missing course_id parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleGetCourseContents(ctx, s.client, tools.GetCourseContentsInput{CourseID: courseID})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleGetGrades(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	courseID := s.getIntQuery(r, "course_id")
	if courseID == 0 {
		http.Error(w, "Missing course_id parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleGetGrades(ctx, s.client, tools.GetGradesInput{CourseID: courseID})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleGradesOverview(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleGetGradesOverview(ctx, s.client, tools.GetGradesOverviewInput{})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleGetAssignments(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	courseID := s.getIntQuery(r, "course_id")
	if courseID == 0 {
		http.Error(w, "Missing course_id parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleGetAssignments(ctx, s.client, tools.GetAssignmentsInput{CourseID: courseID})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleUpcomingAssignments(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	daysAhead := s.getIntQuery(r, "days_ahead")
	if daysAhead == 0 {
		daysAhead = 30
	}

	ctx := r.Context()
	result, err := tools.HandleGetUpcomingAssignments(ctx, s.client, tools.GetUpcomingAssignmentsInput{DaysAhead: daysAhead})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleSubmitAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		AssignmentID int    `json:"assignment_id"`
		Text         string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleSubmitAssignment(ctx, s.client, tools.SubmitAssignmentInput{
		AssignmentID: req.AssignmentID,
		Text:         req.Text,
	})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleUpdateAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		AssignmentID int    `json:"assignment_id"`
		Text         string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	result, err := tools.HandleUpdateAssignment(ctx, s.client, tools.UpdateAssignmentInput{
		AssignmentID: req.AssignmentID,
		Text:         req.Text,
	})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleGetJournalEntry(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	journalID := s.getIntQuery(r, "journal_id")
	ctx := r.Context()
	result, err := tools.HandleGetJournalEntry(ctx, s.client, tools.GetJournalEntryInput{JournalID: journalID})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleSubmitJournal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	var req struct {
		JournalID int    `json:"journal_id"`
		Text      string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	result, err := tools.HandleSubmitJournal(ctx, s.client, tools.SubmitJournalInput{
		JournalID: req.JournalID,
		Text:      req.Text,
	})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleCalendarEvents(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	daysAhead := s.getIntQuery(r, "days_ahead")
	if daysAhead == 0 {
		daysAhead = 30
	}

	ctx := r.Context()
	result, err := tools.HandleGetCalendarEvents(ctx, s.client, tools.GetCalendarEventsInput{DaysAhead: daysAhead})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleUpcomingDeadlines(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	daysAhead := s.getIntQuery(r, "days_ahead")
	if daysAhead == 0 {
		daysAhead = 14
	}

	ctx := r.Context()
	result, err := tools.HandleGetUpcomingDeadlines(ctx, s.client, tools.GetUpcomingDeadlinesInput{DaysAhead: daysAhead})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleNotifications(w http.ResponseWriter, r *http.Request) {
	if !s.client.IsAuthenticated() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	limit := s.getIntQuery(r, "limit")
	if limit == 0 {
		limit = 20
	}

	unreadOnly := r.URL.Query().Get("unread_only") != "false"

	ctx := r.Context()
	result, err := tools.HandleGetNotifications(ctx, s.client, tools.GetNotificationsInput{
		Limit: limit, UnreadOnly: unreadOnly,
	})
	s.jsonResponse(w, result, err)
}

func (s *RESTServer) handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	spec := getOpenAPISpec()
	json.NewEncoder(w).Encode(spec)
}

// --- Utilities ---

func (s *RESTServer) jsonResponse(w http.ResponseWriter, data string, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func (s *RESTServer) getIntQuery(r *http.Request, param string) int {
	val := r.URL.Query().Get(param)
	if val == "" {
		return 0
	}
	i, _ := strconv.Atoi(val)
	return i
}

func getOpenAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Moodle API",
			"description": "REST API for Moodle LMS integration",
			"version":     "1.0.0",
		},
		"servers": []map[string]interface{}{
			{
				"url":         "http://localhost:8080",
				"description": "Local development server",
			},
		},
		"paths": map[string]interface{}{
			"/api/login": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Login to Moodle",
					"operationId": "login",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"moodle_url": map[string]interface{}{"type": "string"},
										"username":   map[string]interface{}{"type": "string"},
										"password":   map[string]interface{}{"type": "string"},
									},
									"required": []string{"moodle_url", "username", "password"},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Logged in successfully",
						},
					},
				},
			},
			"/api/courses": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "List enrolled courses",
					"operationId": "listCourses",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "List of courses",
						},
					},
				},
			},
			"/api/grades": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get grades for a course",
					"operationId": "getGrades",
					"parameters": []map[string]interface{}{
						{
							"name":        "course_id",
							"in":          "query",
							"required":    true,
							"schema":      map[string]interface{}{"type": "integer"},
							"description": "Course ID",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Grades data",
						},
					},
				},
			},
			"/api/assignments": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get assignments for a course",
					"operationId": "getAssignments",
					"parameters": []map[string]interface{}{
						{
							"name":        "course_id",
							"in":          "query",
							"required":    true,
							"schema":      map[string]interface{}{"type": "integer"},
							"description": "Course ID",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Assignments",
						},
					},
				},
			},
		},
	}
}

// truncate is a helper function
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}
