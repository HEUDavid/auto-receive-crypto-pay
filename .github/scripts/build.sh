#!/bin/bash
set -e

# -----------------------------------------------------------------------------
# Script: build.sh
# Description: Builds the Go binary for the specified OS and Architecture.
# Environment Variables:
#   APP_NAME : Name of the application (e.g., receivepay)
#   GOOS     : Target Operating System
#   GOARCH   : Target Architecture
# -----------------------------------------------------------------------------

# Validate environment variables
if [[ -z "$APP_NAME" || -z "$GOOS" || -z "$GOARCH" ]]; then
  echo "Error: APP_NAME, GOOS, and GOARCH must be set."
  exit 1
fi

# Determine binary name with extension for Windows
BINARY_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
if [[ "$GOOS" == "windows" ]]; then
  BINARY_NAME="${BINARY_NAME}.exe"
fi

echo "--------------------------------------------------"
echo "Building ${BINARY_NAME} ..."
echo "OS: ${GOOS}, Arch: ${GOARCH}"
echo "--------------------------------------------------"

# Build the binary
go build -v -o "$BINARY_NAME" ./cmd

echo "Build complete: ${BINARY_NAME}"
