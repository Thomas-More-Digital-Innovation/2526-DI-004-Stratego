#!/bin/sh

set -e

cd "$(dirname "$0")"

echo "Stopping Stratego development environment..."

docker compose -f docker-compose.dev.yml down

