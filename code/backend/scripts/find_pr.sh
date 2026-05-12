#!/bin/bash
# Find the Pull Request ID for the current branch

REPO=$1
OWNER=$2
BRANCH=$3
TOKEN=$4

if [ -z "$TOKEN" ]; then
    echo "No GitHub token provided."
    exit 0
fi

PR_DATA=$(curl -s -H "Authorization: token $TOKEN" \
    "https://api.github.com/repos/$REPO/pulls?head=$OWNER:$BRANCH&state=open")

PR_ID=$(echo "$PR_DATA" | jq -r '.[0].number')

if [ "$PR_ID" != "null" ] && [ "$PR_ID" != "" ]; then
    echo "$PR_ID"
else
    exit 1
fi
