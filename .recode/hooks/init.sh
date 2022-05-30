#!/bin/bash
set -euo pipefail

log () {
  echo -e "${1}" >&2
}

log "Downloading dependencies listed in go.mod"

go mod download
