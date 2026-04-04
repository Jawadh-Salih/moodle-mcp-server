#!/bin/bash
# Moodle MCP Server - macOS/Linux Installer
# Run this to automatically download and install

set -e

# Detect OS
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  Darwin)
    if [ "$ARCH" = "arm64" ]; then
      BINARY_NAME="moodle-mcp-macos-arm64"
    else
      BINARY_NAME="moodle-mcp-macos-amd64"
    fi
    ;;
  Linux)
    BINARY_NAME="moodle-mcp-linux-amd64"
    ;;
  *)
    echo "ERROR: Unsupported OS: $OS"
    exit 1
    ;;
esac

INSTALL_DIR="$HOME/.moodle-mcp"
BINARY_PATH="$INSTALL_DIR/moodle-mcp"
DOWNLOAD_URL="https://github.com/Jawadh-Salih/moodle-mcp-server/releases/download/v1.0.0/$BINARY_NAME"

echo "========================================"
echo "Moodle MCP Server - Installer"
echo "========================================"
echo ""

# Create installation directory
mkdir -p "$INSTALL_DIR"
echo "✓ Created directory: $INSTALL_DIR"

# Download the binary
echo "Downloading Moodle MCP binary..."
if command -v curl &> /dev/null; then
  curl -L -o "$BINARY_PATH" "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
  wget -O "$BINARY_PATH" "$DOWNLOAD_URL"
else
  echo "ERROR: Please install curl or wget"
  exit 1
fi

# Make binary executable
chmod +x "$BINARY_PATH"
echo "✓ Downloaded to: $BINARY_PATH"

echo ""
echo "========================================"
echo "Installation Complete!"
echo "========================================"
echo ""
echo "Next steps:"
echo "1. Open Claude Desktop configuration file:"
echo "   macOS:  ~/Library/Application Support/Claude/claude_desktop_config.json"
echo "   Linux:  ~/.config/Claude/claude_desktop_config.json"
echo ""
echo "2. Add this to your MCP servers (edit the JSON):"
echo ""
cat << 'EOF'
{
  "mcpServers": {
    "moodle": {
      "command": "$HOME/.moodle-mcp/moodle-mcp"
    }
  }
}
EOF
echo ""
echo "3. Restart Claude Desktop"
echo "4. In Claude, use the 'login' tool to authenticate with your Moodle account"
echo ""
echo "Questions? See README.md or visit: https://github.com/Jawadh-Salih/moodle-mcp-server"
echo ""
