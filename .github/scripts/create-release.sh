#!/bin/bash
set -e

# -----------------------------------------------------------------------------
# Script: create-release.sh
# Description: Creates a GitHub Release and uploads assets.
# Environment Variables:
#   GITHUB_TOKEN : GitHub Token for authentication
# -----------------------------------------------------------------------------

# Generate Release Tag
TAG_NAME=$(date +'v%Y.%m.%d-%H%M%S')

echo "--------------------------------------------------"
echo "Creating Release: $TAG_NAME"
echo "--------------------------------------------------"

RELEASE_NOTES="## Auto Receive Crypto Pay

### Release $TAG_NAME

Instructions: https://github.com/HEUDavid/auto-receive-crypto-pay"

# Create Release using gh CLI
gh release create "$TAG_NAME" release-packages/* \
  --title "Release $TAG_NAME" \
  --notes "$RELEASE_NOTES" \
  --draft=false

echo "Release $TAG_NAME created successfully."
