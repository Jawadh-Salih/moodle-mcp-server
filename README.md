# Moodle MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that gives Claude direct, structured access to any Moodle LMS instance. Written in Go — single binary, no runtime dependencies.

Built for students. Ask Claude natural questions about your courses, grades, assignments, and deadlines instead of clicking through Moodle's UI.

---

## What You Can Do

| Tool | Description |
|------|-------------|
| `login` | Authenticate interactively at runtime |
| `get_site_info` | View Moodle site and current user info |
| `get_user_profile` | View your full profile |
| `list_courses` | List all enrolled courses |
| `get_course_contents` | Browse sections, resources, and activities |
| `get_course_details` | View course metadata (format, dates, category) |
| `get_grades` | View grades for a specific course |
| `get_grades_overview` | Grade summary across all courses |
| `get_assignments` | Assignments for a specific course |
| `get_upcoming_assignments` | All upcoming assignments across courses, sorted by due date |
| `submit_assignment` | Submit text content for an online text assignment |
| `get_calendar_events` | Upcoming calendar events from all courses |
| `get_upcoming_deadlines` | Consolidated deadlines sorted by urgency |
| `get_notifications` | Messages and notifications |

**Example questions you can ask Claude:**

- *"What are my grades across all courses?"*
- *"Do I have any assignments due this week?"*
- *"Show me the contents of my Digital Signal Processing course."*
- *"What's the most urgent deadline I have right now?"*
- *"Do I have any unread notifications?"*

---

## Requirements

- [Go 1.23+](https://go.dev/dl/)
- A Moodle account at any institution that has the Mobile Web Service enabled

---

## Installation

```bash
git clone https://github.com/Jawadh-Salih/moodle-mcp-server.git
cd moodle-mcp-server
make build
```

The binary is written to `./moodle-mcp`.

---

## Authentication

Two authentication modes are supported.

### Mode 1: Token (Recommended)

If you have a Moodle API token (available from your Moodle profile under *Preferences > Security keys*), use it directly. The token never passes through Claude's context window.

```bash
export MOODLE_URL=https://moodle.youruniversity.edu
export MOODLE_TOKEN=your-api-token
```

### Mode 2: Interactive Login

Start the server with no configuration and authenticate through Claude at runtime using the `login` tool:

> *"Log in to my Moodle at https://moodle.youruniversity.edu with username student123"*

> **Security note:** The `login` tool accepts a password as a tool argument. This means the password passes through Claude's context window and may be retained in conversation history. Use token-based authentication (Mode 1) in shared or persistent environments.

### Mode 3: Credentials via Environment

```bash
export MOODLE_URL=https://moodle.youruniversity.edu
export MOODLE_USERNAME=your-username
export MOODLE_PASSWORD=your-password
```

Credentials are exchanged for a token at startup; the password is not stored after that.

---

## Configuration

### Claude Desktop

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

**Token-based (recommended):**

```json
{
  "mcpServers": {
    "moodle": {
      "command": "/absolute/path/to/moodle-mcp",
      "env": {
        "MOODLE_URL": "https://moodle.youruniversity.edu",
        "MOODLE_TOKEN": "your-api-token"
      }
    }
  }
}
```

**Interactive login (no credentials in config):**

```json
{
  "mcpServers": {
    "moodle": {
      "command": "/absolute/path/to/moodle-mcp"
    }
  }
}
```

### Claude Code

Add to `.claude/settings.json`:

```json
{
  "mcpServers": {
    "moodle": {
      "command": "/absolute/path/to/moodle-mcp",
      "env": {
        "MOODLE_URL": "https://moodle.youruniversity.edu",
        "MOODLE_TOKEN": "your-api-token"
      }
    }
  }
}
```

---

## Development

```bash
# Run tests
make test

# Run tests with race detector and verbose output
make test-v

# Run all benchmarks
make bench

# Generate HTML coverage report
make coverage

# Print coverage summary
make coverage-text

# Run go vet
make lint
```

### Benchmark Results (Apple M1 Pro)

| Benchmark | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| `ClientCall` (mock server round-trip) | 55,939 | 8,715 | 104 |
| `HandleListCourses` (10 courses, mock) | 87,715 | 21,165 | 265 |
| `NormalizeURL` | 10 | 0 | 0 |
| `TruncateASCII` | 434 | 448 | 3 |
| `StripHTMLSimple` | 548 | 129 | 6 |
| `StripHTMLComplex` | 2,601 | 482 | 8 |

### Project Structure

```
.
├── cmd/moodle-mcp/       # Entry point: MCP server setup and tool registration
├── internal/
│   ├── api/              # HTTP client, authentication, error types
│   ├── config/           # Environment config, URL normalisation
│   └── tools/            # One file per domain (courses, grades, assignments…)
├── Makefile
└── README.md
```

**Layer responsibilities:**

- `internal/api` — all HTTP communication with Moodle. No business logic.
- `internal/config` — load and validate configuration. No HTTP.
- `internal/tools` — one handler per MCP tool. Calls `api.Client`, formats output. No HTTP directly.
- `cmd/moodle-mcp/main.go` — wires everything together, registers tools with the MCP server.

---

## How It Works

This server uses Moodle's [Web Services REST API](https://docs.moodle.org/dev/Web_service_API_functions) with the `moodle_mobile_app` service, which is enabled by default on virtually all Moodle installations. No admin configuration is required.

- Authentication credentials are sent in the **HTTP POST body**, not the URL, to avoid appearing in server access logs.
- API responses are capped at **10 MB** to prevent memory exhaustion.
- Course names are fetched once and reused across all tools — no N+1 HTTP calls.
- HTTP clients have explicit **30-second timeouts** on all connections.

---

## Troubleshooting

**"Invalid login" error:** Double-check your username and password. Some institutions use an email address as the username; others use a student ID.

**"Web service not available":** Your Moodle admin may have disabled the Mobile Web Service. Ask them to enable it under *Site administration → Plugins → Web services → Mobile*.

**"User ID not set":** Call `get_site_info` after logging in, or use token-based auth which sets the user ID automatically at startup.

**Grades not showing:** The grade report API requires the `gradereport/user:view` capability. This is standard for students but may be restricted on some sites.

---

## License

MIT
