# Google Gemini Setup Guide (Windows)

Use Moodle MCP Server with Google Gemini on Windows. No coding experience required.

---

## Prerequisites

Before starting, make sure you have the following ready:

### Required Software

| Software | What It's For | How to Get It |
|----------|--------------|---------------|
| **Windows 10/11** | Operating system | You probably already have this |
| **PowerShell** | Running commands | Built into Windows (search "PowerShell" in Start menu) |
| **A web browser** | Google Apps Script, Gemini | Chrome, Edge, Firefox, etc. |
| **ngrok** | Exposing your local server to the internet | https://ngrok.com/download (free account required) |

### Required Accounts

| Account | What It's For | How to Get It |
|---------|--------------|---------------|
| **Google account** | Apps Script + Gemini access | https://accounts.google.com/signup (free) |
| **ngrok account** | Auth token for tunneling | https://dashboard.ngrok.com/signup (free) |
| **Moodle account** | Your university/institution LMS | Provided by your institution |

### Required Information

You will need these details during setup:

- **Moodle URL** -- your institution's Moodle address (e.g. `https://online.uom.lk`)
- **Moodle username** -- your login username or email
- **Moodle password** -- your login password
- **ngrok auth token** -- from https://dashboard.ngrok.com/get-started/your-authtoken (after signing up)

### Optional (For 24/7 Access)

If you want the server running permanently without keeping your PC on, you'll need a cloud account. See [Deployment Guide](DEPLOYMENT_GUIDE.md) after completing this setup.

| Platform | Cost |
|----------|------|
| **Google Cloud Run** | Free tier available |
| **Heroku** | $5-7/month |

---

## Step 1: Download the Server

Download **moodle-mcp-windows-amd64.exe** from:
https://github.com/Jawadh-Salih/moodle-mcp-server/releases/latest

Save it to a folder, e.g. `C:\moodle-mcp\`

---

## Step 2: Start the REST API Server

Open **PowerShell** (search for "PowerShell" in the Start menu) and run:

```powershell
cd C:\moodle-mcp
.\moodle-mcp-windows-amd64.exe -mode rest -port 8080
```

You should see:
```
REST API listening on http://localhost:8080
```

**Leave this window open** -- closing it stops the server.

---

## Step 3: Expose Your Server to the Internet

Gemini needs a public URL to reach your local server.

### Install ngrok

1. Go to https://ngrok.com/download
2. Download the **Windows** version
3. Extract the zip file to `C:\ngrok\`
4. Sign up for a free ngrok account at https://dashboard.ngrok.com/signup
5. Copy your auth token from the dashboard

### Run ngrok

Open a **new PowerShell window** (keep the server running in the first one):

```powershell
cd C:\ngrok
.\ngrok.exe authtoken YOUR_AUTH_TOKEN
.\ngrok.exe http 8080
```

You'll see something like:
```
Forwarding  https://abc123.ngrok-free.app -> http://localhost:8080
```

**Copy the `https://` URL** -- you'll need it in the next step.

### Test It

Open your browser and go to:
```
https://abc123.ngrok-free.app/health
```
(Replace with your actual ngrok URL)

You should see: `{"status":"ok","authenticated":false}`

---

## Step 4: Set Up Google Apps Script

This creates the bridge between Gemini and your Moodle server.

