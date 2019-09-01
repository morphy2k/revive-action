#!/bin/bash

set -e

cd "$GITHUB_WORKSPACE"

IFS=';' read -ra ADDR <<< "$INPUT_EXCLUDE"
for i in "${ADDR[@]}"; do
  EXCLUDES="$EXCLUDES -exclude="$i""
done

if [ ! -z "${INPUT_CONFIG}" ]; then CONFIG="-config=$INPUT_CONFIG"; fi

if [ ! -z "${GITHUB_TOKEN}" ];
then
  sh -c "revive $CONFIG $EXCLUDES -formatter=ndjson ./... | revive-action"
else
  echo "Annotations inactive. No GitHub token provided"
  sh -c "revive $CONFIG $EXCLUDES -formatter friendly ./..."
fi
