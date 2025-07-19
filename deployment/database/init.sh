#!/bin/bash
set -e

echo "⏳ Waiting for Postgres to be ready..."

until pg_isready -U "$POSTGRES_USER"; do
  sleep 1
done

if [ "$PRELOAD_SPQR" = "true" ]; then
  echo "✅ Running spqr.sql..."
  # spqr.sql is renamed to spqr.sql.disabled during image build so that it is not run automatically
  psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/spqr.sql.disabled
else
  echo "⏭ Skipping spqr.sql.disabled preload."
fi