1. Go to [Google Apps Script](https://script.google.com) in your browser
2. Click **New project**
3. Delete everything in the editor
4. Paste the full script below
5. **Replace `YOUR_URL`** on line 1 with your ngrok URL (e.g. `https://abc123.ngrok-free.app`)

```javascript
const MOODLE_API = "https://YOUR_URL";

// -- Authentication -------------------------------------------------------

function moodleLogin(moodleUrl, username, password) {
  return post_("/api/login", { moodle_url: moodleUrl, username, password });
}

function getSiteInfo() {
  return get_("/api/site-info");
}

function getUserProfile() {
  return get_("/api/user-profile");
}

// -- Courses --------------------------------------------------------------

function listCourses() {
  return get_("/api/courses");
}

function getCourseDetails(courseId) {
  return get_("/api/courses/details?course_id=" + courseId);
}

function getCourseContents(courseId) {
  return get_("/api/courses/contents?course_id=" + courseId);
}

// -- Grades ---------------------------------------------------------------

function getGrades(courseId) {
  return get_("/api/grades?course_id=" + courseId);
}

function getGradesOverview() {
  return get_("/api/grades/overview");
}

// -- Assignments ----------------------------------------------------------

function getAssignments(courseId) {
  return get_("/api/assignments?course_id=" + courseId);
}

function getUpcomingAssignments(daysAhead) {
  daysAhead = daysAhead || 30;
  return get_("/api/assignments/upcoming?days_ahead=" + daysAhead);
}

function submitAssignment(assignmentId, text) {
  return post_("/api/assignments/submit", { assignment_id: assignmentId, text: text });
}

function updateAssignment(assignmentId, text) {
  return post_("/api/assignments/update", { assignment_id: assignmentId, text: text });
}

// -- Journal --------------------------------------------------------------

function getJournalEntry(journalId) {
  return get_("/api/journal/entry?journal_id=" + journalId);
}

function submitJournal(journalId, text) {
  return post_("/api/journal/submit", { journal_id: journalId, text: text });
}

// -- Calendar -------------------------------------------------------------

function getCalendarEvents(daysAhead) {
  daysAhead = daysAhead || 30;
  return get_("/api/calendar/events?days_ahead=" + daysAhead);
}

function getUpcomingDeadlines(daysAhead) {
  daysAhead = daysAhead || 14;
  return get_("/api/calendar/deadlines?days_ahead=" + daysAhead);
}

// -- Notifications --------------------------------------------------------

function getNotifications(limit, unreadOnly) {
  limit = limit || 20;
  unreadOnly = (unreadOnly !== undefined) ? unreadOnly : true;
  return get_("/api/notifications?limit=" + limit + "&unread_only=" + unreadOnly);
}

// -- Resources ------------------------------------------------------------

function listResources(courseId) {
  return get_("/api/resources?course_id=" + courseId);
}

// -- Internal helpers -----------------------------------------------------

function get_(path) {
  var resp = UrlFetchApp.fetch(MOODLE_API + path, {
    method: "get",
    muteHttpExceptions: true,
    headers: { "Content-Type": "application/json" }
  });
  return JSON.parse(resp.getContentText());
}

function post_(path, body) {
  var resp = UrlFetchApp.fetch(MOODLE_API + path, {
    method: "post",
    payload: JSON.stringify(body),
    headers: { "Content-Type": "application/json" },
    muteHttpExceptions: true
  });
  return JSON.parse(resp.getContentText());
}
```

6. Click **Save** (Ctrl+S)
7. Name the project something like "Moodle API"

---

## Step 5: Deploy as a Gemini Extension

1. In Google Apps Script, click **Deploy** (top right) > **New deployment**
2. Click the gear icon next to "Select type" and choose **Web app**
3. Set:
   - **Description**: Moodle API
   - **Execute as**: Me
   - **Who has access**: Only myself
4. Click **Deploy**
5. Click **Authorize access** and sign in with your Google account
6. Copy the **Web app URL**

---

## Step 6: Test with Gemini

Open [Google Gemini](https://gemini.google.com) and try:

- "Use the Moodle API to log in with my credentials"
- "Call listCourses() and show me what courses I have"
- "Get my grades overview"
- "What assignments are due in the next 7 days?"
- "Get the course contents for course 31778 to find my journal IDs"
- "Submit my journal entry (ID 527772) with this text: ..."

---

## Step 7: Set Up Automated Reminders (Optional)

Get daily email summaries of your Moodle deadlines:

1. In Google Apps Script, click **Triggers** (clock icon in the left sidebar)
2. Click **+ Add Trigger**
3. Configure:
   - **Function**: `dailyMoodleCheck`
   - **Event source**: Time-driven
   - **Type**: Day timer
   - **Time of day**: 7am - 8am (or whenever you want)
4. Click **Save**

Add this function to your script:

```javascript
function dailyMoodleCheck() {
  var deadlines = getUpcomingDeadlines(7);
  var notifications = getNotifications(10, true);
  MailApp.sendEmail(
    Session.getActiveUser().getEmail(),
    "Daily Moodle Summary",
    "Deadlines:\n" + JSON.stringify(deadlines, null, 2) +
    "\n\nNotifications:\n" + JSON.stringify(notifications, null, 2)
  );
}
```

---

## Step 8: Use in Google Workspace (Optional)

- **Google Sheets**: Pull grades into a spreadsheet with `getGrades()` or `getGradesOverview()`
- **Google Docs**: Draft and submit journal entries directly via `submitJournal()`
- **Gmail**: Auto-trigger reminders via `getUpcomingDeadlines()`

---

## Production: 24/7 Access Without ngrok

For always-on access (no need to keep PowerShell open), deploy to Google Cloud Run:

See [Deployment Guide](DEPLOYMENT_GUIDE.md) for step-by-step cloud deployment instructions. Google Cloud Run is recommended for Gemini as it integrates seamlessly with Google services.

---

## Troubleshooting

**"URL not reachable"**
- Make sure both PowerShell windows are open (server + ngrok)
- Check that your ngrok URL hasn't expired (free plan URLs change on restart)

**"Authorization failed" in Apps Script**
- Test the endpoint in your browser first: `https://YOUR_NGROK_URL/health`
- Make sure you authorized the Apps Script when deploying

**Server closes when I close PowerShell**
- The server runs inside PowerShell -- it must stay open
- For background running, use: `Start-Process -WindowStyle Hidden .\moodle-mcp-windows-amd64.exe -ArgumentList "-mode","rest","-port","8080"`

**ngrok URL keeps changing**
- Free ngrok accounts get a new URL each time you restart
- Update the `MOODLE_API` variable in your Apps Script each time
- For a permanent URL, upgrade to ngrok paid plan or deploy to cloud (see Deployment Guide)

**Can't find journal ID**
- Call `getCourseContents(courseId)` and look for entries where `modname` is `"journal"`. The `id` field is your `journal_id`.

---

## Need Help?

- Open an issue: https://github.com/Jawadh-Salih/moodle-mcp-server/issues
- See the main [README](README.md) for general troubleshooting
