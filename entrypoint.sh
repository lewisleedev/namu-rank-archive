#!/bin/ash

echo "[$(date "+%Y-%m-%d %H:%M:%S")] Container started"

echo "[$(date "+%Y-%m-%d %H:%M:%S")] Migration started"
/app/namu-rank-archive migrate
echo "[$(date "+%Y-%m-%d %H:%M:%S")] Migration finished"

echo "*/30 * * * * /app/namu-rank-archive archive
# " >scheduler.txt

crontab scheduler.txt

echo "[$(date "+%Y-%m-%d %H:%M:%S")] Crontab installed"

crond -f
