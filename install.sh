#!/bin/bash

owner="song940"
repo="chatgpt-cli"
target="/usr/local/bin/chatgpt"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# fix for macOS
[ "$ARCH" = "x86_64" ] && ARCH="amd64"

# Download the latest release from GitHub
DOWNLOAD_URL=$(curl -s "https://api.github.com/repos/${owner}/${repo}/releases/latest" | grep "browser_download_url.*$OS-$ARCH" | cut -d '"' -f 4)
if [ -z "$DOWNLOAD_URL" ]; then
    echo "Failed to get download URL for $OS/$ARCH architecture"
    exit 1
fi

curl -#L $DOWNLOAD_URL -o "$target"
chmod +x "$target"

echo "Successfully installed"