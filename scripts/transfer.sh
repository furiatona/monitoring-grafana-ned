#!/usr/bin/env bash
# -----------------------------------------------------------------------------
# Script: transfer.sh
# Author: Dheny (@furiatona on GitHub)
# Description: Transfers all files from local BUILD_DIR to remote SERVER_FILE_DIR via SSH/SCP
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
    echo -e "${YELLOW}🔧 Loading environment from ${ENV_FILE}${NC}"
    set -a
    source "${ENV_FILE}"
    set +a
else
    echo -e "${RED}❌ .env file not found at ${ENV_FILE}${NC}"
    exit 1
fi

: "${SSH_HOST:?Missing SSH_HOST in .env or environment}"
: "${SSH_USER:?Missing SSH_USER in .env or environment}"
: "${SERVER_FILE_DIR:?Missing SERVER_FILE_DIR in .env or environment}"

SSH_PORT="${SSH_PORT:-22}"
USE_SSH_PASS=false
[[ -n "${SSH_PASS:-}" ]] && USE_SSH_PASS=true

transfer_with_retries() {
    local retries=5
    local delay=3
    local attempt=1

    while (( attempt <= retries )); do
        echo -e "${CYAN}Attempt ${attempt}/${retries} to transfer files...${NC}"

        if $USE_SSH_PASS; then
            sshpass -p "$SSH_PASS" scp -r -P "$SSH_PORT" "${BUILD_DIR}/." "${SSH_USER}@${SSH_HOST}:${SERVER_FILE_DIR}/" && return 0
        else
            scp -r -P "$SSH_PORT" "${BUILD_DIR}/." "${SSH_USER}@${SSH_HOST}:${SERVER_FILE_DIR}/" && return 0
        fi

        echo -e "${YELLOW}⚠️  Transfer attempt ${attempt} failed. Retrying in ${delay}s...${NC}"
        ((attempt++))
        sleep "$delay"
    done

    echo -e "${RED}❌ All ${retries} attempts failed to transfer files.${NC}"
    return 1
}

transfer_all_files() {
    echo -e "${CYAN}🚀 Transferring all files in ${BUILD_DIR} to ${SSH_USER}@${SSH_HOST}:${SERVER_FILE_DIR}${NC}"

    if [[ ! -d "${BUILD_DIR}" ]]; then
        echo -e "${RED}❌ Build directory not found: ${BUILD_DIR}${NC}"
        return 1
    fi

    transfer_with_retries
    echo -e "${GREEN}✅ All files transferred successfully to ${SERVER_FILE_DIR}${NC}"
}

main() {
    transfer_all_files
}

main "$@"