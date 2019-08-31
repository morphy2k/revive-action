#!/bin/sh

set -e

if [ ! -z "${GITHUB_TOKEN}" ];
then
  sh -c "cd $GITHUB_WORKSPACE && revive -config $INPUT_CONFIG -formatter ndjson ./... | revive-action"
else
  sh -c "cd $GITHUB_WORKSPACE && revive -config $INPUT_CONFIG -formatter friendly ./..."
fi
