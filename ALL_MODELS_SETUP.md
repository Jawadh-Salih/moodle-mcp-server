# Complete Multi-AI Setup Guide

Use your Moodle MCP Server with **Claude, ChatGPT, Gemini**, and any other AI model.

---

## 🚀 Quick Overview

| AI Model | Method | Difficulty | Free |
|----------|--------|-----------|------|
| **Claude** | MCP (Built-in) | Easy | ✅ Yes |
| **ChatGPT** | REST API + Actions | Medium | ⏳ Requires ChatGPT+ |
| **Google Gemini** | REST API + Apps Script | Medium | ✅ Free |
| **Any Other AI** | REST API (Direct HTTP) | Easy | ✅ Yes |

---

## Setup for Each AI Model

### 1️⃣ Claude (MCP - Recommended)

**Best for:** Students who use Claude

**Setup time:** 2 minutes

**Steps:**
1. Download and run the installer (Windows/Mac)
2. Add to Claude config
3. Start using!

👉 **[Follow Windows Setup](WINDOWS_SETUP.md)** | **[Follow macOS Setup](MAC_SETUP.md)**

---

### 2️⃣ ChatGPT (REST API + Custom Actions)

**Best for:** ChatGPT Plus subscribers who want seamless integration

**Setup time:** 15 minutes

**Requirements:**
- ChatGPT Plus subscription
- Deploy REST API to cloud (see below)

**Steps:**
1. Start REST API server
2. Deploy to cloud (Google Cloud Run, Heroku, etc.)
3. Create Custom GPT
4. Add OpenAPI spec
5. Start using!

👉 **[Follow ChatGPT Setup](CHATGPT_SETUP.md)**

---

### 3️⃣ Google Gemini (REST API + Apps Script)

**Best for:** Google Workspace users, integrating with Google Sheets/Docs

**Setup time:** 20 minutes

**Requirements:**
- Google account
- Deploy REST API to cloud (see below)

**Steps:**
1. Start REST API server
2. Deploy to cloud
3. Set up Google Apps Script
4. Use in Gemini chat or Google Workspace
5. (Optional) Schedule daily updates

👉 **[Follow Gemini Setup](GEMINI_SETUP.md)**

---

### 4️⃣ Any Other AI (Direct REST API)

**Best for:** Custom integrations, scripts, other AI models

**Setup time:** 10 minutes

**Requirements:**
- REST API endpoint (cloud or local)

**Usage:**
```bash
# Login
curl -X POST https://your-server/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "moodle_url": "https://your-moodle.edu",
    "username": "your-username",
    "password": "your-password"
  }'

# Get courses
curl https://your-server/api/courses

# Get grades
curl https://your-server/api/grades?course_id=123

# View all endpoints
curl https://your-server/api/docs
```

---

## Choosing Your Deployment Strategy

### Local Testing (Start Here)

Run REST API locally with ngrok for temporary internet access:

```bash
# Terminal 1: Start server
cd ~/fajr
go run ./cmd/moodle-mcp/ -mode rest -port 8080

# Terminal 2: Expose to internet (temporary)
ngrok http 8080
# Get URL: https://abc123.ngrok.io
```

**Pros:** Free, instant, no setup
**Cons:** Temporary, need ngrok running always

---

### Production (Long-term)

Deploy to cloud for permanent, always-on access:

| Platform | Setup | Cost | Best For |
|----------|-------|------|----------|
| **Google Cloud Run** | 5 min | Free tier | Gemini, quick test |
| **Heroku** | 5 min | $5-7/mo | Simple production |
| **DigitalOcean** | 10 min | $5-12/mo | Full control |
| **AWS Lambda** | 20 min | Pay-as-you-go | High traffic |
| **Docker VPS** | 15 min | Varies | Self-hosted |

👉 **[Follow Deployment Guide](DEPLOYMENT_GUIDE.md)**

---

## Recommended Setup Paths

### Path 1: Claude Only (Simplest)
1. Run the MCP installer
2. Done!
- Time: 2 minutes
- Cost: $0
- Complexity: ⭐

