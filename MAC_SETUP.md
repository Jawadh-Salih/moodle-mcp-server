# macOS Setup Guide (Easy Mode)

For macOS users - follow these **4 simple steps**.

## Step 1: Run the Installer

Open **Terminal** (search for "Terminal" in Spotlight or Applications > Utilities)

Copy and paste this command:

```bash
curl -fsSL https://raw.githubusercontent.com/jawadh/moodle-mcp-server/main/install.sh | bash
```

Press **Enter** and wait for it to finish.

You'll see a message like:
```
✓ Downloaded to: /Users/yourname/.moodle-mcp/moodle-mcp
```

**Note:** The path is `~/.moodle-mcp/moodle-mcp` - you'll need this in Step 3.

## Step 2: Open Claude Desktop Config

In Terminal, run:

```bash
open ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

This opens the file in your default editor.

## Step 3: Add the MCP Server

Find the `"mcpServers"` section. If it doesn't exist, copy this whole thing:

```json
{
  "mcpServers": {
    "moodle": {
      "command": "/Users/yourname/.moodle-mcp/moodle-mcp"
    }
  }
}
```

**Important:** Replace `yourname` with your actual macOS username!

If `"mcpServers"` already exists, just add the `"moodle"` part inside it.

**Save** the file (⌘+S)

## Step 4: Restart Claude Desktop

- Close Claude Desktop completely
- Open Claude Desktop again
- You're done! ✓

## Using It

In Claude, ask anything like:
- "Log in to my Moodle account"
- "Show my courses"
- "What assignments are due?"
- "Check my grades"

When you ask to log in, you'll be asked for:
1. Your Moodle URL (e.g., `https://online.uom.lk`)
2. Your username/email
3. Your password

That's it! Your friends can follow the same steps.

---

## Troubleshooting

**"curl: command not found"?**
- macOS should have curl built-in. Try closing Terminal and opening a new one.

**Can't open the config file?**
- The path might be different on some Mac versions. Try:
  ```bash
  open ~/.config/Claude/claude_desktop_config.json
  ```

**"Invalid login" error?**
- Double-check your Moodle username and password
- Some universities use email as username

**"Permission denied" error?**
- Run this to fix permissions:
  ```bash
  chmod +x ~/.moodle-mcp/moodle-mcp
  ```

Need help? Check the main README.md file or ask your friends!
