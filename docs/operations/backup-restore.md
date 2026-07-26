# Backup & Restore

The `backup` service in `docker-compose.yml` runs `scripts/backup.sh`, which takes a
gzipped `pg_dump` of the database on a schedule and keeps a bounded number of dumps in the
`db_backups` volume.

## Configuration (in `app.env`)

- `BACKUP_INTERVAL_SECONDS` — seconds between dumps (default `86400`, i.e. daily)
- `BACKUP_RETENTION` — how many dumps to keep (default `14`)

Uses the same `POSTGRES_*` credentials already in `app.env`.

## Where dumps live

Inside the named volume `db_backups`, as `${POSTGRES_DB}-YYYYMMDD-HHMMSS.sql.gz`.

List them:

```
docker compose exec backup ls -lh /backups
```

Copy one to the host:

```
docker compose cp backup:/backups/<file>.sql.gz ./<file>.sql.gz
```

## Restore

Restore a dump into a running database (this overwrites current data — be sure):

```
gunzip -c <file>.sql.gz | docker compose exec -T db psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"
```

To restore into a clean database instead, drop and recreate it first, then run the command above.

## Restore drill (do this periodically)

A backup you have never restored is not a backup. Every so often:

1. Spin up a scratch Postgres (or a throwaway compose project).
2. Restore the newest dump into it.
3. Confirm row counts for `users`, `accounts`, `transactions` look right.

If the restore fails, fix it before you need it for real.
