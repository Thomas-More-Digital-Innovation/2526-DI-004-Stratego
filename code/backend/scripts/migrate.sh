#!/usr/bin/env bash
# scripts/migrate.sh - Helper to run database migrations manually
set -euo pipefail

# Find project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/.."

cd "$BACKEND_DIR"

echo "=== Stratego Database Migration Tool ==="
go run scripts/migrate.go
echo "======================================="
