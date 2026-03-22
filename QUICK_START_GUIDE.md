# Moodle API — Quick Start for Windows (ChatGPT & Gemini)

No coding required. Just download, run, and connect.

---

## Step 1: Download the Binary

Download **moodle-mcp-windows-amd64.exe** from:
https://github.com/Jawadh-Salih/moodle-mcp-server/releases/latest

Save it somewhere easy, e.g. `C:\moodle-mcp\`

---

## Step 2: Run the Server

Open **PowerShell** and run:

```powershell
cd C:\moodle-mcp
.\moodle-mcp-windows-amd64.exe -mode rest -port 8080
```

You should see:
```
REST API listening on http://localhost:8080
```

Leave this window open — the server must stay running.

---

## Step 3: Expose It to the Internet

ChatGPT and Gemini need a public URL to reach your server.

### Download ngrok
1. Go to https://ngrok.com/download and download for Windows
2. Extract and open `ngrok.exe`
3. Run:
```powershell
.\ngrok.exe http 8080
```
4. Copy the `https://` URL it gives you (e.g. `https://abc123.ngrok-free.app`)

---

## Step 4A: Use with ChatGPT

1. Go to https://chat.openai.com/gpts/editor
2. Create a new GPT → **Configure** → **Actions** → **Create new action**
3. Paste the OpenAPI schema from `CHATGPT_SETUP.md` (replace `YOUR_URL` with your ngrok URL)
4. Save and test

**Try asking:**
- "Log me into Moodle at online.uom.lk with username X and password Y"
- "Show my courses"
- "What are my upcoming deadlines?"
- "Submit my Technical Article Review (journal ID 527772) with this text: ..."

---

## Step 4B: Use with Gemini (Google Apps Script)

1. Go to https://script.google.com → New project
2. Paste the full script from `GEMINI_SETUP.md` (replace `YOUR_URL` with your ngrok URL)
3. Save → Deploy as web app
4. In Gemini, call the functions:
   - `moodleLogin("https://online.uom.lk", "username", "password")`
   - `listCourses()`
   - `getGradesOverview()`
   - `getUpcomingDeadlines(7)`
   - `submitJournal(527772, "My review text...")`

---

## Quick Reference: All Available Endpoints

| What you want | How to ask |
|---|---|
| Login | POST `/api/login` |
| My courses | GET `/api/courses` |
| Course contents (find journal IDs) | GET `/api/courses/contents?course_id=X` |
| My grades | GET `/api/grades/overview` |
| Upcoming assignments | GET `/api/assignments/upcoming` |
| Submit assignment | POST `/api/assignments/submit` |
| Read journal entry | GET `/api/journal/entry?journal_id=X` |
| Submit/update journal | POST `/api/journal/submit` |
| Upcoming deadlines | GET `/api/calendar/deadlines` |
| Notifications | GET `/api/notifications` |

Full API docs available at `http://localhost:8080/api/docs` when the server is running.

---

## Tips

- **Keep PowerShell open** — closing it stops the server
- **ngrok URL changes** every time you restart ngrok (free plan) — update your ChatGPT/Gemini config when it does
- **Find your journal ID** — ask ChatGPT/Gemini to call `getCourseContents` for your course and look for entries where `modname` is `journal`
- **Need help?** Contact Jawadh or open an issue at https://github.com/Jawadh-Salih/moodle-mcp-server/issues
