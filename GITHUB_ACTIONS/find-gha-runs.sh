#!/bin/bash

# Script to find GitHub Actions runs for HEAD and the previous commit that had a run
# Outputs just the URLs for easy comparison

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null; then
    print_error "GitHub CLI (gh) is not installed. Please install it first."
    exit 1
fi

# Check if we're authenticated
if ! gh auth status &> /dev/null; then
    print_error "Not authenticated with GitHub CLI. Please run 'gh auth login' first."
    exit 1
fi

# Get current commit hash
print_status "Getting current commit (HEAD)..."
HEAD_COMMIT=$(git rev-parse HEAD)
HEAD_SHORT=$(git rev-parse --short HEAD)

print_status "Current commit: $HEAD_SHORT ($HEAD_COMMIT)"

# Check if HEAD has GitHub Actions runs
print_status "Checking for GitHub Actions runs on current commit..."
HEAD_RUNS=$(gh run list --commit "$HEAD_COMMIT" --json databaseId,url,conclusion,createdAt --jq '.[0]')

if [ "$HEAD_RUNS" = "null" ] || [ -z "$HEAD_RUNS" ]; then
    print_warning "No GitHub Actions runs found for current commit ($HEAD_SHORT)"
    print_status "This might be a release commit or the commit wasn't pushed to main/alpha/beta branch"
    HEAD_URL=""
else
    HEAD_RUN_ID=$(echo "$HEAD_RUNS" | jq -r '.databaseId')
    HEAD_URL=$(echo "$HEAD_RUNS" | jq -r '.url')
    HEAD_STATUS=$(echo "$HEAD_RUNS" | jq -r '.conclusion')
    HEAD_DATE=$(echo "$HEAD_RUNS" | jq -r '.createdAt')
    
    print_status "Found run for current commit: $HEAD_RUN_ID ($HEAD_STATUS)"
fi

# Find the previous commit that had a GitHub Actions run
print_status "Finding previous commit with GitHub Actions runs..."

PREV_COMMIT=""
PREV_URL=""
PREV_RUN_ID=""
PREV_STATUS=""
PREV_DATE=""

# Start checking from HEAD~1 and go back until we find a run
for i in {1..10}; do
    COMMIT_HASH=$(git rev-parse "HEAD~$i" 2>/dev/null || echo "")
    
    if [ -z "$COMMIT_HASH" ]; then
        print_warning "Reached end of commit history after checking $((i-1)) commits back"
        break
    fi
    
    COMMIT_SHORT=$(git rev-parse --short "HEAD~$i")
    print_status "Checking commit HEAD~$i: $COMMIT_SHORT"
    
    # Check if this commit has GitHub Actions runs
    RUNS=$(gh run list --commit "$COMMIT_HASH" --json databaseId,url,conclusion,createdAt --jq '.[0]')
    
    if [ "$RUNS" != "null" ] && [ -n "$RUNS" ]; then
        PREV_COMMIT="$COMMIT_HASH"
        PREV_RUN_ID=$(echo "$RUNS" | jq -r '.databaseId')
        PREV_URL=$(echo "$RUNS" | jq -r '.url')
        PREV_STATUS=$(echo "$RUNS" | jq -r '.conclusion')
        PREV_DATE=$(echo "$RUNS" | jq -r '.createdAt')
        
        print_status "Found previous commit with runs: $COMMIT_SHORT ($PREV_STATUS)"
        break
    else
        print_status "No runs found for $COMMIT_SHORT, checking next commit..."
    fi
done

# Output results
echo ""
echo "=========================================="
echo "GitHub Actions Run Comparison URLs"
echo "=========================================="
echo ""

if [ -n "$HEAD_URL" ]; then
    echo "Current Commit ($HEAD_SHORT):"
    echo "  URL: $HEAD_URL"
    echo "  Status: $HEAD_STATUS"
    echo "  Date: $HEAD_DATE"
    echo ""
else
    echo "Current Commit ($HEAD_SHORT):"
    echo "  No GitHub Actions runs found"
    echo ""
fi

if [ -n "$PREV_URL" ]; then
    PREV_SHORT=$(git rev-parse --short "$PREV_COMMIT")
    echo "Previous Commit ($PREV_SHORT):"
    echo "  URL: $PREV_URL"
    echo "  Status: $PREV_STATUS"
    echo "  Date: $PREV_DATE"
    echo ""
else
    echo "Previous Commit:"
    echo "  No GitHub Actions runs found in last 10 commits"
    echo ""
fi

# Summary
echo "=========================================="
echo "Summary:"
echo "=========================================="

if [ -n "$HEAD_URL" ] && [ -n "$PREV_URL" ]; then
    echo "✅ Both commits have GitHub Actions runs"
    echo "📊 You can compare the runs using the URLs above"
elif [ -n "$HEAD_URL" ]; then
    echo "⚠️  Only current commit has runs (previous commits may be release-only)"
elif [ -n "$PREV_URL" ]; then
    echo "⚠️  Only previous commit has runs (current commit may be release-only)"
else
    echo "❌ No GitHub Actions runs found for either commit"
fi

echo ""
echo "💡 Tip: Open both URLs in separate browser tabs to compare the runs"