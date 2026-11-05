#!/usr/bin/env bash
# eжедневный бэкап БД бота: дамп в ~/backups, хранится KEEP_DAYS дней
# eсли задан BACKUP_REMOTE (user@host:path/) - копия туда же через rsync
set -euo pipefail

BACKUP_DIR="$HOME/backups"
KEEP_DAYS=14
CONTAINER=vpnbot-env-postgres

mkdir -p "$BACKUP_DIR"
file="$BACKUP_DIR/vpnbot-$(date -u +%Y-%m-%dT%H-%M).dump"

# пользователь и имя БД берётся из окружения самого контейнера postgres
# .env не читается из-за значений с пробелами (bash их не переварит)
docker exec "$CONTAINER" sh -c 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc' > "$file.tmp"
mv "$file.tmp" "$file"

find "$BACKUP_DIR" -name 'vpnbot-*.dump' -mtime +"$KEEP_DAYS" -delete

if [ -n "${BACKUP_REMOTE:-}" ]; then
    rsync -a --delete "$BACKUP_DIR/" "$BACKUP_REMOTE"
fi

echo "$(date -u +%FT%TZ) ok $file ($(du -h "$file" | cut -f1))"