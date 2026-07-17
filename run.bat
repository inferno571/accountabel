@echo off
echo Starting Accountabel AI Bot and Web Server...
go run main.go -telegram_bot_token="YOUR_TELEGRAM_TOKEN_HERE" -gemini_token="YOUR_GEMINI_TOKEN_HERE" -type="gemini"
pause
