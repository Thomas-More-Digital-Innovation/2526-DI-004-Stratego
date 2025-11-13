#!/bin/sh
# Lightweight entrypoint used in dev container to wait for db then run reflex
set -e

# wait for DB
if [ -n "$DB_HOST" ]; then
  echo "Waiting for DB at $DB_HOST:5432..."
  until nc -z $(echo ${DB_HOST} | cut -d: -f1) 5432; do
    sleep 0.5
  done
fi

exec /go/bin/reflex -r '\\.go$' -s -- sh -c 'exec go run . --server'
