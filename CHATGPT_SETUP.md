# ChatGPT Setup Guide

Use Moodle MCP Server with ChatGPT via the REST API and Custom GPT Actions.

## Step 1: Start the REST API Server

First, you need to run the REST API server on your machine or a cloud server.

### On Your Mac/Linux

```bash
cd ~/fajr
go run ./cmd/moodle-mcp/ -mode rest -port 8080
```

### On Windows (PowerShell)

```powershell
cd C:\Users\YourName\moodle-mcp-server
go run .\cmd\moodle-mcp\ -mode rest -port 8080
```

You should see:
```
Starting REST API on port 8080
OpenAPI docs at http://localhost:8080/api/docs
```

## Step 2: Expose Your Server to the Internet (Optional)

For ChatGPT to access it, the server needs to be reachable from the internet. You have a few options:

### Option A: Use ngrok (Easiest for Testing)

```bash
brew install ngrok  # macOS
# or download from https://ngrok.com

ngrok http 8080
```

You'll get a URL like: `https://abc123.ngrok.io`

### Option B: Deploy to Cloud (Better for Production)

See `DEPLOYMENT_GUIDE.md` for Docker and cloud hosting options.

### Option C: Use Your Own Domain/Server

If you have a server with a domain, point it to your local machine using port forwarding.

## Step 3: Create a Custom GPT in ChatGPT

1. Go to [ChatGPT Custom GPT Builder](https://chat.openai.com/gpts/editor)
2. Click "Create New GPT"
3. Give it a name: "Moodle Assistant"
4. Click "Configure" tab
5. Scroll down to "Actions"
6. Click "Create new action"

## Step 4: Add the OpenAPI Schema

In the Action configuration, paste this schema (replace `YOUR_URL` with your actual server URL):

```yaml
openapi: 3.0.0
info:
  title: Moodle API
  version: 1.0.0
servers:
  - url: https://YOUR_URL  # Replace with your ngrok URL or domain

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
              properties:
                moodle_url:
                  type: string
                  example: "https://online.uom.lk"
                username:
                  type: string
                  example: "student@uom.lk"
                password:
                  type: string
                  example: "your-password"
              required: [moodle_url, username, password]
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

  /api/grades/overview:
    get:
      operationId: gradesOverview
      summary: Get grade summary across all courses
      responses:
        '200':
          description: Grades summary

  /api/assignments/upcoming:
    get:
      operationId: upcomingAssignments
      summary: Get upcoming assignments
      parameters:
        - name: days_ahead
          in: query
          schema:
            type: integer
            default: 30
      responses:
        '200':
          description: List of upcoming assignments

  /api/calendar/deadlines:
    get:
      operationId: upcomingDeadlines
      summary: Get upcoming deadlines
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
      responses:
        '200':
          description: List of notifications
```

## Step 5: Save and Test

1. Click "Save"
2. In the preview chat, try asking:
   - "Log me into my Moodle"
   - "Show my courses"
   - "What are my grades?"
   - "What assignments are due?"

3. You'll be prompted for your Moodle URL, username, and password

## Step 6: Share Your Custom GPT (Optional)

1. Click "Share" to make it public
2. Share the link with your friends
3. They can use it immediately!

---

## Troubleshooting

**"URL not reachable" error?**
- Make sure your REST server is running
- If using ngrok, make sure it's still active
- Check that the URL in the OpenAPI spec matches your server URL

**"Action failed" error?**
- Make sure you're logged in first with the `/api/login` action
- Check that your Moodle URL is correct
- Verify your username and password

**Need to update the OpenAPI spec?**
- Edit the Custom GPT
- Go to the Action
- Update the schema
- Click Save

---

## Production: Keep Your API Running 24/7

For production, see `DEPLOYMENT_GUIDE.md` to deploy to:
- AWS Lambda
- Google Cloud Run
- Heroku
- DigitalOcean
- Docker
