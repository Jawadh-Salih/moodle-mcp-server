# ChatGPT Setup Guide

Use Moodle MCP Server with ChatGPT via the REST API and Custom GPT Actions.

## Step 1: Start the REST API Server

### On Your Mac/Linux

```bash
cd ~/moodle-mcp-server
./moodle-mcp -mode rest -port 8080
```

### On Windows (PowerShell)

```powershell
cd C:\Users\YourName\moodle-mcp-server
.\moodle-mcp.exe -mode rest -port 8080
```

You should see:
```
REST API listening on http://localhost:8080
OpenAPI docs at http://localhost:8080/api/docs
```

## Step 2: Expose Your Server to the Internet

For ChatGPT to reach it, the server must be publicly accessible.

### Option A: ngrok (Easiest for Testing)

```bash
brew install ngrok  # macOS
ngrok http 8080
```

You'll get a URL like: `https://abc123.ngrok.io`

### Option B: Deploy to Cloud (Recommended for Production)

See `DEPLOYMENT_GUIDE.md` for Docker and cloud hosting options.

## Step 3: Create a Custom GPT

1. Go to [ChatGPT Custom GPT Builder](https://chat.openai.com/gpts/editor)
2. Click **Create New GPT** → **Configure** tab
3. Name it: **Moodle Assistant**
4. Scroll down to **Actions** → **Create new action**

## Step 4: Paste the OpenAPI Schema

Replace `YOUR_URL` with your ngrok URL or domain:

```yaml
openapi: 3.0.0
info:
  title: Moodle API
  version: 1.2.0
servers:
  - url: https://YOUR_URL

paths:
  /api/login:
    post:
      operationId: login
      summary: Login to Moodle
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [moodle_url, username, password]
              properties:
                moodle_url:
                  type: string
                  example: "https://online.uom.lk"
                username:
                  type: string
                password:
                  type: string
      responses:
        '200':
          description: Logged in successfully

  /api/courses:
    get:
      operationId: listCourses
      summary: List all enrolled courses
      responses:
        '200':
          description: List of courses

  /api/courses/contents:
    get:
      operationId: getCourseContents
      summary: Get course sections, resources, and activities (use this to find journal IDs)
      parameters:
        - name: course_id
          in: query
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Course contents

  /api/grades/overview:
    get:
      operationId: getGradesOverview
      summary: Get grade summary across all courses
      responses:
        '200':
          description: Grades summary

  /api/grades:
    get:
      operationId: getGrades
      summary: Get grades for a specific course
      parameters:
        - name: course_id
          in: query
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Grades

  /api/assignments:
    get:
      operationId: getAssignments
      summary: Get assignments for a specific course
      parameters:
        - name: course_id
          in: query
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Assignments

  /api/assignments/upcoming:
    get:
      operationId: getUpcomingAssignments
      summary: Get upcoming assignments across all courses
      parameters:
        - name: days_ahead
          in: query
          schema:
            type: integer
            default: 30
      responses:
        '200':
          description: Upcoming assignments

  /api/assignments/submit:
    post:
      operationId: submitAssignment
      summary: Submit a text assignment
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [assignment_id, text]
              properties:
                assignment_id:
                  type: integer
                text:
                  type: string
      responses:
        '200':
          description: Submitted successfully

  /api/assignments/update:
    post:
      operationId: updateAssignment
      summary: Update (overwrite) an existing assignment submission
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [assignment_id, text]
              properties:
                assignment_id:
                  type: integer
                text:
                  type: string
      responses:
        '200':
          description: Updated successfully

  /api/journal/entry:
    get:
      operationId: getJournalEntry
      summary: Get current journal entry (e.g. Technical Article Review, Research Paper Review)
      parameters:
        - name: journal_id
          in: query
          required: true
          description: Journal module ID from course contents (modname=journal)
          schema:
            type: integer
      responses:
        '200':
          description: Journal entry text

  /api/journal/submit:
    post:
      operationId: submitJournal
      summary: Submit or update a journal entry (e.g. Technical Article Review, Research Paper Review)
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [journal_id, text]
              properties:
                journal_id:
                  type: integer
                  description: Journal module ID from course contents (modname=journal)
                text:
                  type: string
                  description: Journal entry text (HTML supported)
      responses:
        '200':
          description: Saved successfully

  /api/calendar/deadlines:
    get:
      operationId: getUpcomingDeadlines
      summary: Get upcoming deadlines sorted by urgency
      parameters:
        - name: days_ahead
          in: query
          schema:
            type: integer
            default: 14
      responses:
        '200':
          description: List of deadlines

  /api/notifications:
    get:
      operationId: getNotifications
      summary: Get messages and notifications
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
        - name: unread_only
          in: query
          schema:
            type: boolean
            default: true
      responses:
        '200':
          description: Notifications
```

## Step 5: Save and Test

Try asking your Custom GPT:
- "Log me into my Moodle at online.uom.lk"
- "Show my enrolled courses"
- "What are my grades?"
- "What assignments are due this week?"
- "Show my course contents for course 31778" *(to find journal IDs)*
- "Submit my Technical Article Review journal (ID 527772) with this text: ..."
- "What deadlines are coming up?"

---

## Troubleshooting

**"URL not reachable"** — Make sure your REST server is running and ngrok is active.

**"Action failed"** — Make sure you called `/api/login` first before any other action.

**ngrok URL expired** — Get a new one with `ngrok http 8080` and update your GPT action schema.

---

## Production: Keep Your API Running 24/7

See `DEPLOYMENT_GUIDE.md` to deploy to cloud platforms (Heroku, Google Cloud Run, AWS, Docker).