### Path 2: Claude + ChatGPT (Popular)
1. Run MCP installer for Claude
2. Deploy REST API to Heroku (free tier)
3. Create Custom GPT in ChatGPT
4. Done!
- Time: 30 minutes
- Cost: $5-7/month (Heroku, plus ChatGPT+ subscription)
- Complexity: ⭐⭐

### Path 3: All Platforms (Full Power)
1. Run MCP installer for Claude
2. Deploy REST API to Google Cloud Run (free tier)
3. Set up ChatGPT Custom GPT
4. Set up Gemini Apps Script
5. Use with any other AI via REST API
6. Done!
- Time: 1 hour
- Cost: $0 (free tiers) or $5-7/month
- Complexity: ⭐⭐⭐

---

## Server Modes Explained

Your Moodle server can run in three modes:

### Mode 1: MCP (Claude Only)
```bash
go run ./cmd/moodle-mcp/ -mode mcp
```
- ✅ Works with Claude Desktop/Code
- ❌ Doesn't work with ChatGPT/Gemini/others
- Uses: stdio protocol
- Best for: Claude-only users

### Mode 2: REST (ChatGPT/Gemini/Others)
```bash
go run ./cmd/moodle-mcp/ -mode rest -port 8080
```
- ❌ Doesn't work with Claude MCP
- ✅ Works with ChatGPT, Gemini, any HTTP client
- Uses: HTTP REST API
- Best for: ChatGPT/Gemini users

### Mode 3: BOTH (Experimental)
```bash
go run ./cmd/moodle-mcp/ -mode both -port 8080
```
- ✅ Works with Claude (stdio)
- ✅ Works with ChatGPT/Gemini (HTTP on port 8080)
- Uses: Both MCP and REST simultaneously
- Best for: Power users supporting all platforms
- Note: Requires special Claude Desktop config

---

## Common Questions

**Q: Can I use Claude and ChatGPT at the same time?**
A: Yes! Set Claude to use MCP (mode: mcp), and ChatGPT to use the REST API (mode: rest).

**Q: Do I need to pay for all services?**
A: No! Claude and Gemini are free. ChatGPT+ requires subscription ($20/mo). Hosting varies.

**Q: Can I keep it running 24/7?**
A: Yes, deploy to a cloud service (see Deployment Guide). MCP servers run when Claude is open.

**Q: Which platform should I choose?**
A:
- **Start small:** Use local ngrok for testing
- **One platform:** Deploy to Heroku or Google Cloud Run
- **Multiple platforms:** Deploy to DigitalOcean or AWS for more control

**Q: Can I self-host?**
A: Yes! Use Docker on your own server, or use a VPS like DigitalOcean.

**Q: Is my password secure?**
A: MCP: Stored in Claude config (local)
REST: Stored as environment variables (cloud), or passed at login (stateless)

**Q: Can I use without coding?**
A: Yes! Use the installers for Claude. For ChatGPT/Gemini, just follow the setup guides.

---

## Architecture Diagram

```
Your Moodle Account (https://online.uom.lk)
         ↓
    API Client (shared)
    ├─ Handles auth
    ├─ Makes API calls
    └─ Manages session
         ↓
   ┌─────┴─────┐
   ↓           ↓
MCP Server   REST API Server
   ↓           ↓
Claude    ChatGPT/Gemini/Others
   ↓           ↓
Your Computer  Cloud (Heroku, GCP, etc.)
```

---

## Next Steps

1. **Choose your AI:** Claude, ChatGPT, Gemini, or multiple
2. **Follow the relevant setup guide** above
3. **Start using it!** Ask your AI to:
   - "Log in to my Moodle"
   - "Show my courses"
   - "What assignments are due?"
   - "Check my grades"
4. **Share with friends** - they can follow the same steps!

---

## Need Help?

- **Claude issues:** Check [Claude Setup Guides](WINDOWS_SETUP.md)
- **ChatGPT issues:** Check [ChatGPT Setup](CHATGPT_SETUP.md)
- **Gemini issues:** Check [Gemini Setup](GEMINI_SETUP.md)
- **Deployment issues:** Check [Deployment Guide](DEPLOYMENT_GUIDE.md)
- **Technical questions:** Check main [README](README.md)

Enjoy! 🚀
