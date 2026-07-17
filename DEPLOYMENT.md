# Production Deployment Guide for Accountabel AI

This document outlines the steps to prepare and deploy the Accountabel AI platform to a production environment.

## 1. Security & Configuration
- **Remove Hardcoded Tokens:** Avoid placing your Telegram Bot Token or Gemini API Token in source-controlled scripts (like `run.bat`). Use environment variables or a secure configuration file on the production server.
- **Google OAuth Origins:** In the [Google Cloud Console](https://console.cloud.google.com/), navigate to your OAuth Client ID (ending in `.apps.googleusercontent.com`) and update the **Authorized JavaScript origins** and **Authorized redirect URIs** to match your production domain (e.g., `https://yourdomain.com`). Localhost is allowed for development, but production domains must be explicitly authorized.
- **HTTPS/TLS:** The web dashboard uses Web APIs (like `localStorage` and crypto) which require a secure context (HTTPS) in modern browsers. You can either use a reverse proxy like Nginx/Caddy to handle SSL, or configure the Go server directly using the `-crt_file` and `-key_file` flags.

## 2. Building for Production
Do not use `go run` in production as it is slower and requires the Go SDK. Compile a standalone executable instead.

**Windows:**
Run `build.bat` or use the command:
`go build -ldflags "-s -w" -o AccountabelAI.exe main.go`

**Linux/macOS:**
Run `build.sh` or use the command:
`go build -ldflags "-s -w" -o AccountabelAI main.go`

## 3. Running in Production
Use the generated executable and pass your tokens securely. A template for starting the server is provided in `run_prod.bat` and `run_prod.sh`.

**Example start command:**
`./AccountabelAI -telegram_bot_token="YOUR_TELEGRAM_TOKEN" -gemini_token="YOUR_GEMINI_TOKEN" -type="gemini" -http_host=":80"`

*(It is highly recommended to run the Go server behind a reverse proxy like Nginx or Caddy on port 80/443, while keeping the Go app running on an internal port like `:36060`)*

## 4. Data Persistence
The SQLite database is stored in the `data/` directory by default. Ensure that your production environment (if using Docker or ephemeral cloud instances) mounts this directory as a persistent volume so user profiles and check-in logs are not lost during deployments or restarts.
