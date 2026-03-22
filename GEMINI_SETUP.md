# Google Gemini Setup Guide

Use Moodle MCP Server with Google Gemini via the REST API and Google Apps Script.

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

## Step 2: Expose Your Server to the Internet

### Option A: ngrok (Easiest for Testing)

```bash
brew install ngrok  # macOS
ngrok http 8080
```

You'll get a URL like: `https://abc123.ngrok.io`

### Option B: Deploy to Cloud (Recommended)

See `DEPLOYMENT_GUIDE.md`. Google Cloud Run is the best fit for Gemini.

## Step 3: Set Up Google Apps Script

1. Go to [Google Apps Script](https://script.google.com)
2. Create a new project
3. Paste the full script below (replace `YOUR_URL` with your server URL)
4. **Save** and **Deploy as web app**

```javascript
const MOODLE_API = "https://YOUR_URL";

// ── Authentication ──────────────────────────────────────────────

function moodleLogin(moodleUrl, username, password) {
  return post_("/api/login", { moodle_url: moodleUrl, username, password });
}

function getSiteInfo() {
  return get_("/api/site-info");
}

function getUserProfile() {
  return get_("/api/user-profile");
}

// ── Courses ─────────────────────────────────────────────────────

function listCourses() {
  return get_("/api/courses");
}

function getCourseDetails(courseId) {
  return get_("/api/courses/details?course_id=" + courseId);
}

function getCourseContents(courseId) {
  // Use this to find journal IDs (look for modname=journal in the response)
  return get_("/api/courses/contents?course_id=" + courseId);
}

// ── Grades ──────────────────────────────────────────────────────

function getGrades(courseId) {
  return get_("/api/grades?course_id=" + courseId);
}

function getGradesOverview() {
  return get_("/api/grades/overview");
}

// ── Assignments ─────────────────────────────────────────────────

function getAssignments(courseId) {
  return get_("/api/assignments?course_id=" + courseId);
}

function getUpcomingAssignments(daysAhead = 30) {
  return get_("/api/assignments/upcoming?days_ahead=" + daysAhead);
}

function submitAssignment(assignmentId, text) {
  return post_("/api/assignments/submit", { assignment_id: assignmentId, text });
}

function updateAssignment(assignmentId, text) {
  return post_("/api/assignments/update", { assignment_id: assignmentId, text });
}

// ── Journal ─────────────────────────────────────────────────────
// Journal activities include: Technical Article Review, Research Paper Review, etc.
// To find journal_id: call getCourseContents(courseId) and look for modules where modname="journal"

function getJournalEntry(journalId) {
  return get_("/api/journal/entry?journal_id=" + journalId);
}

function submitJournal(journalId, text) {
  return post_("/api/journal/submit", { journal_id: journalId, text });
}

// ── Calendar ────────────────────────────────────────────────────

function getCalendarEvents(daysAhead = 30) {
  return get_("/api/calendar/events?days_ahead=" + daysAhead);
}

function getUpcomingDeadlines(daysAhead = 14) {
  return get_("/api/calendar/deadlines?days_ahead=" + daysAhead);
}

// ── Notifications ───────────────────────────────────────────────

function getNotifications(limit = 20, unreadOnly = true) {
  return get_("/api/notifications?limit=" + limit + "&unread_only=" + unreadOnly);
}

// ── Internal helpers ────────────────────────────────────────────

function get_(path) {
  const resp = UrlFetchApp.fetch(MOODLE_API + path, {
    method: "get", muteHttpExceptions: true,
    headers: { "Content-Type": "application/json" }
  });
  return JSON.parse(resp.getContentText());
}

function post_(path, body) {
  const resp = UrlFetchApp.fetch(MOODLE_API + path, {
    method: "post",
    payload: JSON.stringify(body),
    headers: { "Content-Type": "application/json" },
    muteHttpExceptions: true
  });
  return JSON.parse(resp.getContentText());
}
```

## Step 4: Test with Gemini

In [Google Gemini Advanced](https://gemini.google.com), try:

- "Use the Moodle API to log in with my credentials"
- "Call listCourses() and show me what courses I have"
- "Get my grades overview"
- "What assignments are due in the next 7 days?"
- "Get the course contents for course 31778 to find my journal IDs"
- "Submit my journal entry (ID 527772) with this text: ..."
- "Read my current Technical Article Review 2 entry (journal ID 527772)"

## Step 5: Set Up Automated Reminders (Optional)

1. In Google Apps Script, click **Triggers** (⏰)
2. Click **Create new trigger**
3. Set a daily trigger for:

```javascript
function dailyMoodleCheck() {
  const deadlines = getUpcomingDeadlines(7);
  const notifications = getNotifications(10, true);
  // Send yourself an email summary
  MailApp.sendEmail(
    Session.getActiveUser().getEmail(),
    "Daily Moodle Summary",
    "Deadlines: " + JSON.stringify(deadlines, null, 2) +
    "\n\nNotifications: " + JSON.stringify(notifications, null, 2)
  );
}
```

## Step 6: Use in Google Workspace

- **Google Sheets**: Pull grades into a spreadsheet with `getGrades()` or `getGradesOverview()`
- **Google Docs**: Draft and submit journal entries directly via `submitJournal()`
- **Gmail**: Auto-trigger reminders via `getUpcomingDeadlines()`

---

## Troubleshooting

**"URL not reachable"** — Make sure your REST server is running and ngrok is active (ngrok URLs expire).

**"Authorization failed" in Apps Script** — Test the endpoint in a browser first: `https://YOUR_URL/health`

**Can't find journal ID** — Call `getCourseContents(courseId)` and look for entries where `modname` is `"journal"`. The `id` field is your `journal_id`.

---

## Production: Deploy for 24/7 Access

See `DEPLOYMENT_GUIDE.md`. **Google Cloud Run** is recommended for Gemini as it integrates seamlessly with Google services.
