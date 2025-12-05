#!/bin/bash
set -e

# -----------------------------------------------------------------------------
# Script: package-release.sh
# Description: Packages the built binaries and resources into archives.
# Environment Variables:
#   APP_NAME : Name of the application
# -----------------------------------------------------------------------------

# Validate environment variables
if [[ -z "$APP_NAME" ]]; then
  echo "Error: APP_NAME must be set."
  exit 1
fi

OUTPUT_DIR="release-packages"
ARTIFACTS_DIR="all-artifacts"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "--------------------------------------------------"
echo "Packaging Release Assets"
echo "App Name: ${APP_NAME}"
echo "--------------------------------------------------"

# Dynamically find all binaries in the artifacts directory
# We expect files to be named like: ${APP_NAME}-${GOOS}-${GOARCH}[.exe]
find "$ARTIFACTS_DIR" -type f -name "${APP_NAME}-*" | while read -r BIN_PATH; do
  BIN_NAME=$(basename "$BIN_PATH")
  
  # Parse OS and Arch from the filename
  # 1. Remove .exe extension if present
  FILENAME="${BIN_NAME%.exe}"
  
  # 2. Remove APP_NAME prefix (and the hyphen)
  # suffix should be like "linux-amd64" or "windows-arm64"
  SUFFIX="${FILENAME#${APP_NAME}-}"
  
  # 3. Split by hyphen to get OS and Arch
  # Assuming GOOS and GOARCH don't contain hyphens themselves (standard Go ones don't)
  OS=$(echo "$SUFFIX" | cut -d'-' -f1)
  ARCH=$(echo "$SUFFIX" | cut -d'-' -f2)

  if [[ -z "$OS" || -z "$ARCH" ]]; then
     echo "Warning: Could not parse OS/Arch from $BIN_NAME. Skipping."
     continue
  fi

  echo "-> Packaging ${OS}/${ARCH} from $BIN_PATH..."

  # Create a temporary directory for packaging
  PKG_DIR="${APP_NAME}-${OS}-${ARCH}"
  mkdir -p "$PKG_DIR"

  # Copy binary and resources
  cp "$BIN_PATH" "$PKG_DIR/"
  cp -r conf static go.mod "$PKG_DIR/"

  # Create archive
  ARCHIVE_NAME=""
  if [[ "$OS" == "windows" ]]; then
    ARCHIVE_NAME="${PKG_DIR}.zip"
    zip -r "${OUTPUT_DIR}/${ARCHIVE_NAME}" "$PKG_DIR" > /dev/null
  else
    ARCHIVE_NAME="${PKG_DIR}.tar.gz"
    tar -czf "${OUTPUT_DIR}/${ARCHIVE_NAME}" "$PKG_DIR"
  fi

  echo "   Created: ${OUTPUT_DIR}/${ARCHIVE_NAME}"

  # Clean up temp dir
  rm -rf "$PKG_DIR"
done

echo "--------------------------------------------------"
echo "Release content prepared:"
ls -R "$OUTPUT_DIR"
echo "--------------------------------------------------"
