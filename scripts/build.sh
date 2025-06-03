#!/usr/bin/env bash
# -----------------------------------------------------------------------------
# Script: build.sh
# Author: Dheny (@furiatona on GitHub)
# Description: Interactive single-platform Go binary builder (outputs to build/)
# -----------------------------------------------------------------------------

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="${SCRIPT_DIR}/.."
ENV_FILE="${ROOT_DIR}/.env"
BUILD_DIR="${ROOT_DIR}/build"

if [[ -f "${ENV_FILE}" ]]; then
    echo -e "${YELLOW}🔧 Loading .env from ${ENV_FILE}${NC}"
    set -a
    source "${ENV_FILE}"
    set +a
else
    echo -e "${RED}❌ .env file not found at ${ENV_FILE}${NC}"
    exit 1
fi

if [[ -z "${APP_NAME:-}" ]]; then
    echo -e "${RED}❌ APP_NAME not set in .env${NC}"
    exit 1
fi

PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
)

echo -e "${BLUE}Select target platform:${NC}"
for i in "${!PLATFORMS[@]}"; do
  echo "  [$((i+1))] ${PLATFORMS[$i]}"
done

read -rp "$(echo -e "${CYAN}Enter number [1-${#PLATFORMS[@]}]: ${NC}")" PLATFORM_INDEX
PLATFORM_INDEX=$((PLATFORM_INDEX-1))

while [[ $PLATFORM_INDEX -lt 0 || $PLATFORM_INDEX -ge ${#PLATFORMS[@]} ]]; do
  echo -e "${RED}❌ Invalid selection${NC}"
  read -rp "$(echo -e "${CYAN}Enter number [1-${#PLATFORMS[@]}]: ${NC}")" PLATFORM_INDEX
  PLATFORM_INDEX=$((PLATFORM_INDEX-1))
done

IFS="/" read -r GOOS GOARCH <<< "${PLATFORMS[$PLATFORM_INDEX]}"
OUTPUT_NAME="${APP_NAME}"
OUTPUT_PATH="${BUILD_DIR}/${OUTPUT_NAME}"

if [[ -f "${OUTPUT_PATH}" ]]; then
  read -rp "$(echo -e "${YELLOW}⚠️  ${OUTPUT_NAME} already exists. Overwrite? [Y/n]: ${NC}")" CONFIRM
  CONFIRM=${CONFIRM:-Y}
  if [[ ! "${CONFIRM}" =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}❌ Skipping build.${NC}"
    exit 0
  fi
fi

echo -e "${CYAN}🔧 Building ${APP_NAME} for ${GOOS}/${GOARCH}...${NC}"
env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "${OUTPUT_PATH}" "${ROOT_DIR}"

if [[ $? -ne 0 ]]; then
  echo -e "${RED}❌ Build failed.${NC}"
  exit 1
else
  echo -e "${GREEN}✔️  Build succeeded, you can run your app with ./build/${APP_NAME}${NC}"
fi

# Check if 'build/' is in .gitignore
if ! grep -qx "build/" "${ROOT_DIR}/.gitignore"; then
  echo -e "${YELLOW}⚠️  Warning: 'build/' directory not excluded in .gitignore${NC}"
  echo -e "${BLUE}👉 Tip: Add to .gitignore:${NC} ${YELLOW}build/${NC}"
fi

echo -e "${GREEN}✅ Done.${NC}"