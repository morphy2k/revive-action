#!/bin/bash

set -e
set -o pipefail

cd "$GITHUB_WORKSPACE"

LINT_PATH="./..."

if [ ! -z "${INPUT_PATH}" ]; then LINT_PATH=$INPUT_PATH; fi

IFS=';' read -ra ADDR <<< "$INPUT_EXCLUDE"
for i in "${ADDR[@]}"; do
  EXCLUDES="$EXCLUDES -exclude="$i""
done

if [ ! -z "${INPUT_CONFIG}" ]; then CONFIG="-config=$INPUT_CONFIG"; fi

eval "revive $CONFIG $EXCLUDES -formatter ndjson $LINT_PATH | revive-action"
