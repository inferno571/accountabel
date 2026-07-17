@echo off
echo Starting Accountabel AI Production Server...

:: WARNING: DO NOT COMMIT THIS FILE IF YOU HARDCODE TOKENS HERE!
:: It is recommended to use environment variables instead, e.g., %TELEGRAM_TOKEN%

set TELEGRAM_TOKEN="YOUR_TELEGRAM_BOT_TOKEN_HERE"
set GEMINI_TOKEN="YOUR_GEMINI_API_TOKEN_HERE"

:: Note: In production, you might want to run this behind Nginx or Caddy.
:: -http_host=":36060" binds the server to port 36060.
.\AccountabelAI.exe -telegram_bot_token=%TELEGRAM_TOKEN% -gemini_token=%GEMINI_TOKEN% -type="gemini" -http_host=":36060"
