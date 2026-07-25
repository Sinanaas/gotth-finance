#!/bin/sh
# Periodic Postgres backup: timestamped gzipped pg_dump with bounded retention.
# Runs in a loop inside the `backup` sidecar. Configure via env (see app.env):
#   BACKUP_INTERVAL_SECONDS (default 86400 = daily)
#   BACKUP_RETENTION        (default 14 = keep newest 14 dumps)
set -eu

INTERVAL="${BACKUP_INTERVAL_SECONDS:-86400}"
RETENTION="${BACKUP_RETENTION:-14}"
OUT_DIR=/backups

mkdir -p "$OUT_DIR"

while true; do
	TS=$(date +%Y%m%d-%H%M%S)
	FILE="$OUT_DIR/${POSTGRES_DB}-${TS}.sql.gz"
	echo "[backup] dumping ${POSTGRES_DB} -> ${FILE}"
	if PGPASSWORD="$POSTGRES_PASSWORD" pg_dump -h db -U "$POSTGRES_USER" "$POSTGRES_DB" | gzip > "$FILE"; then
		echo "[backup] ok: ${FILE}"
	else
		echo "[backup] FAILED at ${TS}" >&2
		rm -f "$FILE"
	fi

	# Prune: keep only the newest $RETENTION dumps.
	ls -1t "$OUT_DIR"/${POSTGRES_DB}-*.sql.gz 2>/dev/null | tail -n +$((RETENTION + 1)) | while read -r old; do
		echo "[backup] pruning ${old}"
		rm -f "$old"
	done

	sleep "$INTERVAL"
done
