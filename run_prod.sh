#!/bin/bash

echo "Starting Accountabel AI Production Server..."

# WARNING: DO NOT COMMIT THIS FILE IF YOU HARDCODE TOKENS HERE!
# It is recommended to use environment variables exported by your server (e.g. systemd).

TELEGRAM_TOKEN="YOUR_TELEGRAM_BOT_TOKEN_HERE"
GEMINI_TOKEN="YOUR_GEMINI_API_TOKEN_HERE"

# Ensure the executable exists
if [ ! -f "./AccountabelAI" ]; then
    echo "Executable not found. Please run ./build.sh first."
    exit 1
fi

# Run the server
./AccountabelAI -telegram_bot_token="$TELEGRAM_TOKEN" -gemini_token="$GEMINI_TOKEN" -type="gemini" -http_host=":36060"
