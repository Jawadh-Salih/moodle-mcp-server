# Moodle MCP Server - Windows Installer
# Run this in PowerShell to automatically download and install

param(
    [string]$InstallDir = "$env:USERPROFILE\moodle-mcp"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Moodle MCP Server - Windows Installer" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if curl/wget is available
$hasWget = $null -ne (Get-Command wget -ErrorAction SilentlyContinue)
$hasCurl = $null -ne (Get-Command curl -ErrorAction SilentlyContinue)

if (-not $hasWget -and -not $hasCurl) {
    Write-Host "ERROR: Please install curl or wget to download the binary" -ForegroundColor Red
    exit 1
}

# Create installation directory
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
    Write-Host "✓ Created directory: $InstallDir" -ForegroundColor Green
}

# Download the Windows binary
$BinaryName = "moodle-mcp.exe"
$BinaryPath = Join-Path $InstallDir $BinaryName
$DownloadURL = "https://github.com/Jawadh-Salih/moodle-mcp-server/releases/download/v1.0.0/moodle-mcp-windows-amd64.exe"

Write-Host "Downloading Moodle MCP binary..." -ForegroundColor Yellow
try {
    if ($hasCurl) {
        curl -L -o $BinaryPath $DownloadURL
    } else {
        wget -O $BinaryPath $DownloadURL
    }
    Write-Host "✓ Downloaded to: $BinaryPath" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Failed to download binary" -ForegroundColor Red
    Write-Host "Manual download: $DownloadURL" -ForegroundColor Yellow
    exit 1
}

# Make binary executable (Windows doesn't require this, but good practice)
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Installation Complete!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Open Claude Desktop configuration file:" -ForegroundColor White
Write-Host "   Path: %APPDATA%\Claude\claude_desktop_config.json" -ForegroundColor Cyan
Write-Host ""
Write-Host "2. Add this to your MCP servers (edit the JSON):" -ForegroundColor White
Write-Host ""
Write-Host '{
  "mcpServers": {
    "moodle": {
      "command": "' + $BinaryPath + '"
    }
  }
}' -ForegroundColor Cyan
Write-Host ""
Write-Host "3. Restart Claude Desktop" -ForegroundColor White
Write-Host "4. In Claude, use the 'login' tool to authenticate with your Moodle account" -ForegroundColor White
Write-Host ""
Write-Host "Questions? See README.md or visit: https://github.com/Jawadh-Salih/moodle-mcp-server" -ForegroundColor Yellow
Write-Host ""
