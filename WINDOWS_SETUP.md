# Windows Setup Guide (Easy Mode)

For Windows users - follow these **4 simple steps**.

## Step 1: Run the Installer

Open **PowerShell** (search for "PowerShell" in Windows):

Right-click PowerShell → **Run as Administrator**

Copy and paste this command:

```powershell
irm https://raw.githubusercontent.com/jawadh/moodle-mcp-server/main/install.ps1 | iex
```

Press **Enter** and wait for it to finish.

You'll see a message like:
```
✓ Downloaded to: C:\Users\YourName\moodle-mcp\moodle-mcp.exe
```

**Copy this path** - you'll need it in Step 3.

## Step 2: Open Claude Desktop Config

Open **File Explorer** and paste this in the address bar:

```
%APPDATA%\Claude
```

Press Enter. You should see a file called `claude_desktop_config.json`

Right-click it → **Open with** → **Notepad**

## Step 3: Add the MCP Server

Find the `"mcpServers"` section in the file. If it doesn't exist, copy this whole thing:

```json
{
  "mcpServers": {
    "moodle": {
      "command": "C:\\Users\\YourName\\moodle-mcp\\moodle-mcp.exe"
    }
  }
}
```

**Important:** Replace `YourName` with your actual Windows username!

If `"mcpServers"` already exists, just add the `"moodle"` part inside it.

**Save** the file (Ctrl+S)

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

**"Command not found" error?**
- Make sure PowerShell is running as Administrator
- Try opening a new PowerShell window

**Can't find `claude_desktop_config.json`?**
- Make sure you pasted `%APPDATA%\Claude` exactly in File Explorer
- Check that Claude Desktop is installed

**"Invalid login" error?**
- Double-check your Moodle username and password
- Some universities use email as username

Need help? Check the main README.md file or ask your friends!
