#!/bin/bash
# =====================================================
# Accountabel AI – DigitalOcean Droplet Setup Script
# Run this once on a fresh Ubuntu 22.04 Droplet
# =====================================================
set -e

echo "=== Installing Docker ==="
apt-get update -y
apt-get install -y docker.io docker-compose curl git
systemctl enable docker
systemctl start docker

echo "=== Cloning repository ==="
cd /opt
git clone https://github.com/inferno571/accountabel.git accountabel
cd accountabel

echo "=== Creating data directory ==="
mkdir -p /opt/accountabel/data

echo "=== Building Docker image ==="
docker build -t accountabel-ai .

echo "=== Starting the application ==="
# Replace the values below with your real tokens before running!
docker run -d \
  --name accountabel \
  --restart always \
  -p 80:36060 \
  -v /opt/accountabel/data:/app/data \
  -e TELEGRAM_BOT_TOKEN="${TELEGRAM_BOT_TOKEN}" \
  -e GEMINI_TOKEN="${GEMINI_TOKEN}" \
  -e DISCORD_BOT_TOKEN="${DISCORD_BOT_TOKEN}" \
  -e TYPE="gemini" \
  -e HTTP_HOST=":36060" \
  accountabel-ai

echo ""
echo "=== Accountabel AI is running! ==="
echo "Access it at: http://$(curl -s ifconfig.me)"
echo "SQLite data stored persistently at: /opt/accountabel/data/"
