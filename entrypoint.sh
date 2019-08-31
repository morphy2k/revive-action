#!/bin/sh

set -e

CONFIG=""

if [ ! -z "${INPUT_CONFIG}" ]; then CONFIG="-config=$INPUT_CONFIG"; fi

if [ ! -z "${GITHUB_TOKEN}" ];
then
  sh -c "cd $GITHUB_WORKSPACE && revive $CONFIG -formatter ndjson ./... | revive-action"
else
  echo "Annotations inactive. No GitHub token provided"
  sh -c "cd $GITHUB_WORKSPACE && revive $CONFIG -formatter friendly ./..."
fi
