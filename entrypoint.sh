#!/bin/bash

set -e
set -o pipefail

cd "$GITHUB_WORKSPACE"

ACTION_VERSION=$(revive-action -version)
REVIVE_VERSION=$(revive -version | gawk '{match($0,"v[0-9].[0-9].[0-9]",a)}END{print a[0]}')

LINT_PATH="./..."

if [ ! -z "${INPUT_PATH}" ]; then LINT_PATH=$INPUT_PATH; fi

IFS=';' read -ra ADDR <<< "$INPUT_EXCLUDE"
for i in "${ADDR[@]}"; do
  EXCLUDES="$EXCLUDES -exclude="$i""
done

if [ ! -z "${INPUT_CONFIG}" ]; then CONFIG="-config=$INPUT_CONFIG"; fi

echo "ACTION: $ACTION_VERSION
REVIVE: $REVIVE_VERSION"

eval "revive $CONFIG $EXCLUDES -formatter ndjson $LINT_PATH | revive-action"
