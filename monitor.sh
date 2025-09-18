#!/bin/bash

# Discord webhook URL (replace with your webhook)
DISCORD_WEBHOOK="https://discord.com/api/webhooks/YOUR_WEBHOOK_URL_HERE"

# Check if InstaFix container is running
if ! docker compose ps instafix | grep -q "Up"; then
    # Send Discord notification
    curl -H "Content-Type: application/json" \
         -X POST \
         -d "{\"content\":\"🚨 **InstaFix Alert** 🚨\n\nInstaFix container has stopped running on $(hostname)\n\nTime: $(date)\n\nPlease check the logs: \`docker compose logs instafix\`\"}" \
         $DISCORD_WEBHOOK

    echo "$(date): InstaFix container is down, notification sent to Discord"
fi