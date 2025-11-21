#!/usr/bin/env bash
set -euo pipefail

API_URL=${API_URL:-http://localhost:8080}
COOKIEJAR="/tmp/stratego-seed-cookies.txt"

echo "Seeding development data to $API_URL"

# Clean previous cookiejar
rm -f "$COOKIEJAR"

# Helper to POST JSON and pretty-print response
post_json() {
  local url="$1"
  local data="$2"
  local cookie_opt="${3:-}"
  echo "POST $url"
  curl -sS $cookie_opt -H "Content-Type: application/json" -d "$data" "$url"
  echo -e "\n---"
}

# 1) Create user 'epicGamer' (backend will hash the password)
echo "Creating user 'epicGamer'..."
post_json "$API_URL/api/users/register" '{"username":"epicGamer","password":"password"}' "-c $COOKIEJAR" || true

# If registration above returned an error (user exists), try to login to obtain cookie
if ! grep -q "session_id" "$COOKIEJAR" 2>/dev/null; then
  echo "Attempting to login as 'epicGamer' to obtain session cookie..."
  post_json "$API_URL/api/users/login" '{"username":"epicGamer","password":"password"}' "-c $COOKIEJAR" || true
fi

# Verify cookie
if ! grep -q "session_id" "$COOKIEJAR" 2>/dev/null; then
  echo "WARNING: session cookie not found in $COOKIEJAR. Board setups creation will likely fail."
else
  echo "Session cookie obtained. Creating sample board setups..."

  # Create a couple of example board setups
  post_json "$API_URL/api/board-setups" '{"name":"Standard Starter","description":"A starter setup","setup_data":"{\"pieces\":[]}","is_default":true}' "-b $COOKIEJAR"

  post_json "$API_URL/api/board-setups" '{"name":"Aggressive","description":"Front-loaded pieces","setup_data":"{\"pieces\":[]}","is_default":false}' "-b $COOKIEJAR"

  echo "Sample board setups created."
fi

# Cleanup cookiejar
rm -f "$COOKIEJAR"

echo "Seeding complete."
