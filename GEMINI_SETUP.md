# Google Gemini Setup Guide

Use Moodle MCP Server with Google Gemini via the REST API and Google Extensions.

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

## Step 2: Expose Your Server to the Internet

For Gemini to access it, the server needs to be reachable from the internet.

### Option A: Use ngrok (Easiest for Testing)

```bash
brew install ngrok  # macOS
# or download from https://ngrok.com

ngrok http 8080
```

You'll get a URL like: `https://abc123.ngrok.io`

### Option B: Deploy to Cloud (Better for Production)

See `DEPLOYMENT_GUIDE.md` for Docker and cloud hosting options.

## Step 3: Access Gemini with Extensions

Google Gemini uses a different extension system than ChatGPT.

### Option A: Use via Google Apps Script (Recommended)

1. Go to [Google Apps Script](https://script.google.com)
2. Create a new project
3. Paste this code:

```javascript
function getMoodleCourses() {
  const response = UrlFetchApp.fetch("https://YOUR_URL/api/courses", {
    method: "get",
    muteHttpExceptions: true,
    headers: {
      "Content-Type": "application/json"
    }
  });
  return JSON.parse(response.getContentText());
}

function moodleLogin(moodleUrl, username, password) {
  const response = UrlFetchApp.fetch("https://YOUR_URL/api/login", {
    method: "post",
    payload: JSON.stringify({
      moodle_url: moodleUrl,
      username: username,
      password: password
    }),
    headers: {
      "Content-Type": "application/json"
    },
    muteHttpExceptions: true
  });
  return JSON.parse(response.getContentText());
}

function getMoodleGrades() {
  const response = UrlFetchApp.fetch("https://YOUR_URL/api/grades/overview", {
    method: "get",
    muteHttpExceptions: true
  });
  return JSON.parse(response.getContentText());
}

function getMoodleAssignments(daysAhead = 30) {
  const response = UrlFetchApp.fetch(`https://YOUR_URL/api/assignments/upcoming?days_ahead=${daysAhead}`, {
    method: "get",
    muteHttpExceptions: true
  });
  return JSON.parse(response.getContentText());
}

function getMoodleDeadlines(daysAhead = 14) {
  const response = UrlFetchApp.fetch(`https://YOUR_URL/api/calendar/deadlines?days_ahead=${daysAhead}`, {
    method: "get",
    muteHttpExceptions: true
  });
  return JSON.parse(response.getContentText());
}

function getMoodleNotifications(limit = 20) {
  const response = UrlFetchApp.fetch(`https://YOUR_URL/api/notifications?limit=${limit}`, {
    method: "get",
    muteHttpExceptions: true
  });
  return JSON.parse(response.getContentText());
}
```

4. **Save** and **Deploy** as API
5. In Gemini, you can reference these functions:
   - "Call getMoodleCourses()"
   - "Call moodleLogin(...)"
   - etc.

### Option B: Direct REST API Usage

Gemini can also call REST APIs directly. In your Gemini chat:

```
Call this API: https://YOUR_URL/api/courses
Method: GET
```

Gemini will make the request and show you the results.

## Step 4: Test with Gemini

1. Go to [Google Gemini Advanced](https://gemini.google.com)
2. Try these prompts:
   - "Use the Moodle API to show my courses"
   - "Call the login function with my Moodle credentials"
   - "Get my grade overview"
   - "Show me upcoming assignments"

## Step 5: Use in Google Workspace

You can integrate Moodle data into:
- **Google Sheets**: Fetch data with `=ImportData()` or Apps Script
- **Google Docs**: Embed Moodle data via Apps Script
- **Google Classroom**: Sync assignments and deadlines
- **Gmail**: Create filters/rules based on Moodle data

---

## Option C: Use Gemini's Scheduled Functions

For automated daily updates, you can schedule Apps Script functions:

1. In Google Apps Script, click **Triggers** (⏰)
2. Click **Create new trigger**
3. Set up a daily trigger for:
   - `getMoodleDeadlines()` to send you daily reminders
   - `getMoodleAssignments()` to check for new work

---

## Troubleshooting

**"URL not reachable" error?**
- Make sure your REST server is running
- If using ngrok, make sure it's still active
- Check that the URL matches your server URL
- ngrok URLs expire - get a new one if needed

**Apps Script says "Authorization failed"?**
- Make sure your Moodle server is accessible
- Check that `MOODLE_URL` in your REST API is correct
- Try testing the endpoint in a browser first

**Want to use with Gemini in Google Search?**
- Gemini's search extensions are limited
- For now, use the Google Apps Script approach above
- Or deploy to cloud and access via browser

---

## Production: Deploy for 24/7 Access

For production, see `DEPLOYMENT_GUIDE.md` to deploy to:
- Google Cloud Run (easiest for Gemini)
- AWS Lambda
- Heroku
- DigitalOcean
- Docker

**Recommended for Gemini:** Google Cloud Run - it integrates seamlessly with Google services!
