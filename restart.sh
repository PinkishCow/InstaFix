#!/bin/bash
# Monthly container restart script.
# Pulls latest code, rebuilds, and restarts containers.
# Gives up after 5 health check attempts.
#
# Install as a cron job:
#   crontab -e
#   0 4 1 * * cd /path/to/InstaFix && ./restart.sh >> /var/log/instafix-restart.log 2>&1

set -euo pipefail

MAX_RETRIES=5
RETRY_INTERVAL=30
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

log "Starting monthly restart"

# Pull latest changes
if git pull --ff-only; then
    log "Pulled latest changes"
else
    log "WARNING: git pull failed, rebuilding with current code"
fi

# Rebuild and restart
log "Rebuilding containers"
docker compose up --build -d

# Wait for healthy status
for i in $(seq 1 $MAX_RETRIES); do
    sleep $RETRY_INTERVAL
    if docker compose ps instafix | grep -q "healthy"; then
        log "InstaFix is healthy after attempt $i"
        exit 0
    fi
    log "Health check attempt $i/$MAX_RETRIES failed"
done

log "ERROR: InstaFix failed to become healthy after $MAX_RETRIES attempts"
log "Container logs:"
docker compose logs --tail=30 instafix
exit 1
