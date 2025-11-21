#!/bin/sh

set -e

if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker first."
    exit 1
fi
echo "Docker is running"

cd "$(dirname "$0")"

docker compose -f docker-compose.dev.yml up --build


